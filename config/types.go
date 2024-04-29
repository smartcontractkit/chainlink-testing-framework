package config

import "github.com/smartcontractkit/seth"

type SethConfig interface {
	GetSethConfig() *seth.Config
}
