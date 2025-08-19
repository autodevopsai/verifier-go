package main

import (
	"fmt"
	"github.com/autodevops/verifier-go/internal/provider"
)

func main() {
	p := provider.NewAnthropicProvider("test-api-key", "claude-2")
	fmt.Println(p)
}
