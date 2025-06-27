# Docker Compose

Go module for `framework` packages that wrap `docker-compose.yml` files.

This module is separated from the `framework`, because `testcontainers-go` module that adds Docker Compose support pulls in a lot of dependencies and we want to limit the blast radius as much as possible.