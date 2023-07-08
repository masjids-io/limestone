package auth


import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

type Session struct {
	Data struct {
		Name      string `json:"name"` // unique name of the org 
		ID        int    `json:"id"`    // org id
		FirstName string `json:"firstName"` // user first name
		LastName  string `json:"lastName"`	// user last name
		UserID    int    `json:"userId"`	// user id
	} `json:"data"`
}

//update .env file with REDWOOD_URL=http://dev.mosque.icu:8910


func CheckAuth(req *http.Request) (*Session, error) {

	var id string = req.URL.Query().Get("id") //public user id passed in as query param

	cookie := req.Header.Get("Cookie")
	url := os.Getenv("REDWOOD_URL") + "/api/upload?id=" + id
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", cookie)
	req.Header.Set("auth-provider", "dbAuth")
	req.Header.Set("Authorization", "Bearer "+ id)

	client := &http.Client{}
	response, err := client.Do(req)
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

	return &session, nil
}
