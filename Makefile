COMPOSEV2 := $(shell docker compose version 2> /dev/null)

ifdef COMPOSEV2
    COMMAND_DOCKER=docker compose
else
    COMMAND_DOCKER=docker-compose
endif

up:
	$(COMMAND_DOCKER) up -d

logs:
	docker logs --tail 50 -f go_wbc

down:
	$(COMMAND_DOCKER) down