#!/bin/bash

# Travis Docker deploy script
IMAGE_NAME="boynux/squid-exporter"
IMAGE_TAG=${TRAVIS_TAG:-latest}

docker --version
docker build -t $IMAGE_NAME:$IMAGE_TAG .
echo $DOCKER_API_KEY | docker login -u boynux --password-stdin
docker push $IMAGE_NAME:$IMAGE_TAG
