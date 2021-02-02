SHELL := /bin/bash

export PROJECT = go-party-finder

# ==============================================================================
# Development

run: up dev

up:
	docker-compose up -d db

dev:
	go run ./cmd/web

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

# ==============================================================================
# Running tests within the local computer

test:
	go test ./... -count=1
	staticcheck ./...
