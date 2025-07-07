package main

import (
	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
)

func main() {
	_, err := fake.NewFakeDataProvider(nil)
	err = fake.JSON("GET", "/static-fake", map[string]any{
		"response": "I'm a static fake, put JSON or anything inside me!",
	}, 200)
	err = fake.Func("GET", "/dynamic-fake", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"response": "I'm a dynamic fake, write Go to define my logic!",
		})
	})
	if err != nil {
		panic(err)
	}
	select {}
}
