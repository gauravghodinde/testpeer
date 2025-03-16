package main

import (
	"fmt"
	"log"
	"time"

	"github.com/schollz/peerdiscovery"
)

func main() {
	fmt.Println("🔍 Starting peer discovery...")

	discoveries, err := peerdiscovery.Discover(peerdiscovery.Settings{
		Limit:     -1,               // Find unlimited peers
		TimeLimit: 30 * time.Second, // Run for 30 seconds
		Delay:     2000,             // Wait 2 seconds before retrying
		AllowSelf: false,            // Allow self-discovery (for debugging)
	})

	if err != nil {
		log.Fatalf("❌ Error discovering peers: %v", err)
	}

	// Check if no peers were found
	if len(discoveries) == 0 {
		fmt.Println("⚠️ No peers found. Check firewall, AP Isolation, and network settings.")
	} else {
		// Print discovered peers
		for _, d := range discoveries {
			fmt.Printf("✅ Discovered peer: %s\n", d.Address)
		}
	}
}
