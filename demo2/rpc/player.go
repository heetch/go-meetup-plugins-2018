package rpc

import (
	"context"
	"log"

	"github.com/heetch/go-meetup-plugins-2018/demo2"
	"github.com/heetch/go-meetup-plugins-2018/demo2/rpc/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type PlayerService struct {
	Player             *demo2.Player
	OnCodecRegisteredC chan demo2.Codec
}

func NewPlayerService(p *demo2.Player) *PlayerService {
	return &PlayerService{
		Player:             p,
		OnCodecRegisteredC: make(chan demo2.Codec),
	}
}

func (p *PlayerService) RegisterCodec(ctx context.Context, req *pb.RegisterCodecRequest) (*pb.RegisterCodecResponse, error) {
	conn, err := grpc.Dial(req.Addr, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to codec plugin")
	}

	cd := &codec{
		client: pb.NewCodecClient(conn),
	}
	p.Player.RegisterCodec(req.Ext, cd)

	p.OnCodecRegisteredC <- cd

	log.Printf("Codec for '%s' extension successfully registered", req.Ext)
	return new(pb.RegisterCodecResponse), nil
}

type codec struct {
	addr   string
	client pb.CodecClient
}

func (c *codec) Decoder(path string) (demo2.Decoder, error) {
	resp, err := c.client.AudioFileMetadata(context.Background(), &pb.AudioFileMetadataRequest{Path: path})
	if err != nil {
		return nil, err
	}

	stream, err := c.client.DecodeAudioFile(context.Background(), &pb.DecodeAudioFileRequest{Path: path})
	if err != nil {
		return nil, err
	}

	return &decoder{
		sampleRate: int(resp.SampleRate),
		stream:     stream,
	}, nil
}

type decoder struct {
	sampleRate int
	stream     pb.Codec_DecodeAudioFileClient
}

func (d *decoder) SampleRate() int {
	return d.sampleRate
}

func (d *decoder) Read(p []byte) (n int, err error) {
	resp, err := d.stream.Recv()
	if err != nil {
		return 0, err
	}

	copy(p, resp.Chunk)
	return len(resp.Chunk), nil
}
