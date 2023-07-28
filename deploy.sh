#!/bin/bash

# Check if we're in the build directory and warn the user if they are
current_directory=$(basename "$PWD")
if [ "$current_directory" = "build" ]
then
    echo "This script should not be run from the build directory. Please move to the root directory of the project."
    exit 1
fi

read -p "Docker Hub username: " DOCKER_USERNAME
read -sp "Docker Hub password: " DOCKER_PASSWORD
echo
read -p "Docker Hub repository name: " DOCKER_REPO_NAME
echo

echo "Logging in to Docker Hub..."
echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USERNAME" --password-stdin
if [ $? -ne 0 ]; then
    echo "Docker login failed."
    exit 1
fi
echo "Login Succeeded"
echo

echo "Building Docker image..."
docker build -t "$DOCKER_USERNAME/$DOCKER_REPO_NAME:latest" .
if [ $? -ne 0 ]; then
    echo "Docker build failed."
    exit 1
fi
echo "Build Succeeded"
echo

echo "Pushing Docker image to Docker Hub..."
docker push "$DOCKER_USERNAME/$DOCKER_REPO_NAME:latest"
if [ $? -ne 0 ]; then
    echo "Docker push failed."
    exit 1
fi
echo "Push Succeeded"
