# Variables
## Development Variables
SERVER_AIR_CONF="./config/server.air.toml"
AGENT_AIR_CONF="./config/agent.air.toml"
DASHBOARD_BASE_API="/api/v1"

dev_server:
	@trap 'rm -f ./tmp/server' EXIT; air -c ${SERVER_AIR_CONF}

build_server:
	@bash ./scripts/build_dashboard.bash ${DASHBOARD_BASE_API} && go build -o ./build/archivo ./cmd/server

dev_agent:
	@trap 'rm -f ./tmp/agent' EXIT; air -c ${AGENT_AIR_CONF}

build_agent:
	@go build -o ./build/agent ./cmd/agent

format:
	@gofmt -l -s -w . && go mod tidy

deploy_compose_dbonly:
	@docker-compose -f ./docker-compose.dbonly.yaml up

release_archivo:
	@bash ./scripts/build_dashboard.bash ${DASHBOARD_BASE_API} && bash ./scripts/build_cli.bash archivo ${PWD}/cmd/server/main.go

release_agent:
	@bash ./scripts/build_cli.bash agent ${PWD}/cmd/agent/main.go
