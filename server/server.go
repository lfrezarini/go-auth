package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/LucasFrezarini/go-auth-manager/generated"
	"github.com/LucasFrezarini/go-auth-manager/resolvers"
	"github.com/vektah/gqlparser/gqlerror"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(makeExecutableSchema(), makeErrorPresenter()))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func makeExecutableSchema() graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{}})
}

func makeErrorPresenter() handler.Option {
	return handler.ErrorPresenter(
		func(ctx context.Context, e error) *gqlerror.Error {
			if e.Error() == "internal system error" {
				return &gqlerror.Error{
					Message: "Internal server error",
					Path:    graphql.GetResolverContext(ctx).Path(),
					Extensions: map[string]interface{}{
						"code": "INTERNAL_SERVER_ERROR",
					},
				}
			}

			return graphql.DefaultErrorPresenter(ctx, e)
		},
	)
}
