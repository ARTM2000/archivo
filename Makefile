# Variables
## Development Variables
SERVER_AIR_CONF="./config/server.air.toml"
AGENT_AIR_CONF="./config/agent.air.toml"

dev_server:
	@trap 'rm -f ./tmp/server' EXIT; air -c ${SERVER_AIR_CONF}

build_server:
	@go build -o ./build/agent ./cmd/agent

dev_agent:
	@trap 'rm -f ./tmp/agent' EXIT; air -c ${AGENT_AIR_CONF}

build_agent:
	@go build -o ./build/agent ./cmd/agent

format:
	@gofmt -l -s -w .	
