package auth

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var DefaultRequireRole = func(ctx context.Context, allowedRoles []string, methodName string) error {

	userRole, ok := ctx.Value("userRole").(string)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "authentication required: user role not found in context for %s", methodName)
	}

	for _, role := range allowedRoles {
		if userRole == role {
			return nil
		}
	}
	return status.Errorf(codes.PermissionDenied, "access denied: role '%s' not allowed for method '%s'", userRole, methodName)
}

var RequireRole = DefaultRequireRole

func ResetRequireRole() {
	RequireRole = DefaultRequireRole
}
