package soak

import (
	"testing"
)

func TestOCRSoak(t *testing.T) {
	RunSoakTest("@soak-ocr", "chainlink-soak-ocr", 6)
}

func TestKeeperSoak(t *testing.T) {
	RunSoakTest("@soak-keeper-block-time", "chainlink-soak-keeper", 6)
}
