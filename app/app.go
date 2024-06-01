package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mlange-42/tom/api"
	"github.com/mlange-42/tom/config"
	"github.com/mlange-42/tom/data"
	"github.com/mlange-42/tom/render"
	"github.com/rivo/tview"
)

type App struct {
	cliArgs config.CliArgs

	data *config.MeteoResult
}

func New(cliArgs config.CliArgs) *App {
	return &App{
		cliArgs: cliArgs,
	}
}

func (a *App) Run() error {
	var err error
	a.data, err = api.GetMeteo(a.cliArgs)
	if err != nil {
		return err
	}

	var now time.Time
	loc, err := time.LoadLocation(a.data.Location.TimeZone)
	if err == nil {
		now = time.Now().In(loc)
	}

	app := tview.NewApplication()
	pages := tview.NewPages()

	grid := tview.NewGrid().
		SetRows(3, 0, 1).
		SetColumns(len(data.DayLayout[0]) + 2).
		SetBorders(false)

	renderer := render.NewRenderer(a.data)

	current := tview.NewTextView().
		SetWrap(false).
		SetDynamicColors(true).
		SetText(
			fmt.Sprintf("%s (%0.2f°N, %0.2f°E)  %s | %s",
				a.cliArgs.Location, a.data.Location.Lat, a.data.Location.Lon,
				now.Format(config.TimeLayout),
				renderer.Current(),
			))
	current.SetBorder(true)
	current.SetTitle(" Current weather ")
	grid.AddItem(current, 0, 0, 1, 1, 0, 0, false)

	builder := strings.Builder{}
	for i, t := range a.data.DailyTime {
		_, err := builder.WriteString(
			fmt.Sprintf("%-11s | %s\n%s\n",
				t.Format(config.DateLayoutShort), renderer.DaySummary(i), renderer.DaySixHourly(i*4)),
		)
		if err != nil {
			return err
		}
	}
	forecast := tview.NewTextView().
		SetWrap(false).
		SetDynamicColors(true).
		SetText(builder.String())
	forecast.SetBorder(true)
	forecast.SetTitle(" 7 days forecast ")

	grid.AddItem(forecast, 1, 0, 1, 1, 0, 0, true)

	help := tview.NewTextView().
		SetWrap(false).
		SetText("Exit: Esc  Focus: Tab  Scroll: ←→↕")
	grid.AddItem(help, 2, 0, 1, 1, 0, 0, false)

	pages.AddAndSwitchToPage("forecast", grid, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.Stop()
			return nil
		} else if event.Key() == tcell.KeyTab {
			if current.HasFocus() {
				app.SetFocus(forecast)
			} else {
				app.SetFocus(current)
			}
			return nil
		}
		return event
	})

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}

	return nil
}
