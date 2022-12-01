#!/usr/bin/env bash

docker run -d --name db \
              --platform=linux/x86_64 \
              -p 5432:5432 \
              romuloslv/customdb:1.0

docker run -d --name pg \
              --platform=linux/x86_64 \
              -e PGADMIN_DEFAULT_EMAIL=admin@admin.com \
              -e PGADMIN_DEFAULT_PASSWORD=test12321 \
              -p 80:80 \
              --link db \
              dpage/pgadmin4

docker run -d --name api \
              --platform=linux/x86_64 \
              -e APP_POSTGRES_HOST=db \
              -e APP_POSTGRES_PASSWORD=test12321 \
              -p 8080:8080 \
              --link db \
              romuloslv/customapp:1.0