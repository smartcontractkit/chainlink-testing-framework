# Docker

We are not removing volumes and images when you are working locally to allow you to debug, however, to clean up some space use:
```
docker volume prune -f
docker system prune -f
```
