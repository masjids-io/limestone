package auth

var UnprotectedRoutes = map[string]bool{
	"/limestone.UserService/CreateUser":       true,
	"/limestone.UserService/AuthenticateUser": true,
	"/limestone.UserService/RefreshToken":     true,
}

var UnprotectedRoutesHTTP = map[string]bool{
	"/users/authenticate":                     true,
	"/limestone.UserService/AuthenticateUser": true,
	"/limestone.UserService/RefreshToken":     true,
}
