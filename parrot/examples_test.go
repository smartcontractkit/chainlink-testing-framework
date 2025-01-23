package parrot_test

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/parrot"
)

func ExampleRegister() {
	// Create a new parrot instance with no logging and a custom save file
	saveFile := "register_example.json"
	p, err := parrot.Wake(parrot.WithLogLevel(zerolog.NoLevel), parrot.WithSaveFile(saveFile))
	if err != nil {
		panic(err)
	}
	defer func() { // Cleanup the parrot instance
		err = p.Shutdown(context.Background()) // Gracefully shutdown the parrot instance
		if err != nil {
			panic(err)
		}
		p.WaitShutdown()    // Wait for the parrot instance to shutdown. Usually unnecessary, but we want to clean up the save file
		os.Remove(saveFile) // Cleanup the save file for the example
	}()

	// Create a new route /test that will return a 200 status code with a text/plain response body of "Squawk"
	route := &parrot.Route{
		Method:             http.MethodGet,
		Path:               "/test",
		RawResponseBody:    "Squawk",
		ResponseStatusCode: http.StatusOK,
	}

	// Register the route with the parrot instance
	err = p.Register(route)
	if err != nil {
		panic(err)
	}

	// Call the route
	resp, err := p.Call(http.MethodGet, "/test")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.StatusCode())
	fmt.Println(string(resp.Body()))
	// Output:
	// 200
	// Squawk
}

func ExampleRoute() {
	// Run the parrot server as a separate instance, like in a Docker container
	saveFile := "route_example.json"
	p, err := parrot.Wake(parrot.WithPort(9090), parrot.WithLogLevel(zerolog.NoLevel), parrot.WithSaveFile(saveFile))
	if err != nil {
		panic(err)
	}
	defer func() { // Cleanup the parrot instance
		err = p.Shutdown(context.Background()) // Gracefully shutdown the parrot instance
		if err != nil {
			panic(err)
		}
		p.WaitShutdown()    // Wait for the parrot instance to shutdown. Usually unnecessary, but we want to clean up the save file
		os.Remove(saveFile) // Cleanup the save file for the example
	}()

	// Code that calls the parrot server from another service
	// Use resty to make HTTP calls to the parrot server
	client := resty.New()

	// Register a new route /test that will return a 200 status code with a text/plain response body of "Squawk"
	route := &parrot.Route{
		Method:             http.MethodGet,
		Path:               "/test",
		RawResponseBody:    "Squawk",
		ResponseStatusCode: http.StatusOK,
	}
	resp, err := client.R().SetBody(route).Post("http://localhost:9090/routes")
	if err != nil {
		panic(err)
	}
	defer resp.RawResponse.Body.Close()
	fmt.Println(resp.StatusCode())

	// Get all routes from the parrot server
	routes := make([]*parrot.Route, 0)
	resp, err = client.R().SetResult(&routes).Get("http://localhost:9090/routes")
	if err != nil {
		panic(err)
	}
	defer resp.RawResponse.Body.Close()
	fmt.Println(resp.StatusCode())
	fmt.Println(len(routes))

	// Delete the route
	req := &parrot.RouteRequest{
		ID: route.ID(),
	}
	resp, err = client.R().SetBody(req).Delete("http://localhost:9090/routes")
	if err != nil {
		panic(err)
	}
	defer resp.RawResponse.Body.Close()
	fmt.Println(resp.StatusCode())

	// Get all routes from the parrot server
	routes = make([]*parrot.Route, 0)
	resp, err = client.R().SetResult(&routes).Get("http://localhost:9090/routes")
	if err != nil {
		panic(err)
	}
	defer resp.RawResponse.Body.Close()
	fmt.Println(len(routes))

	// Output:
	// 201
	// 200
	// 1
	// 204
	// 0
}

func ExampleRecorder() {
	saveFile := "recorder_example.json"
	p, err := parrot.Wake(parrot.WithLogLevel(zerolog.NoLevel), parrot.WithSaveFile(saveFile))
	if err != nil {
		panic(err)
	}
	defer func() { // Cleanup the parrot instance
		err = p.Shutdown(context.Background()) // Gracefully shutdown the parrot instance
		if err != nil {
			panic(err)
		}
		p.WaitShutdown()    // Wait for the parrot instance to shutdown. Usually unnecessary, but we want to clean up the save file
		os.Remove(saveFile) // Cleanup the save file for the example
	}()

	// Create a new recorder
	recorder, err := parrot.NewRecorder()
	if err != nil {
		panic(err)
	}

	// Register the recorder with the parrot instance
	err = p.Record(recorder.URL)
	if err != nil {
		panic(err)
	}
	defer recorder.Close()

	// Register a new route /test that will return a 200 status code with a text/plain response body of "Squawk"
	route := &parrot.Route{
		Method:             http.MethodGet,
		Path:               "/test",
		RawResponseBody:    "Squawk",
		ResponseStatusCode: http.StatusOK,
	}
	err = p.Register(route)
	if err != nil {
		panic(err)
	}

	// Call the route
	go func() {
		_, err := p.Call(http.MethodGet, "/test")
		if err != nil {
			panic(err)
		}
	}()

	// Record the route call
	for {
		select {
		case recordedRouteCall := <-recorder.Record():
			if recordedRouteCall.RouteID == route.ID() {
				fmt.Println(recordedRouteCall.RouteID)
				fmt.Println(recordedRouteCall.Request.Method)
				fmt.Println(recordedRouteCall.Response.StatusCode)
				fmt.Println(string(recordedRouteCall.Response.Body))
				return
			}
		case err := <-recorder.Err():
			panic(err)
		}
	}
	// Output:
	// GET:/test
	// GET
	// 200
	// Squawk
}
