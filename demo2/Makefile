NAME := demo2

.PHONY: all $(NAME) test testrace run build mp3 wav clean

all: $(NAME)

build: $(NAME)

$(NAME):
	go build ./cmd/$@

mp3:
	go build -o ./bin/player-mp3codec ./cmd/mp3codec

wav:
	go build -o ./bin/player-wavCodec ./cmd/wavCodec

clean:
	rm $(NAME)
	rm ./bin/player-mp3codec
	rm ./bin/player-wavCodec

test:
	go test -v -cover -timeout=1m ./...

testrace:
	go test -v -race -cover -timeout=2m ./...
