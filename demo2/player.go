package demo2

import (
	"context"
	"io"
	"path/filepath"
	"sync"

	"github.com/hajimehoshi/oto"
	"github.com/pkg/errors"
)

type Player struct {
	Middlewares []Middleware

	codecs map[string]Codec
	mu     sync.Mutex
}

func NewPlayer() *Player {
	return &Player{
		codecs: make(map[string]Codec),
	}
}

func (p *Player) RegisterCodec(ext string, codec Codec) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.codecs[ext] = codec
}

func (p *Player) Run(ctx context.Context, path string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	ext := filepath.Ext(path)
	codec, ok := p.codecs[ext]
	if !ok {
		return errors.Errorf("no registered codec found for extension '%s'", ext)
	}

	dec, err := codec.Decoder(path)
	if err != nil {
		return err
	}

	player, err := oto.NewPlayer(dec.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer player.Close()

	var reader io.Reader = dec
	for _, m := range p.Middlewares {
		reader = m.Pipe(reader)
	}

	done := make(chan struct{})
	go func() {
		_, err = io.Copy(player, reader)
		close(done)
	}()

	select {
	case <-ctx.Done():
		player.Close()
		<-done
	case <-done:
	}

	return err
}

type Codec interface {
	Decoder(string) (Decoder, error)
}

type Decoder interface {
	io.Reader

	SampleRate() int
}

type Middleware interface {
	Pipe(io.Reader) io.Reader
}
