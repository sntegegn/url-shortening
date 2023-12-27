#!/bin/bash

docker compose down

docker rm -vf $(docker ps -aq)
docker rmi -f $(docker images -aq)
