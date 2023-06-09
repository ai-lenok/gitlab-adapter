BINARY_NAME=gitlab-adapter

build: 
	mkdir -p out/bin
	go build -o out/bin/$(BINARY_NAME) .

clean:
	rm -fr ./out

start-server:
	go run *.go start-server