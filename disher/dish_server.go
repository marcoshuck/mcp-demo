package disher

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Dish struct {
	Name        string
	DisplayName string
}

type ServeDishRequest struct {
	Name string
}

type ServeDishResponse struct {
	Dish Dish
}

type ListDishesRequest struct {
	Limit int
}

type ListDishesResponse struct {
	Dishes []Dish
}

type DishServer interface {
	ServeDish(ctx context.Context, req *ServeDishRequest) (*ServeDishResponse, error)
	ListDishes(ctx context.Context, req *ListDishesRequest) (*ListDishesResponse, error)
}

type dishServer struct {
	// dishes contains a map of dish name to dish display name.
	dishes map[string]string
}

func (svc *dishServer) ServeDish(ctx context.Context, req *ServeDishRequest) (*ServeDishResponse, error) {
	displayName, ok := svc.dishes[req.Name]
	if !ok {
		return nil, errors.New("not found")
	}
	return &ServeDishResponse{
		Dish: Dish{
			Name:        req.Name,
			DisplayName: displayName,
		},
	}, nil
}

func (svc *dishServer) ListDishes(ctx context.Context, req *ListDishesRequest) (*ListDishesResponse, error) {
	var res ListDishesResponse
	res.Dishes = make([]Dish, 0, req.Limit)
	var count int
	for k, v := range svc.dishes {
		if count == req.Limit {
			break
		}
		res.Dishes = append(res.Dishes, Dish{
			Name:        k,
			DisplayName: v,
		})
		count++
	}
	return &res, nil
}

func NewDishServer(dishes map[string]string) DishServer {
	return &dishServer{
		dishes: dishes,
	}
}

func RegisterTools(s *mcp.Server, ds DishServer) {
	s.AddTools(
		mcp.NewServerTool[ListDishesRequest, ListDishesResponse](
			"list-dishes",
			"List the available dishes at the restaurant",
			ListDishesHandler(ds),
			mcp.Input(
				mcp.Property(
					"limit",
					mcp.Description("the number of items to return when listing the available dishes"),
				),
			),
		),
		mcp.NewServerTool[ServeDishRequest, ServeDishResponse](
			"serve-dish",
			"Serve a dish for a customer to eat at the restaurant",
			ServeDishHandler(ds),
			mcp.Input(
				mcp.Property(
					"name",
					mcp.Description("the name of the dish to serve at the restaurant"),
				),
			),
		),
	)
}

func ListDishesHandler(ds DishServer) mcp.ToolHandlerFor[ListDishesRequest, ListDishesResponse] {
	return func(ctx context.Context, session *mcp.ServerSession, c *mcp.CallToolParamsFor[ListDishesRequest]) (*mcp.CallToolResultFor[ListDishesResponse], error) {
		res, err := ds.ListDishes(ctx, &ListDishesRequest{Limit: c.Arguments.Limit})
		if err != nil {
			return nil, err
		}
		b, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}
		return &mcp.CallToolResultFor[ListDishesResponse]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(b)},
			},
			StructuredContent: *res,
		}, nil
	}
}

func ServeDishHandler(ds DishServer) mcp.ToolHandlerFor[ServeDishRequest, ServeDishResponse] {
	return func(ctx context.Context, session *mcp.ServerSession, c *mcp.CallToolParamsFor[ServeDishRequest]) (*mcp.CallToolResultFor[ServeDishResponse], error) {
		res, err := ds.ServeDish(ctx, &ServeDishRequest{Name: c.Arguments.Name})
		if err != nil {
			return nil, err
		}
		b, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}
		return &mcp.CallToolResultFor[ServeDishResponse]{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(b)},
			},
			StructuredContent: *res,
		}, nil
	}
}
