package resolvers

import (
	"context"

	"github.com/LucasFrezarini/go-auth-manager/jsonwebtoken"
)

type claimsResolver struct{ *Resolver }

func (r *claimsResolver) Iss(ctx context.Context, obj *jsonwebtoken.Claims) (string, error) {
	return obj.Issuer, nil
}

func (r *claimsResolver) Sub(ctx context.Context, obj *jsonwebtoken.Claims) (string, error) {
	return obj.Subject, nil
}

func (r *claimsResolver) Exp(ctx context.Context, obj *jsonwebtoken.Claims) (int, error) {
	return int(obj.ExpiresAt), nil
}

func (r *claimsResolver) Iat(ctx context.Context, obj *jsonwebtoken.Claims) (int, error) {
	return int(obj.IssuedAt), nil
}
