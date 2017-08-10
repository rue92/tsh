package twitch

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const twitchAPIBaseURL = "https://api.twitch.tv/kraken/"
const tshClientID = "grjvdyr25euc3o1r3peeplpqed75icb"
const twitchAPIVersionToken = "application/vnd.twitchtv.v5+json"

type tokenContainer struct {
	Token token `json:"token"`
}

type token struct {
	Valid    bool   `json:"valid"`
	UserName string `json:"user_name"`
	UserID   string `json:"user_id"`
}

var twitchClient = &http.Client{
	Timeout: time.Second * 10,
}

func prepareRequest(endpoint string, method string) (*http.Request, error) {
	requestURL := twitchAPIBaseURL + endpoint
	req, err := http.NewRequest(method, requestURL, nil)
	if err != nil {
		log.Println("Couldn't create request")
	}
	req.Header.Add("Accept", twitchAPIVersionToken)
	req.Header.Add("Client-ID", tshClientID)
	return req, err
}

// SendRequest Sends the request for using the given method invoked on the given endpoint
func SendRequest(method string, endpoint string, queryParams map[string]string) ([]byte, error) {
	req, _ := prepareRequest(endpoint, method)
	q := req.URL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := twitchClient.Do(req)
	if err != nil {
		log.Println("Couldn't perform HTTP request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Couldn't read message body")
	}

	return body, err
}

// Authenticate communicates with the REST endpoint to retrieve an access token for this
// session corresponding to the given OAuth token.
func Authenticate(oauth string) bool {
	client := &http.Client{}
	authRequest, _ := prepareRequest("", "GET")
	authRequest.Header.Add("Authorization", "OAuth "+oauth)

	resp, err := client.Do(authRequest)
	if err != nil {
		log.Println("Couldn't perform HTTP request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Couldn't read message body")
	}

	request, err := unmarshalAuthorization(body)
	if err != nil {
		panic("Error unmarshaling JSON")
	}
	return request.Token.Valid
}

func unmarshalAuthorization(body []byte) (*tokenContainer, error) {
	var req = new(tokenContainer)
	err := json.Unmarshal(body, &req)
	return req, err
}
