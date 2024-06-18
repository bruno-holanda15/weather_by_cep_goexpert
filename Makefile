COMPOSEV2 := $(shell docker compose version 2> /dev/null)

ifdef COMPOSEV2
    COMMAND_DOCKER=docker compose
else
    COMMAND_DOCKER=docker-compose
endif

up:
	$(COMMAND_DOCKER) up -d

logs-cep-validator:
	docker logs --tail 50 -f go_wbc1

logs-wbc:
	docker logs --tail 50 -f go_wb2

down:
	$(COMMAND_DOCKER) down

tests:
	go test --count=1 -coverprofile=coverage.out $(shell go list ./... | grep internal/domain)