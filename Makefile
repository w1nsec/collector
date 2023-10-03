AGENT=cmd/agent/agent
SERVER=cmd/server/server
PORT=9000



all: clean server agent


clean:
	rm -rf ${SERVER} ${AGENT}

server:
	go build -o ${SERVER} cmd/server/server_main.go

agent:
	go build -o ${AGENT} cmd/agent/agent_main.go

check1:
	metricstest -test.v -test.run=^TestIteration1$$ -binary-path=${SERVER}


check2:
	metricstest -test.v -test.run=^TestIteration2[AB]$$ -binary-path=${SERVER} -agent-binary-path=${AGENT} -source-path=./

check3:
	metricstest -test.v -test.run=^TestIteration3$$ -binary-path=${SERVER} -agent-binary-path=${AGENT} -source-path=./ -server-port=${PORT}

check4:
	metricstest -test.v -test.run=^TestIteration4$$ -binary-path=${SERVER} -agent-binary-path=${AGENT} -source-path=./ -server-port=${PORT}

check5:
	metricstest -test.v -test.run=^TestIteration5$$ -binary-path=${SERVER} -agent-binary-path=${AGENT} -source-path=./ -server-port=${PORT}

check6:
	metricstest -test.v -test.run=^TestIteration6$$ -binary-path=${SERVER} -agent-binary-path=${AGENT} -source-path=./ -server-port=${PORT}