package main

import "fmt"
import "os"
import "log"
import "github.com/rue92/tsh/twitch"
import ui "github.com/gizak/termui"

const (
	R_GAMES = iota
	R_STREAMS
)

type ShellState struct {
	pager Pager
}

type Pager interface {
	next() interface{}
	prev() interface{}
	current() interface{}
}

type Requester struct {
	offset uint32
	limit  uint8
	total  uint32
	name   string
}

type GameRequester Requester
type StreamRequester Requester

func (requester *GameRequester) next() interface{} {
	var games []twitch.Game
	requester.offset += uint32(requester.limit)
	if requester.offset > requester.total {
		requester.offset = 0
	}
	games, requester.total = twitch.GetGames(requester.limit, requester.offset)
	return games
}

func (requester *GameRequester) prev() interface{} {
	var games []twitch.Game
	requester.offset -= uint32(requester.limit)
	if requester.offset > requester.total {
		requester.offset = requester.total - uint32(requester.limit)
	}
	games, requester.total = twitch.GetGames(requester.limit, requester.offset)
	return games
}

func (requester *GameRequester) current() interface{} {
	var games []twitch.Game
	games, requester.total = twitch.GetGames(requester.limit, requester.offset)
	return games
}

func (requester *StreamRequester) next() interface{} {
	var streams []twitch.Stream
	requester.offset += uint32(requester.limit)
	if requester.offset > requester.total {
		requester.offset = 0
	}
	streams, requester.total = twitch.GetStreams(requester.limit, requester.offset)
	return streams
}

func (requester *StreamRequester) prev() interface{} {
	var streams []twitch.Stream
	requester.offset -= uint32(requester.limit)
	if requester.offset > requester.total {
		requester.offset = requester.total - uint32(requester.limit)
	}
	streams, requester.total = twitch.GetStreams(requester.limit, requester.offset)
	return streams
}

func (requester *StreamRequester) current() interface{} {
	var streams []twitch.Stream
	streams, requester.total = twitch.GetStreams(requester.limit, requester.offset)
	return streams
}

func main() {
	f, _ := os.OpenFile("tsh.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer f.Close()

	log.SetOutput(f)

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	strs := []string{
		"[1] Games",
		"[2] Channels",
		"[3] Followed"}
	gameRequester := GameRequester{0, 24, 0, "Games"}
	streamRequester := StreamRequester{0, 24, 0, "Streams"}
	state := ShellState{&gameRequester}

	ls := ui.NewList()
	ls.Items = strs
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "Choices"
	ls.Height = len(strs) + 2
	ls.Width = ui.TermWidth()

	results := ui.NewList()
	results.Height = ui.TermHeight() - len(strs) - 2
	results.Width = ui.TermWidth()
	results.Items = twitch.GamesToStrings(state.pager.current().([]twitch.Game))
	results.BorderLabel = gameRequester.name

	logPar := ui.NewPar(state.String())
	logPar.BorderLabel = "State"
	logPar.Height = len(strs) + 2

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, ls),
			ui.NewCol(6, 0, logPar)),
		ui.NewRow(
			ui.NewCol(12, 0, results)))
	ui.Body.Align()

	ui.Handle("/sys/kbd/r", func(ui.Event) {
		results.Items = RequestToString(state.pager.current())
		logPar.Text = state.String()
		ui.Clear()
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/1", func(ui.Event) {
		state.pager = &gameRequester
		results.Items = RequestToString(state.pager.current())
		results.BorderLabel = gameRequester.name
		logPar.Text = state.String()
		ui.Clear()
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/n", func(ui.Event) {
		results.Items = RequestToString(state.pager.next())
		logPar.Text = state.String()
		ui.Clear()
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/p", func(ui.Event) {
		results.Items = RequestToString(state.pager.prev())
		logPar.Text = state.String()
		ui.Clear()
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/2", func(ui.Event) {
		state.pager = &streamRequester
		results.Items = RequestToString(state.pager.current())
		results.BorderLabel = streamRequester.name
		logPar.Text = state.String()
		ui.Clear()
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Render(ui.Body)
	ui.Loop()
}

func (state ShellState) String() string {
	return fmt.Sprintf("%s", state.pager)
}

func (requester GameRequester) String() string {
	return fmt.Sprintf("Offset: %d, Limit: %d, Total: %d", requester.offset, requester.limit, requester.total)
}

func (requester StreamRequester) String() string {
	return fmt.Sprintf("Offset: %d, Limit: %d, Total: %d", requester.offset, requester.limit, requester.total)
}

func RequestToString(data interface{}) []string {
	switch request := data.(type) {
	case []twitch.Game:
		return twitch.GamesToStrings(request)
	case []twitch.Stream:
		return twitch.StreamsToStrings(request)
	}
	return make([]string, 1)
}
