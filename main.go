package main

import (
	"github.com/NitroSniper/indigo/server"
)

func main() {
	config := server.ServerConfig{Name: "./example.md"}

	config.HostServer()
}
