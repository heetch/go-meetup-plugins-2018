package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/heetch/go-meetup-plugins-2018/demo2"
	"github.com/heetch/go-meetup-plugins-2018/demo2/rpc"
	"github.com/heetch/go-meetup-plugins-2018/demo2/rpc/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	f := flag.String("f", "", "path to file")
	flag.Parse()

	if *f == "" {
		flag.Usage()
		os.Exit(2)
	}

	err := start(*f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
	}
}

func start(path string) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lis, err := net.Listen("tcp", ":")
	if err != nil {
		return err
	}

	p := demo2.NewPlayer()
	ps := rpc.NewPlayerService(p)

	errC := make(chan error, 2)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := runServer(ctx, lis, ps)
		if err != nil {
			errC <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()

		err := playFile(ctx, path, p, lis.Addr().String(), ps.OnCodecRegisteredC)
		if err != nil {
			errC <- err
		}
	}()

	select {
	case <-quit:
		cancel()
	case err = <-errC:
		cancel()
	case <-ctx.Done():
	}

	wg.Wait()
	return err
}

func runServer(ctx context.Context, l net.Listener, ps *rpc.PlayerService) error {
	server := grpc.NewServer()
	pb.RegisterPlayerServer(server, ps)
	reflection.Register(server)
	go func() {
		<-ctx.Done()
		server.Stop()
	}()

	return server.Serve(l)
}

func playFile(ctx context.Context, path string, p *demo2.Player, playerAddr string, ch chan demo2.Codec) error {
	loader := demo2.PluginLoader{Path: "./bin", PlayerAddr: playerAddr}
	err := loader.Load()
	if err != nil {
		return errors.Wrap(err, "player: failed to load plugins")
	}
	defer loader.Stop()
	var registeredPlugins int
	lctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	for registeredPlugins < len(loader.Plugins) {
		select {
		case <-lctx.Done():
			if lctx.Err() == context.DeadlineExceeded {
				return errors.Errorf("timeout occured when waiting for plugins to register. %d out of %d successfully registered", registeredPlugins, len(loader.Plugins))
			}

			return lctx.Err()
		case <-ch:
			registeredPlugins++
		}
	}

	return p.Run(ctx, path)
}
