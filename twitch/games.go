package twitch

import "log"
import "encoding/json"
import "io/ioutil"
import "fmt"
import "net/http"
import "strconv"

type GamesRequest struct {
	Games      []Game `json:"top"`
	TotalGames uint32 `json:"_total"`
}

type Game struct {
	Info     GameInfo `json:"game"`
	Viewers  uint64   `json:"viewers"`
	Channels uint64   `json:"channels"`
}

type GameInfo struct {
	Name string `json:"name"`
	Id   uint64 `json:"_id"`
}

func GetGames(limit uint8, offset uint32) (games []Game, total uint32) {
	client := &http.Client{}
	topGamesUrl := TwitchApiBaseUrl + "games/top/"
	req, err := http.NewRequest("GET", topGamesUrl, nil)
	if err != nil {
		log.Println("Couldn't create request")
	}
	req.Header.Add("Accept", "application/vnd.twitchtv.v3+json")
	req.Header.Add("Client-ID", TshClientId)
	q := req.URL.Query()
	q.Add("offset", strconv.FormatUint(uint64(offset), 10))
	q.Add("limit", strconv.FormatUint(uint64(limit), 10))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Couldn't perform HTTP request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Couldn't read message body")
	}

	request, err := unmarshalGames(body)
	if err != nil {
		panic("Error unmarshaling JSON")
	}
	return request.Games, request.TotalGames
}

func unmarshalGames(body []byte) (*GamesRequest, error) {
	var req = new(GamesRequest)
	err := json.Unmarshal(body, &req)
	return req, err
}

func (request GamesRequest) String() string {
	return fmt.Sprintf("Request: %s\n", request.Games)
}

func GamesToStrings(games []Game) []string {
	toStrings := make([]string, len(games), len(games))
	for i := range games {
		toStrings[i] = games[i].String()
	}
	return toStrings
}

func (game Game) String() string {
	var toString string
	toString = fmt.Sprintf("[%s](fg-blue), %d Viewers, %d Streamers\n", game.Info,
		game.Viewers, game.Channels)
	return toString
}

func (info GameInfo) String() string {
	var toString string
	toString = fmt.Sprintf("%s", info.Name)
	return toString
}
