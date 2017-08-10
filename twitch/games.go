package twitch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type gamesRequest struct {
	Games      []Game `json:"top"`
	TotalGames uint32 `json:"_total"`
}

// Game represents information regarding a Game
type Game struct {
	Info     GameInfo `json:"game"`
	Viewers  uint64   `json:"viewers"`
	Channels uint64   `json:"channels"`
}

// GameInfo contains the name and unique ID of a game
type GameInfo struct {
	Name string `json:"name"`
	ID   uint64 `json:"_id"`
}

// GetGames retrieves a number of games limited to the specified amount
// starting from the given offset of top results
func GetGames(limit uint8, offset uint32) (games []Game, total uint32) {
	params := map[string]string{
		"offset": strconv.FormatUint(uint64(offset), 10),
		"limit":  strconv.FormatUint(uint64(limit), 10),
	}
	body, err := SendRequest(http.MethodGet, "games/top", params)

	request, err := unmarshalGames(body)
	if err != nil {
		panic("Error unmarshaling JSON")
	}
	return request.Games, request.TotalGames
}

func unmarshalGames(body []byte) (*gamesRequest, error) {
	var req = new(gamesRequest)
	err := json.Unmarshal(body, &req)
	return req, err
}

func (request gamesRequest) String() string {
	return fmt.Sprintf("Request: %s\n", request.Games)
}

// GamesToStrings converts a slice of Games into a slice of their string representations
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
