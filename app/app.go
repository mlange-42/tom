package app

import (
	"fmt"
	"strings"

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

	currentWeather *tview.TextView
	forecast       *tview.TextView
	plots          *tview.TextView
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

	if a.cliArgs.SetDefault {
		if err := config.SaveCliArgs(&a.cliArgs); err != nil {
			return err
		}
	}

	app := tview.NewApplication()
	pages := tview.NewPages()

	if err = a.createWidgets(); err != nil {
		return err
	}
	forecasts := a.createForecastsPage()
	plots := a.createPlotsPage()

	pages.AddAndSwitchToPage("forecast", forecasts, true)
	pages.AddPage("plots", plots, true, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.Stop()
			return nil
		} else if event.Key() == tcell.KeyTab {
			if a.currentWeather.HasFocus() {
				_, page := pages.GetFrontPage()
				if page == forecasts {
					app.SetFocus(a.forecast)
				} else {
					app.SetFocus(a.plots)
				}
			} else {
				app.SetFocus(a.currentWeather)
			}
			return nil
		} else if event.Rune() == ' ' {
			_, page := pages.GetFrontPage()
			if page == forecasts {
				pages.SwitchToPage("plots")
			} else {
				pages.SwitchToPage("forecast")
			}
		}
		return event
	})

	if err := app.SetRoot(pages, true).Run(); err != nil {
		return err
	}

	return nil
}

func (a *App) createWidgets() error {
	now := a.data.Current.Time

	renderer := render.NewRenderer(a.data)
	a.currentWeather = tview.NewTextView().
		SetWrap(false).
		SetDynamicColors(true).
		SetText(
			fmt.Sprintf("%s (%0.2f°N, %0.2f°E)  %s | %s",
				strings.ToTitle(a.cliArgs.Location), a.data.Location.Lat, a.data.Location.Lon,
				now.Format(config.TimeLayout),
				renderer.Current(),
			))
	a.currentWeather.SetBorder(true)
	a.currentWeather.SetTitle(" Current weather ")

	builder := strings.Builder{}
	for i, t := range a.data.DailyTime {
		tag := "[-]"
		if t.YearDay() == now.YearDay() {
			tag = "[yellow]"
		}
		newline := ""
		if i < len(a.data.DailyTime)-1 {
			newline = "\n"
		}
		_, err := builder.WriteString(
			fmt.Sprintf("%s%-11s[-] | %s\n%s%s",
				tag, t.Format(config.DateLayoutShort), renderer.DaySummary(i), renderer.DaySixHourly(i*4), newline),
		)
		if err != nil {
			return err
		}
	}

	a.forecast = tview.NewTextView().
		SetWrap(false).
		SetDynamicColors(true).
		SetText(builder.String())
	a.forecast.SetBorder(true)
	a.forecast.SetTitle(fmt.Sprintf(" %s %d days forecast ", a.cliArgs.Service.Description, a.cliArgs.Days))

	a.plots = tview.NewTextView().
		SetWrap(false).
		SetDynamicColors(true).
		SetText(renderer.Charts(now))
	a.plots.SetBorder(true)
	a.plots.SetTitle(fmt.Sprintf(" %s %d days charts ", a.cliArgs.Service.Description, a.cliArgs.Days))

	return nil
}

func (a *App) createForecastsPage() tview.Primitive {
	grid := tview.NewGrid().
		SetRows(3, 0, 1).
		SetColumns(len(data.DayLayout[0]) + 2).
		SetBorders(false)

	grid.AddItem(a.currentWeather, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(a.forecast, 1, 0, 1, 1, 0, 0, true)

	help := tview.NewTextView().
		SetWrap(false).
		SetText("Exit: ESC  Pages: SPACE  Scroll: ←→↕  Focus: TAB")
	grid.AddItem(help, 2, 0, 1, 1, 0, 0, false)

	return grid
}

func (a *App) createPlotsPage() tview.Primitive {
	grid := tview.NewGrid().
		SetRows(3, 0, 1).
		SetColumns(len(data.DayLayout[0]) + 2).
		SetBorders(false)

	grid.AddItem(a.currentWeather, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(a.plots, 1, 0, 1, 1, 0, 0, true)

	help := tview.NewTextView().
		SetWrap(false).
		SetText("Exit: ESC  Pages: SPACE  Scroll: ←→↕  Focus: TAB")
	grid.AddItem(help, 2, 0, 1, 1, 0, 0, false)

	return grid
}
