package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mlange-42/tom/app"
	"github.com/mlange-42/tom/config"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCommand().Execute(); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		fmt.Print("\nRun `tom -h` for help!\n\n")
		os.Exit(1)
	}
}

// rootCommand sets up the CLI
func rootCommand() *cobra.Command {
	cli := config.CliArgs{}

	services := ""
	for _, s := range config.Services {
		services += fmt.Sprintf("    - %-5s - %s\n", s.Name, s.Description)
	}

	root := cobra.Command{
		Use:           "tom",
		Short:         "Terminal for Open-Meteo.",
		Long:          `Terminal for Open-Meteo.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ArbitraryArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("please specify a location")
			}
			cli.Location = strings.ToLower(strings.Join(args, " "))
			forceApi := strings.HasSuffix(cli.Location, "?")
			cli.Location = strings.TrimSuffix(cli.Location, "?")

			if cli.Days < 1 || cli.Days > 16 {
				return fmt.Errorf("parameter --days must be in range [1, 16]")
			}

			forecasterFound := false
			for _, s := range config.Services {
				if strings.EqualFold(s.Name, cli.Service.Name) {
					cli.Service = s
					forecasterFound = true
				}
			}
			if !forecasterFound {
				return fmt.Errorf("service '%s' not found. Available services:\n%s", cli.Service.Name, services)
			}

			cached, err := config.LoadLocations()
			if err != nil {
				return err
			}

			coords, ok := cached[cli.Location]
			if ok && !forceApi {
				cli.Coords = coords
				a := app.New(cli)
				if err := a.Run(); err != nil {
					return err
				}
				return nil
			}

			a := app.NewLocationDialog(cli, cached)
			if err := a.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	root.Flags().IntVarP(&cli.Days, "days", "d", 7, "Number of forecast days in range [1, 16]")
	root.Flags().StringVarP(&cli.Service.Name, "service", "s", "OM", "Forecast service.\n"+services)

	root.Flags().SortFlags = false

	return &root
}
