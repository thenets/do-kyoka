IMAGE_TAG=thenets/do-kyoka

go-run:
	go run main.go

build:
	docker build --pull --rm \
		-f "Dockerfile" \
		-t $(IMAGE_TAG) "."

run:
	docker run --rm -it \
		-e FIREWALL_NAME=$(FIREWALL_NAME) \
		-e FIREWALL_TAG=$(FIREWALL_TAG) \
		-e DO_API_TOKEN=$(DO_API_TOKEN) \
		$(IMAGE_TAG)