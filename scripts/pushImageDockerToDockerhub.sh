#!/bin/bash

docker build -f Dockerfile_Image --platform linux/amd64 . -t=superj80820/golang-process-image-base
docker push superj80820/golang-process-image-base