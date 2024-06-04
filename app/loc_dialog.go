package app

import (
	"fmt"
	"log"
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/mlange-42/tom/api"
	"github.com/mlange-42/tom/config"
	"github.com/rivo/tview"
)

type LocationDialog struct {
	cliArgs config.CliArgs
	cached  map[string]config.Location
}

func NewLocationDialog(cliArgs config.CliArgs, cached map[string]config.Location) *LocationDialog {
	return &LocationDialog{
		cliArgs: cliArgs,
		cached:  cached,
	}
}

func (d *LocationDialog) Run() error {
	locations, err := api.GetLocations(d.cliArgs.Location, 100)
	if err != nil {
		return err
	}

	if len(locations) == 0 {
		return fmt.Errorf("no locations found for '%s'", d.cliArgs.Location)
	}
	if len(locations) == 1 {
		coords, err := d.updateCache(&locations[0])
		if err != nil {
			log.Fatal(err)
		}
		d.cliArgs.Coords = coords
		a := New(d.cliArgs)
		a.Run()
		return nil
	}

	app := tview.NewApplication()

	grid := tview.NewGrid().
		SetRows(3, 0, 1).
		SetColumns(-1).
		SetBorders(false)

	header := tview.NewTextView().
		SetWrap(false).
		SetText("Please select a location")
	header.SetBorder(true)

	grid.AddItem(header, 0, 0, 1, 1, 0, 0, false)

	data := NewLocationTable(locations)
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetContent(&data).
		SetEvaluateAllRows(true).
		SetSeparator(tview.Borders.Vertical)
	table.SetBorder(true)

	grid.AddItem(table, 1, 0, 1, 1, 0, 0, true)

	help := tview.NewTextView().
		SetWrap(false).
		SetText("Exit: ESC  Scroll: ←→↕  Select: ENTER")
	grid.AddItem(help, 2, 0, 1, 1, 0, 0, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.Stop()
			return nil
		} else if event.Key() == tcell.KeyEnter {
			row, _ := table.GetSelection()
			if row == 0 {
				return nil
			}
			app.Stop()
			loc := locations[row-1]

			coords, err := d.updateCache(&loc)
			if err != nil {
				log.Fatal(err)
			}
			d.cliArgs.Coords = coords

			a := New(d.cliArgs)
			a.Run()
			return nil
		}
		return event
	})

	if err := app.SetRoot(grid, true).Run(); err != nil {
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

	d.cached[d.cliArgs.Location] = coords
	err := config.SaveLocations(d.cached)
	if err != nil {
		return coords, err
	}
	return coords, nil
}

type LocationTable struct {
	tview.TableContentReadOnly

	locations []config.GeoResultEntry
	header    []string
}

func NewLocationTable(locations []config.GeoResultEntry) LocationTable {
	return LocationTable{
		locations: locations,
		header:    []string{"Name", "Region", "State", "Coordinates", "Population"},
	}
}

func (t *LocationTable) GetCell(row, column int) *tview.TableCell {
	if row == 0 {
		return tview.NewTableCell(t.header[column])
	}
	var text string
	loc := &t.locations[row-1]

	switch column {
	case 0:
		runes := []rune(loc.Name)
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
	cell := tview.NewTableCell(text)
	if column == 0 {
		cell.SetExpansion(1)
	}

	return cell
}

func (t *LocationTable) GetRowCount() int {
	return len(t.locations) + 1
}

func (t *LocationTable) GetColumnCount() int {
	return 5
}
