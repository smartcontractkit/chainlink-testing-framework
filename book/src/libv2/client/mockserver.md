# MockServer

This client is reponsible for simplifying interaction with MockServer during test execution (e.g. to dynamically create or modify mocked responses).

## Initialising
There are three ways of intializing the client:
* providing the URL
* providing pointer to `Environment` (CTF's representation of a k8s environment)
* providing pointer to raw `MockserverConfig`

```go
var myk8sEnv *environment.Environment

// ... create k8s environment

k8sMockserver := ConnectMockServer(myk8sEnv)

mockServerUrl := "http://my.mockserver.instance.io"
myMockServer := ConnectMockServerURL(mockServerUrl)

customConfig := &MockserverConfig{
    LocalURL: mockServerUrl,
    ClusterURL: mockServerUrl,
    Headers: map[string]string{"x-secret-auth-header": "such-a-necessary-auth-header"},
}

myCustomMockServerClient := NewMockserverClient(customConfig)
```

In most cases you should initialise it with `*environment.Environment`, unless you need to pass custom headers to make the connection possible.

## Typical usages
### Arbitrary mock
You can set any desired behaviour by using `PutExpectations(body interface{}) error` funciton, where the body should be a JSON string conforming to
MockServer's format that consists of a request matcher and corresponding action. For example:
```go
var myk8sEnv *environment.Environment

returnOk := `{
  "httpRequest": {
    "method": "GET",
    "path": "/status"
  },
  "httpResponse": {
    "statusCode": 200,
    "body": {
      "message": "Service is running"
    }
  }
}`

ms := ConnectMockServer(myk8sEnv)
err := ms.PutExpectations(returnOk)
if err != nil {
    panic(err)
}
```

> [!NOTE]
> You can read more about expecations syntax, including OpenAPI v3 or dynamic expecations with JavaScript [here](https://www.mock-server.com/mock_server/creating_expectations.html).

### Returning a random integer
To return random integer in the response body, together with `200` status code for all requests with a given path, use:
```go
err := ms.SetRandomValuePath("/api/v2/my_endpoint")
if err != nil {
    panic(err)
}
```

### Returning a static integer
To return a static integer in the response body, together with `200` status code for all requests with a given path, use:
```go
err := ms.SetValuePath("/api/v2/my_endpoint")
if err != nil {
    panic(err)
}
```

### Returning a static string
To return a static string in the response body, together with `200` status code for all requests with a given path, use:
```go
err := ms.SetStringValuePath("/api/v2/my_endpoint", "oh-my")
if err != nil {
    panic(err)
}
```

### Returning arbitrary data
To return arbitrary data in the response body, together with `200` status code for all requests with a given path, use:
```go
err := ms.SetAnyValueResponse("/api/v2/my_endpoint", []int{1, 2, 3})
if err != nil {
    panic(err)
}
```

### Returning arbitrary response from external adatper
To mock a response for Chainlink's external adapter for all requests with a given path, use:
```go
var mockedResult interface{}
mockedResult = 5
err := ms.SetAnyValuePath("/api/v2/my_endpoint", mockedResult)
if err != nil {
    panic(err)
}
```

It will return a response with following structure:
```json
{
    "id": "",
    "data": {
        "result": 5
    }
}
```

# Troubleshooting
Enabling debug mode for the underlaying HTTP client can be achieved by setting `RESTY_DEBUG` environment variable to `true`.