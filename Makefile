.PHONY: all main go-judge clean

all: main go-judge

main: main.go
	go build -o main

go-judge:
	git clone git@github.com:criyle/go-judge.git && \
	cd ./go-judge && go generate ./cmd/go-judge/version \
	&& CGO_ENABLED=0 go build -v -tags grpcnotrace,nomsgpack -o go-judge ./cmd/go-judge

run:
	./main ./go-judge/go-judge

clean:
	rm -f main ./go-judge/go-judge
