.PHONY: all dev clean build env-up env-down run

##### ENV
env-up:
	@echo "Start environment ..."
	@cd fixtures && docker-compose -f docker-compose-local.yaml up -d
	@echo "Environment up"

env-down:
	@echo "Stop environment ..."
	@cd fixtures && docker-compose -f docker-compose-local.yaml down
	@echo "Environment down"

##### CLEAN
clean: env-down
	@echo "Clean up ..."
	@docker rm -f -v `docker ps -a --no-trunc | grep "dev-peer" | cut -d ' ' -f 1` 2>/dev/null || true
	@docker rmi `docker images --no-trunc | grep "dev-peer" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "Clean up done"

