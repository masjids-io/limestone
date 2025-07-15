package auth

import (
	"context"
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

//func RequireRole(ctx context.Context, allowedRoles []string, operationName string) error {
//	userRole, ok := ctx.Value(UserRoleContextKey).(string)
//	fmt.Println(userRole)
//	fmt.Println(ok)
//	if !ok {
//		return status.Errorf(codes.Unauthenticated, "Authentication context missing for operation '%s'", operationName)
//	}
//
//	if !containsRole(allowedRoles, userRole) {
//		return status.Errorf(codes.PermissionDenied, "User with role '%s' does not have permission for operation '%s'. Required roles: %v", userRole, operationName, allowedRoles)
//	}
//
//	return nil
//}

func RequireSelfOrRole(ctx context.Context, targetUserID string, allowedRoles []string, operationName string) error {
	requesterID, idOk := ctx.Value(UserIDContextKey).(string)
	requesterRole, roleOk := ctx.Value(UserRoleContextKey).(string)

	if !idOk || !roleOk {
		return status.Errorf(codes.Unauthenticated, "Requester information missing from context for operation '%s'.", operationName)
	}

	if requesterID == targetUserID {
		return nil
	}

	if containsRole(allowedRoles, requesterRole) {
		return nil
	}

	return status.Errorf(codes.PermissionDenied, "User with role '%s' does not have permission for operation '%s'. Requires self-access or one of roles: %v", requesterRole, operationName, allowedRoles)
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(string)
	return userID, ok
}

func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	userRole, ok := ctx.Value(UserRoleContextKey).(string)
	return userRole, ok
}
