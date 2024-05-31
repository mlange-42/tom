package app

import (
	"context"
	"fmt"
	"time"

	"github.com/mlange-42/tom/config"
	"github.com/mlange-42/tom/render"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/text"
)

const redrawInterval = 250 * time.Millisecond

// rootID is the ID assigned to the root container.
const rootID = "root"

// Terminal implementations
const (
	TermboxTerminal = "termbox"
	TCellTerminal   = "tcell"
)

type App struct {
	location string
	data     *config.MeteoResult

	current  *text.Text
	forecast *text.Text
}

func New(location string, data *config.MeteoResult) *App {
	return &App{
		location: location,
		data:     data,
	}
}

func (a *App) Run(term string) error {
	var t terminalapi.Terminal
	var err error
	switch term {
	case TermboxTerminal:
		t, err = termbox.New(termbox.ColorMode(terminalapi.ColorMode256))
	case TCellTerminal:
		t, err = tcell.New(tcell.ColorMode(terminalapi.ColorMode256))
	default:
		err = fmt.Errorf("unknown terminal implementation '%s' specified. Please choose between 'termbox' and 'tcell'", term)
	}

	if err != nil {
		return err
	}

	defer t.Close()

	ctx, cancel := context.WithCancel(context.Background())

	if err := a.createWidgets(); err != nil {
		panic(err)
	}

	c, err := container.New(t, container.ID(rootID))
	if err != nil {
		panic(err)
	}

	gridOpts, err := contLayout(a)
	if err != nil {
		panic(err)
	}
	if err := c.Update(rootID, gridOpts...); err != nil {
		panic(err)
	}

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == keyboard.KeyEsc || k.Key == keyboard.KeyCtrlC {
			cancel()
		}
	}

	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(redrawInterval)); err != nil {
		panic(err)
	}

	return nil
}

func (a *App) createWidgets() error {
	var err error
	a.current, err = text.New()
	if err != nil {
		return err
	}

	var now time.Time
	loc, err := time.LoadLocation(a.data.Location.TimeZone)
	if err == nil {
		now = time.Now().In(loc)
	}

	renderer := render.NewRenderer(a.data)

	a.current.Write(
		fmt.Sprintf("%s (%0.2f°N, %0.2f°E)  %s | %s",
			a.location, a.data.Location.Lat, a.data.Location.Lon,
			now.Format(config.TimeLayout),
			renderer.Current(),
		))

	a.forecast, err = text.New(text.ScrollRunes('↑', '↓'))
	if err != nil {
		return err
	}

	for i, t := range a.data.DailyTime {
		a.forecast.Write(fmt.Sprintf("%-11s | %s\n", t.Format(config.DateLayoutShort), renderer.DaySummary(i)))
		a.forecast.Write(renderer.DaySixHourly(i*4) + "\n")
	}

	return nil
}

func contLayout(a *App) ([]container.Option, error) {
	rows := []container.Option{
		container.SplitHorizontal(
			container.Top(
				container.Border(linestyle.Light),
				container.BorderTitle("Current───Press Esc to quit"),
				container.PlaceWidget(a.current),
			),
			container.Bottom(
				container.Border(linestyle.Light),
				container.BorderTitle("Forecast───Scroll with mouse or arrow keys"),
				container.PlaceWidget(a.forecast),
				container.Focused(),
			),
			container.SplitFixed(3),
		),
	}

	return rows, nil
}
