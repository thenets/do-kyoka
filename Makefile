IMAGE_TAG=thenets/do-kyoka

LOAD_ENVS=FIREWALL_NAME=do-kyoka \
	FIREWALL_TAG=do-kyoka \
	DO_API_TOKEN=$$(make -s get-do-api-token) \
	SENTRY_DSN=$$(make -s get-sentry-dsn)

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

get-sentry-dsn:
	@if [ -x "$$(command -v secret-tool)" ]; then \
		secret-tool lookup thenets_dev do_kyoka_sentry_dsn; \
	fi

set-do-api-token:
	secret-tool store \
		--label "TheNets Dev: DigitalOcean API Token" \
		thenets_dev do_api_token $(DO_API_TOKEN)

get-do-api-token:
	@if [ -x "$$(command -v secret-tool)" ]; then \
		secret-tool lookup thenets_dev do_api_token; \
	fi