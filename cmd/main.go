package main

import (
	"context"
	"github.com/marcoshuck/mcp-demo/disher"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"log"
)

func main() {
	ds := disher.NewDishServer(disher.DefaultDishes)
	server := mcp.NewServer("disher", "v1.0.0", nil)
	disher.RegisterTools(server, ds)
	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Fatal(err)
	}
}
