package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/cryptix/wav"
	"github.com/heetch/go-meetup-plugins-2018/demo2/rpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	playerAddr := flag.String("player-addr", "", "address of the player gRPC server")
	flag.Parse()

	if *playerAddr == "" {
		flag.Usage()
		os.Exit(2)
	}

	err := runPlugin(*playerAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "plugin error: '%s'\n", err)
		os.Exit(2)
	}
}

func runPlugin(playerAddr string) error {
	clientConn, err := grpc.Dial(playerAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	defer clientConn.Close()
	client := pb.NewPlayerClient(clientConn)

	lis, err := net.Listen("tcp", ":")
	if err != nil {
		return err
	}
	defer lis.Close()

	go func() {
		_, err = client.RegisterCodec(context.Background(), &pb.RegisterCodecRequest{Addr: lis.Addr().String(), Ext: ".wav"})
		if err != nil {
			log.Printf("plugin: error while trying to register codec to player server: '%s'", err)
			lis.Close()
		}
	}()

	s := grpc.NewServer()
	pb.RegisterCodecServer(s, new(server))
	reflection.Register(s)

	return s.Serve(lis)
}

type server struct {
}

func (s *server) AudioFileMetadata(ctx context.Context, req *pb.AudioFileMetadataRequest) (*pb.AudioFileMetadataResponse, error) {
	f, err := os.Open(req.Path)
	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	d, err := wav.NewReader(f, fi.Size())
	if err != nil {
		return nil, err
	}

	return &pb.AudioFileMetadataResponse{Length: int64(d.GetFile().SoundSize), SampleRate: int32(d.GetSampleRate())}, nil
}

func (s *server) DecodeAudioFile(req *pb.DecodeAudioFileRequest, stream pb.Codec_DecodeAudioFileServer) error {
	f, err := os.Open(req.Path)
	if err != nil {
		return err
	}

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	d, err := wav.NewReader(f, fi.Size())
	if err != nil {
		return err
	}

	r, err := d.GetDumbReader()
	if err != nil {
		return err
	}

	buf := make([]byte, 512)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF {
			return nil
		}

		err = stream.Send(&pb.DecodeAudioFileResponse{
			Chunk: buf[:n],
		})
		if err != nil {
			return err
		}
	}
}
