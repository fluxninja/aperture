//go:generate go run github.com/99designs/gqlgen generate

package graph

import "github.com/fluxninja/aperture/playground/graphql_demo_app/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver tracks the state of todos.
type Resolver struct {
	todos []*model.Todo
}
