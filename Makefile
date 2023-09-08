

all: clean server agent


clean:
	rm -rf cmd/server/server cmd/agent/agent

server:
	go build -o cmd/server/server cmd/server/server_main.go

agent:
	go build -o cmd/agent/agent cmd/agent/agent_main.go

check1:
	metricstest -test.v -test.run=^TestIteration1$ -binary-path=cmd/server/server


check2:
	metricstest -test.v -test.run=^TestIteration2$ -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent