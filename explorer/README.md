# Explorer stub
Stub that stores messages to explorer

Handlers
```
GET /messages
GET /count
```
To rebuild use
```
docker build -t explorer-mock .
docker tag explorer-mock <registry>:latest
docker push <registry>:latest
```