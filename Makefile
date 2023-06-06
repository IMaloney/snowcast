SERVER_BINARY=snowcast_server
CLIENT_BINARY=snowcast_control
LISTENER_BINARY=snowcast_listener

windows:
	go mod download
	go get -u -f -v all
	go build -o ${SERVER_BINARY}.exe ./cmd/snowcast_server/snowcast_server.go 
	go build -o ${CLIENT_BINARY}.exe ./cmd/snowcast_control/snowcast_control.go
	go build -o ${LISTENER_BINARY}.exe ./cmd/snowcast_listener/snowcast_listener.go

build:
	go mod download
	go get -u -v -f all
	go build -o ${SERVER_BINARY} ./cmd/snowcast_server/snowcast_server.go
	go build -o ${CLIENT_BINARY} ./cmd/snowcast_control/snowcast_control.go
	go build -o ${LISTENER_BINARY} ./cmd/snowcast_listener/snowcast_listener.go

server:
	go build -o ${SERVER_BINARY} ./cmd/snowcast_server/snowcast_server.go

client:
	go build -o ${CLIENT_BINARY} ./cmd/snowcast_control/snowcast_control.go 

listener:
	go build -o ${LISTENER_BINARY} ./cmd/snowcast_listener/snowcast_listener.go

test:
	go test -v ./pkg/* --count=1

clean:
	go clean
	rm -f ${SERVER_BINARY}
	rm -f ${CLIENT_BINARY}
	rm -f ${LISTENER_BINARY}