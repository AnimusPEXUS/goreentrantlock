export GONOPROXY=https://github.com/AnimusPEXUS/*

all: get build

get:
		go get -u -v "./..."
		go mod tidy

build:
		go build

