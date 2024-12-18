# Grafana

> [!NOTE]
> Contrary to other clients, you will find Grafana client in a [separate package & go module](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/lib/grafana).

Grafana client encapsulate following functionalities:
* Dashboard creation
* Managing dashboard annotations (CRUD)
* Checking alerts

# New instance
In order to create a new instance you will need:
* URL
* API token

For example:
```go
url := "http://grafana.io"
apiToken := "such-a-secret-1&11n"
gc := NewGrafanaClient(url, apiToken)
```

# Dashboard creation
You can create a new dashboard defined in JSON with:
```go

//define your dashboard here
dashboardJson := ``

request := PostDashboardRequest {
    Dashboard: dashboardJson,
    FolderId: 5 // change to your folder id
}

dr, rawResponse, err := gc.PostDashboard(request)
if err != nil {
    panic(err)
}

if rawResponse.StatusCode() != 200 {
    panic("response code wasn't 200, but " + rawResponse.StatusCode())
}

fmt.Println("Dashboard slug is is " + *dr.Slug)
```

# Posting annotations
You can post annotations in a following way:
```go
annDetails := PostAnnotation {
	DashboardUID: "some-uid",
	Time:         time.Now(),
	TimeEnd:      time.Now().Add(1 * time.Second)
	Text:         "my test annotation"
}

r, rawResponse, err := gc.PostAnnotation(annDetails)
if rawResponse.StatusCode() != 200 {
    panic("response code wasn't 200, but " + rawResponse.StatusCode())
}

fmt.Println("Created annotation with id: " + r.Id)

```

# Checking alerts
You can check alerts firing for a dashboard with UID:
```go
alerts, rawResponse, err := gc.AlertRulerClient.GetAlertsForDashboard("some-uid")
if rawResponse.StatusCode() != 200 {
    panic("response code wasn't 200, but " + rawResponse.StatusCode())
}

for name, value := range alerts {
    fmt.Println("Alert named " + name + "was triggered. Details: " + string(value))
}
```

# Troubleshooting
To enable debug mode for the underlaying HTTP client set `RESTY_DEBUG` environment variable to `true`.