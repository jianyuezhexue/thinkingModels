#!/bin/bash
git pull
docker compose up --build -d
docker rmi -f $(docker images | grep "<none>" | awk '{print $3}')  # ‌:ml-citation{ref="1,2" data="citationList"}
docker system prune --volumes -f
