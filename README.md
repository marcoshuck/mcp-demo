# Go MCP Server

This repository contains a sample Go project demonstrating how to build a server using the Model Context
Protocol (MCP) Go SDK.

The project implements a simple `Disher` service, which simulates a restaurant's dish-serving capabilities. It's
designed as an educational example to showcase the core patterns for exposing business logic as tools that a large
language model (LLM) can interact with via MCP.

# Core concepts

## Business logic

The core logic is defined in the `DishServer` interface and its implementation. This part of the code is
standard Go and knows nothing about MCP. The interface requires two methods: `ServeDish` to handle serving a single
dish and `ListDishes` to get a list of available dishes. This keeps the primary business functionality clean and easy to
test independently.

```go
type DishServer interface {
    ServeDish(ctx context.Context, req *ServeDishRequest) (*ServeDishResponse, error)
    ListDishes(ctx context.Context, req *ListDishesRequest) (*ListDishesResponse, error)
}
```

## MCP endpoint registration

The `RegisterTools` function acts as the bridge to MCP. It takes our `DishServer` implementation and registers its
capabilities on an MCP server. It uses the `mcp.NewServerTool` function to define two distinct tools: `list-dishes` and
`serve-dish`. Each tool is given a name, a human-readable description for the LLM, and a specific input schema, making the
service's functions discoverable and usable by an LLM model.

## Handlers

The `ServeDishHandler` and `ListDishesHandler` functions are adapters. Their job is to translate an incoming tool call
from the MCP server into a standard method call on our DishServer instance. They then format the response from our
service back into the structure that the MCP server expects, completing the request-response cycle.