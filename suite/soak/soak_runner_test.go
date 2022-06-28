package soak_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/actions"
)

func TestOCRSoak(t *testing.T) {
	actions.RunSoakTest(t, "@soak-ocr", "chainlink-soak-ocr", map[string]interface{}{
		"replicas": 6,
	})
}

func TestKeeperSoak(t *testing.T) {
	actions.RunSoakTest(t, "@soak-keeper", "soak-keeper", map[string]interface{}{
		"replicas": 6,
	})
}
