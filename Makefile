

all: clean server

clean:
	rm -rf cmd/server/server cmd/agent/agent

server:
	go build -o cmd/server/server cmd/server/server.go



check1:
	metricstest -test.v -test.run=^TestIteration1$ -binary-path=cmd/server/server