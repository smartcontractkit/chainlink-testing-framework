package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/pods"
	v1 "k8s.io/api/core/v1"
)

func main() {
	ctx := context.Background()
	if pods.K8sEnabled() {
		_, _, err := pods.Run(ctx, &pods.Config{
			Pods: []*pods.PodConfig{
				{
					Name:     pods.Ptr("fakes-alex"),
					Image:    pods.Ptr(os.Getenv("FAKE_IMAGE")),
					Ports:    []string{"8080:8080"},
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
			panic(fmt.Sprintf("failed to deploy container: %s", err.Error()))
		}
		return
	}

	// Create a default Gin router
	r := gin.Default()

	// ========== STATIC RESPONSE ROUTES ==========

	// Simple static JSON response
	r.GET("/api/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"users": []gin.H{
				{"id": 1, "name": "John Doe", "email": "john@example.com"},
				{"id": 2, "name": "Jane Smith", "email": "jane@example.com"},
			},
			"total":  2,
			"status": "success",
		})
	})

	// Static response with query parameter
	r.GET("/api/user", func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "id parameter is required",
			})
			return
		}

		// Still static but uses the ID in response
		c.JSON(http.StatusOK, gin.H{
			"id":    id,
			"name":  "John Doe",
			"email": "john@example.com",
			"role":  "admin",
		})
	})

	// ========== DYNAMIC RESPONSE ROUTES ==========

	// Dynamic response based on path parameter
	r.GET("/api/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Dynamic logic based on the ID
		var user gin.H
		switch id {
		case "1":
			user = gin.H{"id": 1, "name": "John Doe", "email": "john@example.com"}
		case "2":
			user = gin.H{"id": 2, "name": "Jane Smith", "email": "jane@example.com"}
		default:
			user = gin.H{"id": id, "name": "Unknown User", "email": "unknown@example.com"}
		}

		// Add timestamp to show it's dynamic
		c.JSON(http.StatusOK, gin.H{
			"user":         user,
			"timestamp":    time.Now().Unix(),
			"requested_id": id,
		})
	})

	// Start the server on port 8080
	r.Run(":8080")
}
