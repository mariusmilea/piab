NAME ?= piab
BUILD=.build

all: build

clean:
	rm -rf bin *.rpm $(BUILD) $(NAME)

build: clean
	cd app; go get .; go build    

up:
	docker-compose up --build

