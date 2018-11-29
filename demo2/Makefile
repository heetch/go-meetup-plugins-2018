NAME := demo2

.PHONY: all $(NAME) test testrace run build

all: $(NAME)

build: $(NAME)

$(NAME):
	go build ./cmd/$@
	go build -o ./bin/player-mp3codec ./cmd/mp3codec
	go build -o ./bin/player-wavCodec ./cmd/wavCodec

test:
	go test -v -cover -timeout=1m ./...

testrace:
	go test -v -race -cover -timeout=2m ./...

run: build
	./demo2 -f sound1.mp3
