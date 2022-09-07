//go:build tools
// +build tools

package tools

// nolint: typecheck
import (
	_ "entgo.io/ent/cmd/ent"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
