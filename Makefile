IMAGE_TAG=thenets/do-kyoka

LOAD_ENVS=FIREWALL_NAME=do-kyoka \
	FIREWALL_TAG=do-kyoka \
	DO_API_TOKEN=$$(secret-tool lookup thenets_dev do_api_token) \
	SENTRY_DSN=$$(secret-tool lookup thenets_dev do_kyoka_sentry_dsn)

go-run:
	$(LOAD_ENVS) go run main.go

build:
	docker build --pull --rm \
		-f "Dockerfile" \
		-t $(IMAGE_TAG) "."

run:
	$(LOAD_ENVS) docker run --rm -it \
		-e FIREWALL_NAME=$(FIREWALL_NAME) \
		-e FIREWALL_TAG=$(FIREWALL_TAG) \
		-e DO_API_TOKEN=$(DO_API_TOKEN) \
		$(IMAGE_TAG)

load-envs:
	@echo $(LOAD_ENVS)

set-sentry-dsn:
	secret-tool store \
		--label "TheNets Dev: Sentry DSN" \
		thenets_dev do_kyoka_sentry_dsn $(SENTRY_DSN)

set-do-api-token:
	secret-tool store \
		--label "TheNets Dev: DigitalOcean API Token" \
		thenets_dev do_api_token $(DO_API_TOKEN)
