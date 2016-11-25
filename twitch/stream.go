package twitch

import "time"
import "fmt"
import "log"
import "strconv"
import "encoding/json"
import "net/http"
import "io/ioutil"

type StreamRequest struct {
	Streams []Stream `json:"streams"`
	Total   uint32   `json:"_total"`
}

type Channel struct {
	Mature              bool      `json:"mature"`
	Status              string    `json:"status"`
	BroadcasterLanguage string    `json:"broadcaster_language"`
	DisplayName         string    `json:"display_name"`
	Game                string    `json:"game"`
	Delay               float64   `json:"delay"`
	Language            string    `json:"language"`
	Id                  uint64    `json:"_id"`
	Name                string    `json:"name"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	Partner             bool      `json:"partner"`
	Url                 string    `json:"url"`
}

type Stream struct {
	Game        string    `json:"game"`
	Viewers     uint64    `json:"viewers"`
	AverageFps  float64   `json:"average_fps"`
	Delay       float64   `json:"delay"`
	VideoHeight uint64    `json:"video_height"`
	IsPlaylist  bool      `json:"is_playlist"`
	CreatedAt   time.Time `json:"created_at"`
	Id          uint64    `json:"_id"`
	Channel     Channel   `json:"channel"`
}

func GetStream(user string) (stream Stream) {
	client := &http.Client{}
	streamUrl := TwitchApiBaseUrl + "streams/" + user + "/"
	req, err := http.NewRequest("GET", streamUrl, nil)
	if err != nil {
		log.Println("Couldn't create request")
	}
	req.Header.Add("Accept", "application/vnd.twitchtv.v3+json")
	req.Header.Add("Client-ID", TshClientId)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Couldn't perform HTTP request")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Couldn't read message body")
	}
	resp.Body.Close()

	request, err := GetRequestFromJson(body)
	if err != nil {
		panic("Error unmarshaling JSON")
	}
	return request.Streams[0]
}

func GetStreams(limit uint8, offset uint32) (streams []Stream, total uint32) {
	client := &http.Client{}
	streamUrl := TwitchApiBaseUrl + "streams/"
	req, err := http.NewRequest("GET", streamUrl, nil)
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

	request, err := GetRequestFromJson(body)
	if err != nil {
		panic("Error unmarshaling JSON")
	}
	return request.Streams, request.Total
}

func GetRequestFromJson(body []byte) (*StreamRequest, error) {
	var req = new(StreamRequest)
	err := json.Unmarshal(body, &req)
	return req, err
}

func (request StreamRequest) String() string {
	var toString string
	toString = fmt.Sprintf("Request: %s\n", request.Streams)
	return toString
}

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
		channel.Name, channel.UpdatedAt, channel.CreatedAt, channel.Url)
	return toString
}

func (stream *Stream) printRoundedStreamTime() string {
	var toString string = ""
	toString = fmt.Sprintf("%s", time.Now().Round(time.Minute).Sub(stream.CreatedAt.Round(time.Minute)))
	return toString
}
