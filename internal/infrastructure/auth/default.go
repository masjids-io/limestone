package auth

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	// Assuming you have a password hashing function here, if not, remove or add it
	// "golang.org/x/crypto/bcrypt"
)

// DefaultRequireRole is the actual, production-ready implementation for role-based authorization.
// It takes context, a slice of allowed roles, and the method name (for logging/error messages).
var DefaultRequireRole = func(ctx context.Context, allowedRoles []string, methodName string) error {
	// In a real application, you would extract the user's role from the context.
	// This context might have been populated by an interceptor or middleware.
	// For demonstration, let's assume a "userRole" key in context.

	userRole, ok := ctx.Value("userRole").(string)
	if !ok {
		// If no role is found in context, consider it unauthenticated or unauthorized.
		return status.Errorf(codes.Unauthenticated, "authentication required: user role not found in context for %s", methodName)
	}

	// Check if the user's role is among the allowed roles
	for _, role := range allowedRoles {
		if userRole == role {
			return nil // Role is allowed, no error
		}
	}

	// If the loop finishes, the user's role is not allowed
	return status.Errorf(codes.PermissionDenied, "access denied: role '%s' not allowed for method '%s'", userRole, methodName)
}

// RequireRole is the variable that can be reassigned for testing purposes.
// It's initialized with the DefaultRequireRole.
var RequireRole = DefaultRequireRole

// ResetRequireRole resets the RequireRole variable back to its DefaultRequireRole implementation.
// This is crucial for test isolation, ensuring tests don't affect each other's mocks.
func ResetRequireRole() {
	RequireRole = DefaultRequireRole
}

// You might also have a password hashing utility function in this package
// func HashPassword(password string) (string, error) {
//     bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//     return string(bytes), err
// }

// func CheckPasswordHash(password, hash string) bool {
//     err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
//     return err == nil
// }
