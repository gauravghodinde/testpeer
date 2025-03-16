package main

import (
	"flag"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Get local IP
func getLocalIP() string {
	cmd := exec.Command("hostname", "-I") // Works on Linux
	output, _ := cmd.Output()
	localIP := strings.Fields(string(output))
	if len(localIP) > 0 {
		return localIP[0]
	}
	return ""
}

// Get peer IPs from arp -a
func getPeers(myIP string) []string {
	cmd := exec.Command("arp", "-a")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("❌ Error running arp:", err)
		return nil
	}

	re := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	ips := re.FindAllString(string(output), -1)

	var peers []string
	for _, ip := range ips {
		if ip != myIP { // Ignore self
			peers = append(peers, ip)
		}
	}
	return peers
}

// Send UDP discovery message
func sendDiscovery(peers []string) {
	message := []byte("Hello, peer!") // Custom discovery message

	for _, ip := range peers {
		addr := fmt.Sprintf("%s:9999", ip) // Use port 9999
		conn, err := net.Dial("udp", addr)
		if err != nil {
			fmt.Println("⚠️ Could not reach", ip)
			continue
		}
		defer conn.Close()

		conn.Write(message)
		fmt.Println("📤 Sent discovery to", ip)
		time.Sleep(100 * time.Millisecond) // Avoid flooding
	}
}

// Listen for UDP messages
func receiveDiscovery() {
	addr, err := net.ResolveUDPAddr("udp", ":9999")
	if err != nil {
		fmt.Println("❌ Error resolving address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("❌ Error listening:", err)
		return
	}
	defer conn.Close()

	fmt.Println("📡 Listening for peers on port 9999...")

	buffer := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("⚠️ Error reading:", err)
			continue
		}
		fmt.Printf("✅ Received from %s: %s\n", remoteAddr, string(buffer[:n]))
	}
}

func main() {
	// Define command-line flags
	sendFlag := flag.Bool("send", false, "Send UDP discovery message")
	receiveFlag := flag.Bool("receive", false, "Listen for UDP discovery messages")
	flag.Parse()

	if *sendFlag {
		myIP := getLocalIP()
		fmt.Println("🚀 Your IP:", myIP)

		peers := getPeers(myIP)
		if len(peers) == 0 {
			fmt.Println("⚠️ No peers found.")
			return
		}

		fmt.Println("🔍 Peers found:", peers)
		sendDiscovery(peers)
	} else if *receiveFlag {
		receiveDiscovery()
	} else {
		fmt.Println("Usage: go run main.go -send   # To send discovery")
		fmt.Println("       go run main.go -receive # To receive messages")
	}
}
