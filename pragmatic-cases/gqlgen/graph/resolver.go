package graph

//go:generate go run github.com/99designs/gqlgen generate

import "tmp/pragmatic-cases/gqlgen/graph/model"

type Resolver struct {
	todos []*model.Todo
}
