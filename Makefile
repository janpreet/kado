# Makefile

# Variables
BINARY_NAME = kado
VERSION_FILE = VERSION
PYTHON_SCRIPT = bump_version.py
VERSION := $(shell cat $(VERSION_FILE))

# Define your GitHub username and repository name
GITHUB_USERNAME = janpreet
GITHUB_REPOSITORY = kado
DOCKER_IMAGE = ghcr.io/janpreet/kado

# Build the binary
build:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) main.go

# Build the Docker image
docker-build: build
	docker build -t $(DOCKER_IMAGE):latest -t $(DOCKER_IMAGE):$(VERSION) .

# Push the Docker image to GitHub Packages
docker-push:
	docker push $(DOCKER_IMAGE):latest
	docker push $(DOCKER_IMAGE):$(VERSION)

# Bump version, tag, and push
version-bump:
	python3 $(PYTHON_SCRIPT)

tag-and-push: version-bump
	git add $(VERSION_FILE)
	git commit -m "Bump version to $(VERSION)"
	git tag v$(VERSION)
	git push origin main --tags

# Clean up the build
clean:
	rm -f $(BINARY_NAME)

# Default target
all: clean build docker-build docker-push tag-and-push
