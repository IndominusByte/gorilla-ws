#!/bin/bash

export COMPOSE_IGNORE_ORPHANS=True

export BACKEND_IMAGE=gorilla-ws
export BACKEND_IMAGE_TAG=production
export BACKEND_CONTAINER=gorilla-ws-production
export BACKEND_HOST=gorilla-ws.service
export BACKEND_STAGE=production

docker build -t "$BACKEND_IMAGE:$BACKEND_IMAGE_TAG" -f ./manifest-docker/Dockerfile.prod ./
docker-compose -f ./manifest/docker-compose.production.yaml up -d --build
