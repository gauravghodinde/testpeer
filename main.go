package main

import (
	"fmt"
	"net"
)

func main() {
	interfaces, _ := net.Interfaces()

	for _, i := range interfaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			fmt.Println("🔍 Found device IP:", addr.String())
		}
	}
}
