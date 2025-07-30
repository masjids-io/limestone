package auth

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func containsRole(s []string, role string) bool {
	for _, r := range s {
		if r == role {
			return true
		}
	}
	return false
}

func RequireRole(ctx context.Context, allowedRoles []string, operationName string) error {
	userRole, ok := ctx.Value(UserRoleContextKey).(string)
	fmt.Println(userRole)
	fmt.Println(ok)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "role not allowed for %s", operationName)
	}

	if !containsRole(allowedRoles, userRole) {
		return status.Errorf(codes.Unauthenticated, "authentication required: user role not found in context for %s", operationName)
	}

	return nil
}

var TestingRequireRole = func(ctx context.Context, allowedRoles []string, methodName string) error {

	userRole, ok := ctx.Value("userRole").(string)
	fmt.Println(userRole)
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

var RequireRoleTest = TestingRequireRole

func ResetRequireRole() {
	RequireRoleTest = TestingRequireRole
}
