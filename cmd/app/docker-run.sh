#!/usr/bin/env bash

docker run -d --name db \
              -p 5432:5432 \
              romuloslv/customdb:1.0

docker run -d --name pg \
              -e PGADMIN_DEFAULT_EMAIL=admin@admin.com \
              -e PGADMIN_DEFAULT_PASSWORD=postgres \
              -p 80:80 \
              --link db \
              dpage/pgadmin4

docker run -d --name api \
              -e APP_POSTGRES_HOST=db \
              -e APP_POSTGRES_PASSWORD=postgres \
              -p 8080:8080 \
              --link db \
              romuloslv/customapp:1.0