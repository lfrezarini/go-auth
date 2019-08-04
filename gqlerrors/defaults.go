package gqlerrors

import "github.com/vektah/gqlparser/gqlerror"

// CreateInternalServerError creates a default GraphQL internal server error from a request
func CreateInternalServerError(message string) *gqlerror.Error {
	return &gqlerror.Error{
		Message: message,
		Extensions: map[string]interface{}{
			"code": Internal,
		},
	}
}

// CreateAuthorizationError creates a default GraphQL authorization error from a request
func CreateAuthorizationError() *gqlerror.Error {
	return &gqlerror.Error{
		Message: "Unauthorized",
		Extensions: map[string]interface{}{
			"code": Unauthorized,
		},
	}
}

// CreateConflictError creates a default Graphql conflict error from a request
func CreateConflictError(message string) *gqlerror.Error {
	return &gqlerror.Error{
		Message: message,
		Extensions: map[string]interface{}{
			"code": Conflict,
		},
	}
}
