package app

import (
	"context"
	"fmt"
	"time"

	"github.com/mlange-42/tom/api"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/linechart"
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
	data     *api.MeteoResult
}

func New(location string, data *api.MeteoResult) *App {
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

	c, err := container.New(t, container.ID(rootID))
	if err != nil {
		panic(err)
	}

	gridOpts, err := gridLayout(a)
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

func gridLayout(a *App) ([]container.Option, error) {
	t, err := text.New()
	if err != nil {
		return nil, err
	}
	t.Write(fmt.Sprintf("%s  %0.2f°N, %0.2f°E", a.location, a.data.Location.Lat, a.data.Location.Lon))

	xLabels := map[int]string{}
	/*for i := 12; i < len(a.data.HourlyTime); i += 24 {
		t := a.data.HourlyTime[i]
		xLabels[i] = t.Format(api.DateLayoutShort)
	}*/
	for i := 0; i < len(a.data.HourlyTime); i++ {
		t := a.data.HourlyTime[i]
		xLabels[i] = t.Format(api.DateLayoutShort)
	}

	linesTemp, err := linechart.New()
	if err != nil {
		return nil, err
	}
	linesTemp.Series("temperature", a.data.Hourly[string(api.HourlyTemp)],
		linechart.SeriesXLabels(xLabels),
	)

	linesPrecip, err := linechart.New()
	if err != nil {
		return nil, err
	}
	linesPrecip.Series("precip", a.data.Hourly[string(api.HourlyPrecip)],
		linechart.SeriesXLabels(xLabels),
	)

	rows := []grid.Element{
		grid.RowHeightFixed(3,
			grid.Widget(t,
				container.Border(linestyle.Light),
				container.BorderTitle("Press Esc to quit"),
			),
		),
		grid.RowHeightPerc(50,
			grid.Widget(linesTemp,
				container.Border(linestyle.Light),
				container.BorderTitle("Temperature"),
			),
		),
		grid.RowHeightPerc(50,
			grid.Widget(linesPrecip,
				container.Border(linestyle.Light),
				container.BorderTitle("Precipitation"),
			),
		),
	}

	builder := grid.New()
	builder.Add(rows...)

	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return gridOpts, nil
}
