package auth

type UnprotectedRoute struct {
	Path   string
	Method string
}

var UnprotectedRoutes = map[string]bool{
	"/limestone.UserService/CreateUser":       true,
	"/limestone.UserService/AuthenticateUser": true,
	"/limestone.UserService/RefreshToken":     true,
}

var UnprotectedRoutesHTTP = map[UnprotectedRoute]bool{
	{Path: "/v1/users", Method: "POST"}:              true,
	{Path: "/v1/auth/login", Method: "POST"}:         true,
	{Path: "/v1/auth/refresh_token", Method: "POST"}: true,
	{Path: "/v1/masjids", Method: "GET"}:             true,
}
