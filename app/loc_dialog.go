package app

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mlange-42/tom/api"
	"github.com/mlange-42/tom/config"
	"github.com/rivo/tview"
)

type LocationDialog struct {
	location string
	cached   map[string]config.Location
}

func NewLocationDialog(location string, cached map[string]config.Location) *LocationDialog {
	return &LocationDialog{
		location: location,
		cached:   cached,
	}
}

func (d *LocationDialog) Run() error {
	locations, err := api.GetLocations(d.location, 100)
	if err != nil {
		return err
	}

	if len(locations) == 1 {
		coords, err := d.updateCache(&locations[0])
		if err != nil {
			log.Fatal(err)
		}
		a := New(strings.ToTitle(d.location), coords)
		a.Run()
		return nil
	}

	app := tview.NewApplication()
	data := &LocationTable{Locations: locations}
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetContent(data).
		SetEvaluateAllRows(true).
		SetSeparator(tview.Borders.Vertical)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.Stop()
			return nil
		} else if event.Key() == tcell.KeyEnter {
			app.Stop()
			row, _ := table.GetSelection()
			loc := locations[row]

			coords, err := d.updateCache(&loc)
			if err != nil {
				log.Fatal(err)
			}

			a := New(strings.ToTitle(d.location), coords)
			a.Run()
			return nil
		}
		return event
	})

	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
	return nil
}

func (d *LocationDialog) updateCache(loc *config.GeoResultEntry) (config.Location, error) {
	coords := config.Location{
		Lat:      loc.Latitude,
		Lon:      loc.Longitude,
		TimeZone: loc.TimeZone,
	}

	d.cached[d.location] = coords
	err := config.SaveLocations(d.cached)
	if err != nil {
		return coords, err
	}
	return coords, nil
}

type LocationTable struct {
	tview.TableContentReadOnly

	Locations []config.GeoResultEntry
}

func (t *LocationTable) GetCell(row, column int) *tview.TableCell {
	var text string
	loc := &t.Locations[row]

	switch column {
	case 0:
		runes := []rune(loc.Name)
		if len(runes) > 25 {
			runes = runes[:25]
		}
		text = string(runes)
	case 1:
		text = loc.Admin1
	case 2:
		text = loc.Country
	case 3:
		text = fmt.Sprintf("%6.2f°N %7.2f°E", loc.Latitude, loc.Longitude)
	case 4:
		pop := loc.Population
		if pop > 1_000_000 {
			text = fmt.Sprintf("%.1fM", float64(pop/1_000_000))
		} else if pop > 1_000 {
			text = fmt.Sprintf("%.0fk", math.Round(float64(pop/1000)))
		} else {
			text = fmt.Sprintf("%d", pop)
		}
	}

	return tview.NewTableCell(text)
}

func (t *LocationTable) GetRowCount() int {
	return len(t.Locations)
}

func (t *LocationTable) GetColumnCount() int {
	return 5
}
