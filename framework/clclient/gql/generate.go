package main

import (
	"fmt"
	"os"

	"github.com/Khan/genqlient/generate"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient/gql/schema"
)

func main() {
	schema := schema.MustGetRootSchema()

	if err := os.WriteFile("./schema/schema.graphql", []byte(schema), 0o600); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	generate.Main()
}
