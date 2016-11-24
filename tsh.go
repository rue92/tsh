package main

import "fmt"
import "os"
import "log"
import "github.com/rue92/tsh/twitch"
import ui "github.com/gizak/termui"

type ShellState struct {
	currentOffset uint32
	currentLimit  uint8
	lastTotal     uint32
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

	ls := ui.NewList()
	ls.Items = strs
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "Choices"
	ls.Height = len(strs) + 2
	ls.Width = ui.TermWidth()
	//	ls.Y = 0

	results := ui.NewList()
	results.Height = ui.TermHeight() - len(strs) - 2
	results.Width = ui.TermWidth()
	state := ShellState{0, uint8(ui.TermHeight() - len(strs) - 2), 0}
	var games []twitch.Game
	games, state.lastTotal = twitch.GetGames(state.currentLimit, state.currentOffset)
	results.Items = twitch.GamesToStrings(games)
	results.BorderLabel = "Games"

	logPar := ui.NewPar(state.String())
	logPar.BorderLabel = "State"
	logPar.Height = 4

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, ls),
			ui.NewCol(6, 0, logPar)),
		ui.NewRow(
			ui.NewCol(12, 0, results)))
	ui.Body.Align()

	ui.Handle("/sys/kbd/r", func(ui.Event) {
		var games []twitch.Game
		games, state.lastTotal = twitch.GetGames(state.currentLimit, state.currentOffset)
		results.Items = twitch.GamesToStrings(games)
		logPar.Text = state.String()
		ui.Clear()
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/1", func(ui.Event) {
		var games []twitch.Game
		games, state.lastTotal = twitch.GetGames(state.currentLimit, state.currentOffset)
		results.Items = twitch.GamesToStrings(games)
		logPar.Text = state.String()
		ui.Clear()
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/n", func(ui.Event) {
		var games []twitch.Game
		state.currentOffset += uint32(state.currentLimit)
		games, state.lastTotal = twitch.GetGames(state.currentLimit, state.currentOffset)
		results.Items = twitch.GamesToStrings(games)
		logPar.Text = state.String()
		ui.Clear()
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/p", func(ui.Event) {
		state.currentOffset -= uint32(state.currentLimit)
		if state.currentOffset > state.lastTotal {
			state.currentOffset = 0
		}
		var games []twitch.Game
		games, state.lastTotal = twitch.GetGames(state.currentLimit, state.currentOffset)
		results.Items = twitch.GamesToStrings(games)
		logPar.Text = state.String()
		ui.Clear()
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/2", func(ui.Event) {
		results.Items = twitch.StreamsToStrings(twitch.GetStreams())
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
	return fmt.Sprintf("Offset: %d, Limit: %d, LastTotal: %d", state.currentOffset, state.currentLimit, state.lastTotal)
}
