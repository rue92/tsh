package main

import "fmt"
import "os"
import "log"
import "time"
import "github.com/rue92/tsh/twitch"
import ui "github.com/gizak/termui"

type pagerEnum int

const (
	games pagerEnum = iota
	streams
)

type shellState struct {
	pager pager
}

type pager interface {
	next() interface{}
	prev() interface{}
	current() interface{}
	pagerType() pagerEnum
}

type requester struct {
	offset uint32
	limit  uint8
	total  uint32
	name   string
}

type gameRequester requester
type streamRequester requester

func (requester *gameRequester) next() interface{} {
	var games []twitch.Game
	requester.offset += uint32(requester.limit)
	if requester.offset > requester.total {
		requester.offset = 0
	}
	games, requester.total = twitch.GetGames(requester.limit, requester.offset)
	return games
}

func (requester *gameRequester) prev() interface{} {
	var games []twitch.Game
	requester.offset -= uint32(requester.limit)
	if requester.offset > requester.total {
		requester.offset = requester.total - uint32(requester.limit)
	}
	games, requester.total = twitch.GetGames(requester.limit, requester.offset)
	return games
}

func (requester *gameRequester) current() interface{} {
	var games []twitch.Game
	games, requester.total = twitch.GetGames(requester.limit, requester.offset)
	return games
}

func (requester *gameRequester) pagerType() pagerEnum { return games }

func (requester *streamRequester) next() interface{} {
	var streams []twitch.Stream
	requester.offset += uint32(requester.limit)
	if requester.offset > requester.total {
		requester.offset = 0
	}
	streams, requester.total = twitch.GetStreams(requester.limit, requester.offset)
	return streams
}

func (requester *streamRequester) prev() interface{} {
	var streams []twitch.Stream
	requester.offset -= uint32(requester.limit)
	if requester.offset > requester.total {
		requester.offset = requester.total - uint32(requester.limit)
	}
	streams, requester.total = twitch.GetStreams(requester.limit, requester.offset)
	return streams
}

func (requester *streamRequester) current() interface{} {
	var streams []twitch.Stream
	streams, requester.total = twitch.GetStreams(requester.limit, requester.offset)
	return streams
}

func (requester *streamRequester) pagerType() pagerEnum { return streams }

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
		"[1] Top Games",
		"[2] Top Streams",
		"[3] Followed"}
	gameRequester := gameRequester{0, 24, 0, "Games"}
	streamRequester := streamRequester{0, 24, 0, "Streams"}
	state := shellState{&gameRequester}

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

	var refresh = func(e ui.Event) {
		results.Items = requestToString(state.pager.current())
		logPar.Text = state.String()
		ui.Clear()
		ui.Render(ui.Body)
	}

	ui.Handle("/sys/kbd/r", func(e ui.Event) {
		refresh(e)
	})
	ui.Handle("/sys/kbd/1", func(e ui.Event) {
		state.pager = &gameRequester
		results.Items = requestToString(state.pager.current())
		results.BorderLabel = gameRequester.name
		logPar.Text = state.String()
		ui.Clear()
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/2", func(e ui.Event) {
		state.pager = &streamRequester
		results.Items = requestToString(state.pager.current())
		results.BorderLabel = streamRequester.name
		logPar.Text = state.String()
		ui.Clear()
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/n", func(e ui.Event) {
		results.Items = requestToString(state.pager.next())
		logPar.Text = state.String()
		ui.Clear()
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/p", func(e ui.Event) {
		results.Items = requestToString(state.pager.prev())
		logPar.Text = state.String()
		ui.Clear()
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/q", func(e ui.Event) {
		ui.StopLoop()
	})

	ui.Merge("/timer/2s", ui.NewTimerCh(time.Second*2))
	ui.Handle("/timer/2s", func(e ui.Event) {
		if state.pager.pagerType() == streams ||
			state.pager.pagerType() == games {
			refresh(e)
		}
	})
	ui.Render(ui.Body)
	ui.Loop()
}

func (state shellState) String() string {
	return fmt.Sprintf("%s", state.pager)
}

func (requester gameRequester) String() string {
	return fmt.Sprintf("Offset: %d, Limit: %d, Total: %d", requester.offset, requester.limit, requester.total)
}

func (requester streamRequester) String() string {
	return fmt.Sprintf("Offset: %d, Limit: %d, Total: %d", requester.offset, requester.limit, requester.total)
}

func requestToString(data interface{}) []string {
	switch request := data.(type) {
	case []twitch.Game:
		return twitch.GamesToStrings(request)
	case []twitch.Stream:
		return twitch.StreamsToStrings(request)
	}
	return make([]string, 1)
}
