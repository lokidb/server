package main

import (
	"github.com/hvuhsg/lokidb/cli"
	client "github.com/hvuhsg/lokidb/clients/go"
)

func main() {
	client := client.New()
	defer client.Close()

	cli.ShellLoop(client)
}
