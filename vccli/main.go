package main

import (
	"log"

	"github.com/youngjoon-lee/panacea-vc-cli/vccli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
