AGENT=cmd/agent/agent
SERVER=cmd/server/server
LINTER=cmd/staticlint/staticlint

PORT=16738
FSPATH=/tmp/temp.db

DB_ADDR=localhost
DB_USER=user
DB_PASS=password
DB=mydb

KEY=supersecret

GIT_VERSION=`git log --oneline | head -n1 | cut -f1 -d' '`
COMMIT_FLAG=github.com/w1nsec/collector/internal/app/agent.buildCommit=${GIT_VERSION}
DATE=`date -u '+%Y-%m-%d %H:%M:%S'`
DATE_FLAG=github.com/w1nsec/collector/internal/app/agent.buildDate=${DATE}



all: clean server agent


clean:
	rm -rf ${SERVER} ${AGENT}

server:
	go build -o ${SERVER} cmd/server/server_main.go

agent:
	go build -o ${AGENT} -ldflags "-s -w -X \"${COMMIT_FLAG}\" -X \"${DATE_FLAG}\"" cmd/agent/agent_main.go

linter:
	go build -o ${LINTER} cmd/staticlint/staticlint.go

# Checks
check1:
	metricstest -test.v -test.run=^TestIteration1$$ -binary-path=${SERVER}


check2:
	metricstest -test.v -test.run=^TestIteration2[AB]*$$             -source-path=.             -agent-binary-path=cmd/agent/agent

check3:
	metricstest -test.v -test.run=^TestIteration3$$ -binary-path=${SERVER} -agent-binary-path=${AGENT} -source-path=./ -server-port=${PORT}

check4:
	metricstest -test.v -test.run=^TestIteration4$$ -binary-path=${SERVER} -agent-binary-path=${AGENT} -source-path=./ -server-port=${PORT}

check5:
	metricstest -test.v -test.run=^TestIteration5$$ \
				-binary-path=${SERVER} \
				-agent-binary-path=${AGENT} \
				-source-path=./ -server-port=${PORT}

check6:
	metricstest -test.v -test.run=^TestIteration6$$ \
                -agent-binary-path=${AGENT} \
                -binary-path=${SERVER} \
                -server-port=${PORT} \
                -source-path=.

check7:
	metricstest -test.v -test.run=^TestIteration7$$ \
                -agent-binary-path=${AGENT} \
                -binary-path=${SERVER} \
                -server-port=${PORT} \
                -source-path=.

check8:
	metricstest -test.v -test.run=^TestIteration8$$ \
                -agent-binary-path=${AGENT} \
                -binary-path=${SERVER} \
                -server-port=${PORT} \
                -source-path=.

check9:
	metricstest -test.v -test.run=^TestIteration9$$ \
                -agent-binary-path=${AGENT} \
                -binary-path=${SERVER} \
                -server-port=${PORT} \
                -file-storage-path=${FSPATH} \
                -source-path=.


check10:
	metricstest -test.v -test.run=^TestIteration10[AB]$$ \
                -agent-binary-path=${AGENT} \
                -binary-path=${SERVER} \
                -server-port=${PORT} \
                -file-storage-path=${FSPATH} \
                -source-path=. \
                -database-dsn="${DB_USER}:${DB_PASS}@${DB_ADDR}/${DB}"


check11:
	metricstest -test.v -test.run=^TestIteration11$$ \
                -agent-binary-path=${AGENT} \
                -binary-path=${SERVER} \
                -server-port=${PORT} \
                -file-storage-path=${FSPATH} \
                -source-path=. \
                -database-dsn="postgres://${DB_USER}:${DB_PASS}@${DB_ADDR}/${DB}"


check12:
	metricstest -test.v -test.run=^TestIteration12$$ \
                -agent-binary-path=${AGENT} \
                -binary-path=${SERVER} \
                -server-port=${PORT} \
                -file-storage-path=${FSPATH} \
                -source-path=. \
                -database-dsn="postgres://${DB_USER}:${DB_PASS}@${DB_ADDR}/${DB}"

check13:
	metricstest -test.v -test.run=^TestIteration13$$ \
                -agent-binary-path=${AGENT} \
                -binary-path=${SERVER} \
                -server-port=${PORT} \
                -file-storage-path=${FSPATH} \
                -source-path=. \
                -database-dsn="postgres://${DB_USER}:${DB_PASS}@${DB_ADDR}/${DB}"

check14:
	metricstest -test.v -test.run=^TestIteration14$$ \
                -agent-binary-path=${AGENT} \
                -binary-path=${SERVER} \
                -server-port=${PORT} \
                -file-storage-path=${FSPATH} \
                -source-path=. \
                -database-dsn="postgres://${DB_USER}:${DB_PASS}@${DB_ADDR}/${DB}" \
                -key="${KEY}"


staticcheck:
	cmd/staticlint ./...