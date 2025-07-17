package main

import (
	"context"
	"github.com/marcoshuck/mcp-demo/disher"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"log"
)

func main() {
	// Business logic
	ds := disher.NewDishServer(disher.DefaultDishes)
	// Transport layer
	server := mcp.NewServer("disher", "v1.0.0", nil)
	// Register endpoints
	disher.RegisterTools(server, ds)
	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Fatal(err)
	}
}
