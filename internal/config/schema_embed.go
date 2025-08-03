package config

import (
	_ "embed"
)

//go:embed schema.json
var embeddedSchema []byte
