package auth

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
)

type Session struct {
	Data struct {
		Name      string `json:"name"`      // unique name of the org
		ID        int    `json:"id"`        // org id
		FirstName string `json:"firstName"` // user first name
		LastName  string `json:"lastName"`  // user last name
		UserID    int    `json:"userId"`    // user id
	} `json:"data"`
}

//update .env file with REDWOOD_URL=http://dev.mosque.icu:8910

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	id, err := grpc_auth.AuthFromMD(ctx, "id")
	if err != nil {
		return nil, err
	}
	cookie, err := grpc_auth.AuthFromMD(ctx, "Cookie")
	if err != nil {
		return nil, err
	}

	url := os.Getenv("REDWOOD_URL") + "/api/upload?id=" + id
	httpReq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Cookie", cookie)
	httpReq.Header.Set("auth-provider", "dbAuth")
	httpReq.Header.Set("Authorization", "Bearer "+id)

	client := &http.Client{}
	response, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var session Session
	err = json.Unmarshal(bodyBytes, &session)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)
}
