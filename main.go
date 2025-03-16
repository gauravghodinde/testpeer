package main

import (
	"fmt"
	"log"

	"github.com/schollz/peerdiscovery"
)

func main() {
	// Discover peers with a limit of 1
	discoveries, err := peerdiscovery.Discover(peerdiscovery.Settings{
		Limit: 1, // Stops after finding one peer
	})

	if err != nil {
		log.Fatalf("Error discovering peers: %v", err)
	}

	// Print discovered peers
	for _, d := range discoveries {
		fmt.Printf("Discovered peer: %s\n", d.Address)
	}
}
