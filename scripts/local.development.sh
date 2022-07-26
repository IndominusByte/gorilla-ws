#!/bin/bash

export COMPOSE_IGNORE_ORPHANS=True

export BACKEND_IMAGE=gorilla-ws
export BACKEND_IMAGE_TAG=development
export BACKEND_CONTAINER=gorilla-ws-development
export BACKEND_HOST=gorilla-ws.service
export BACKEND_STAGE=development

docker build -t "$BACKEND_IMAGE:$BACKEND_IMAGE_TAG" -f ./manifest-docker/Dockerfile.dev ./
docker-compose -f ./manifest/docker-compose.development.yaml up -d --build
