package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/pods"
	v1 "k8s.io/api/core/v1"
)

var users = map[string]gin.H{
	"1": {"id": 1, "name": "John Doe", "email": "john@example.com"},
	"2": {"id": 2, "name": "Jane Smith", "email": "jane@example.com"},
}

func main() {
	if pods.K8sEnabled() {
		deployToK8s()
		return
	}
	r := gin.Default()
	registerRoutes(r)
	log.Fatal(r.Run(":80"))
}

func deployToK8s() {
	_, _, err := pods.Run(context.Background(), &pods.Config{
		Pods: []*pods.PodConfig{
			{
				Name:     pods.Ptr("fakes-alex"),
				Image:    pods.Ptr(os.Getenv("FAKE_IMAGE")),
				Ports:    []string{"80:80"},
				Requests: pods.ResourcesMedium(),
				Limits:   pods.ResourcesMedium(),
				ContainerSecurityContext: &v1.SecurityContext{
					RunAsUser:  pods.Ptr[int64](999),
					RunAsGroup: pods.Ptr[int64](999),
				},
			},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("failed to deploy container: %s", err))
	}
}

func registerRoutes(r *gin.Engine) {
	// health check
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	webhook := r.Group("/api/webhook")
	webhook.POST("", func(c *gin.Context) { c.Status(http.StatusOK) })
	webhook.GET("/users", handleGetUsers)
	webhook.GET("/user", handleGetUser)
	webhook.GET("/users/:id", handleGetUserByID)
}

func lookupUser(id string) gin.H {
	if user, ok := users[id]; ok {
		return user
	}
	return gin.H{"id": id, "name": "Unknown User", "email": "unknown@example.com"}
}

func handleGetUsers(c *gin.Context) {
	allUsers := make([]gin.H, 0, len(users))
	for _, u := range users {
		allUsers = append(allUsers, u)
	}
	c.JSON(http.StatusOK, gin.H{
		"users":  allUsers,
		"total":  len(users),
		"status": "success",
	})
}

func handleGetUser(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}
	c.JSON(http.StatusOK, lookupUser(id))
}

func handleGetUserByID(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"user":         lookupUser(id),
		"timestamp":    time.Now().Unix(),
		"requested_id": id,
	})
}
