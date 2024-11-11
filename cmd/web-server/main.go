package main

import (
	"flag"

	"github.com/open-wm/blockehr/pkg/web"
)

// Sample:
// go run -v main.go -port :3000
func main() {
	// Parse command line arguments
	port := "3000"
	flag.StringVar(&port, "port", port, "Port to listen on")
	flag.Parse()

	web.StartWebServer(port)
}
