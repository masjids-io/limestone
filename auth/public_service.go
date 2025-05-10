package auth

var UnprotectedRoutes = map[string]bool{
	"/limestone.UserService/CreateUser":       true,
	"/limestone.UserService/AuthenticateUser": true,
	"/limestone.UserService/RefreshToken":     true,
}
