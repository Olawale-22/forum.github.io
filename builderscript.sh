#!/bin/bash

# Set the name and tag for the Docker image
IMAGE_NAME="olawale-goforum-app"
IMAGE_TAG="latest"

# Build the Go program
go build -o olawale-goforum-app .

# Build the Docker image
docker build -t $IMAGE_NAME:$IMAGE_TAG .

# Remove the Go executable
rm olawale-goforum-app

# Create a Docker container
docker create --name $IMAGE_NAME $IMAGE_NAME:$IMAGE_TAG
