package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
)

func main() {
	_, err := fake.NewFakeDataProvider(&fake.Input{Port: 9111})
	err = fake.JSON("GET", "/static-fake", map[string]any{
		"response": "I'm a static fake, put JSON or anything inside me!",
	}, 200)
	err = fake.Func("GET", "/dynamic-fake", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"response": "I'm a dynamic fake, write Go to define my logic!",
		})
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to stark fake server API")
	}
	select {}
}
