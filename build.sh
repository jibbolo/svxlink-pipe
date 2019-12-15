#!/bin/sh

docker build -t jibbolo/svxlink-pipe:build . -f Dockerfile.build

docker create --name extract-svxlink-pipe jibbolo/svxlink-pipe:build  
docker cp extract-svxlink-pipe:go/src/github.com/jibbolo/svxlink-pipe/application ./svxlink-pipe-docker
docker rm -f extract-svxlink-pipe