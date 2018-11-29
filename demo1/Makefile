NAME := demo1

.PHONY: all $(NAME) test testrace run build

all: $(NAME)

build: $(NAME)

$(NAME):
	@ go build ./cmd/$@
	@ go build -o ./bin/player-plugin1 ./cmd/plugin1
	@ go build -o ./bin/player-plugin2 ./cmd/plugin2

run: build
	./demo1
