package middlewares

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/LucasFrezarini/go-auth-manager/generated"
	"github.com/LucasFrezarini/go-auth-manager/gqlerrors"
	"github.com/LucasFrezarini/go-auth-manager/resolvers"
	"github.com/vektah/gqlparser/gqlerror"
)

// MakeHandlers returns the handlers used by server
func MakeHandlers() http.Handler {
	return AuthHandler(handler.GraphQL(makeExecutableSchema(), makeErrorPresenter()))
}

func makeExecutableSchema() graphql.ExecutableSchema {
	c := generated.Config{Resolvers: &resolvers.Resolver{}}
	c.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		userID := ctx.Value("userID")
		if userID != nil {
			return next(ctx)
		}

		return nil, gqlerrors.CreateAuthorizationError()
	}

	return generated.NewExecutableSchema(c)
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
