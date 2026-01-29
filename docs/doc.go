package main

import (
	"log"

	"github.com/nikhilsbhat/ingress-traefik-converter/cmd"
	"github.com/spf13/cobra/doc"
)

//go:generate go run github.com/nikhilsbhat/ingress-traefik-converter/docs
func main() {
	commands := cmd.SetIngressTraefikConverterCommands()

	if err := doc.GenMarkdownTree(commands, "doc"); err != nil {
		log.Fatal(err)
	}
}
