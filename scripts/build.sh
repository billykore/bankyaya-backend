#!/usr/bin/env bash

DOCKERFILE="$(pwd)/Dockerfile"

# check Dockerfile.
if [[ ! -f "$DOCKERFILE" ]]; then
  echo >&2 "error: Dockerfile not found"
  exit 1
fi

IMAGE_NAME=$1

# build Docker image.
echo "Build image..."
docker build -t $IMAGE_NAME -f $DOCKERFILE .

# push image to DockerHub.
echo "Push image..."
docker push $IMAGE_NAME

# delete local Docker image.
echo "Delete image $IMAGE_NAME..."
docker image rm $IMAGE_NAME

