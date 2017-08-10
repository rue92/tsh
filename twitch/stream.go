package twitch

import "time"
import "fmt"

import "strconv"
import "encoding/json"
import "net/http"

type streamRequest struct {
	Streams []Stream `json:"streams"`
	Total   uint32   `json:"_total"`
}

// Channel represents all of the information relevant to a specific channel
type Channel struct {
	Mature              bool      `json:"mature"`
	Status              string    `json:"status"`
	BroadcasterLanguage string    `json:"broadcaster_language"`
	DisplayName         string    `json:"display_name"`
	Game                string    `json:"game"`
	Delay               float64   `json:"delay"`
	Language            string    `json:"language"`
	ID                  uint64    `json:"_id"`
	Name                string    `json:"name"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	Partner             bool      `json:"partner"`
	URL                 string    `json:"url"`
}

// Stream represents all of the information relevant to a specific stream,
// itself representing the instantaneous activity of a channel.
type Stream struct {
	Game        string    `json:"game"`
	Viewers     uint64    `json:"viewers"`
	AverageFps  float64   `json:"average_fps"`
	Delay       float64   `json:"delay"`
	VideoHeight uint64    `json:"video_height"`
	IsPlaylist  bool      `json:"is_playlist"`
	CreatedAt   time.Time `json:"created_at"`
	ID          uint64    `json:"_id"`
	Channel     Channel   `json:"channel"`
}

// GetStream retrieves the stream information relevant to a specific user
func GetStream(user string) (stream Stream) {
	params := map[string]string{}
	body, err := SendRequest(http.MethodGet, "streams/", params)

	request, err := getRequestFromJSON(body)
	if err != nil {
		panic("Error unmarshaling JSON")
	}
	return request.Streams[0]
}

// GetStreams retrieves the given limit number of streams from the given
// offset of the top stream sorted by number of viewers
func GetStreams(limit uint8, offset uint32) (streams []Stream, total uint32) {
	params := map[string]string{
		"offset": strconv.FormatUint(uint64(offset), 10),
		"limit":  strconv.FormatUint(uint64(limit), 10),
	}
	body, err := SendRequest(http.MethodGet, "streams/", params)

	request, err := getRequestFromJSON(body)
	if err != nil {
		panic("Error unmarshaling JSON")
	}
	return request.Streams, request.Total
}

func getRequestFromJSON(body []byte) (*streamRequest, error) {
	var req = new(streamRequest)
	err := json.Unmarshal(body, &req)
	return req, err
}

func (request streamRequest) String() string {
	var toString string
	toString = fmt.Sprintf("Request: %s\n", request.Streams)
	return toString
}

// StreamsToStrings converts a slice of Streams into a slice of their string representations
func StreamsToStrings(streams []Stream) []string {
	toStrings := make([]string, len(streams), len(streams))
	for i := range streams {
		toStrings[i] = streams[i].String()
	}
	return toStrings
}

func (stream Stream) String() string {
	var toString string
	toString = fmt.Sprintf("[%s](fg-cyan) : %s playing [%s](fg-blue) with %d Viewers for %s\n",
		stream.Channel.Status, stream.Channel.DisplayName, stream.Game, stream.Viewers,
		stream.printRoundedStreamTime())
	return toString
}

func (channel Channel) String() string {
	var toString string
	toString = fmt.Sprintf("Mature: %t, Status: %s, Broadcaster Language: %s, ", channel.Mature,
		channel.Status, channel.BroadcasterLanguage)
	toString = fmt.Sprintf("%sDisplay Name: %s, Game: %s, Language: %s, ", toString,
		channel.DisplayName, channel.Game, channel.Language)
	toString = fmt.Sprintf("%sName: %s, Updated At: %s, Created At: %s, URL: %s", toString,
		channel.Name, channel.UpdatedAt, channel.CreatedAt, channel.URL)
	return toString
}

func (stream *Stream) printRoundedStreamTime() string {
	var toString string
	toString = fmt.Sprintf("%s", time.Now().Round(time.Minute).Sub(stream.CreatedAt.Round(time.Minute)))
	return toString
}
