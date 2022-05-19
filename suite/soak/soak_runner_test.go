package soak

import (
	"testing"
)

func TestOCRSoak(t *testing.T) {
	RunSoakTest(t, "@soak-ocr", "chainlink-soak-ocr", 6, nil)
}

func TestKeeperSoak(t *testing.T) {
	RunSoakTest(t, "@soak-keeper-block-time", "chainlink-soak-keeper", 6, nil)
}
