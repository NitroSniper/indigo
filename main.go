package main

import (
	"github.com/NitroSniper/indigo/server"
)

func main() {
	config := server.NewMarkdownServer()

	config.HostServer()
}
