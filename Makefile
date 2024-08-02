BINARY_NAME = kado
VERSION_FILE = VERSION
PYTHON_SCRIPT = bump_version.py
VERSION := $(shell cat $(VERSION_FILE))
GITHUB_USERNAME = janpreet
GITHUB_REPOSITORY = kado
DOCKER_IMAGE = ghcr.io/janpreet/kado

GOARCHES = amd64 arm64

build-all:
	@for os in $(GOOSES); do \
		for arch in $(GOARCHES); do \
			if [ "$$os" != "windows" ] || [ "$$arch" != "arm64" ]; then \
				echo "Building for $$os/$$arch"; \
				GOOS=$$os GOARCH=$$arch go build -o $(BINARY_NAME)-$$os-$$arch$(if $(filter windows,$$os),.exe,) .; \
			fi; \
		done; \
	done

build:
	go build -o $(BINARY_NAME) .

docker-build: build
	docker build -t $(DOCKER_IMAGE):latest -t $(DOCKER_IMAGE):$(VERSION) .

docker-push:
	docker push $(DOCKER_IMAGE):latest
	docker push $(DOCKER_IMAGE):$(VERSION)

version-bump:
	python3 $(PYTHON_SCRIPT)

tag-and-push: version-bump
	git add $(VERSION_FILE)
	git commit -m "Bump version to $(VERSION)"
	git tag v$(VERSION)
	git push origin main --tags

clean:
	rm -f $(BINARY_NAME)*

test:
	go test ./...

all: clean build-all docker-build docker-push tag-and-push

github-release:
	@echo "Creating GitHub release for v$(VERSION)"
	@gh release create v$(VERSION) \
		--title "Release $(VERSION)" \
		--notes "Release notes for version $(VERSION)" \
		$(BINARY_NAME)-*

.PHONY: build-all build docker-build docker-push version-bump tag-and-push clean test all github-release