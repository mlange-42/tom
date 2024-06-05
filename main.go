package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mlange-42/tom/app"
	"github.com/mlange-42/tom/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	cmd, err := rootCommand()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		fmt.Print("\nRun `tom -h` for help!\n\n")
		os.Exit(1)
	}
	if err := cmd.Execute(); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		fmt.Print("\nRun `tom -h` for help!\n\n")
		os.Exit(1)
	}
}

// rootCommand sets up the CLI
func rootCommand() (*cobra.Command, error) {
	defaults, err := config.LoadCliArgs()
	if err != nil {
		return nil, err
	}
	cli := config.CliArgs{}

	services := ""
	for _, s := range config.Services {
		services += fmt.Sprintf("    - %-5s - %s\n", s.Name, s.Description)
	}

	var root cobra.Command
	root = cobra.Command{
		Use:           "tom [LOCATION]",
		Short:         "Terminal for Open-Meteo.",
		Long:          `Terminal for Open-Meteo.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ArbitraryArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			forceSearchLoc := false
			if len(args) == 0 {
				if defaults.Location == "" {
					cmd.Help()
					fmt.Println("\nPlease specify a location!")
					return nil
				}
			} else {
				defaults.Location = strings.ToLower(strings.Join(args, " "))
				forceSearchLoc = strings.HasSuffix(defaults.Location, "?")
				defaults.Location = strings.TrimSuffix(defaults.Location, "?")
			}

			flagUsed := map[string]bool{}
			root.Flags().Visit(func(f *pflag.Flag) {
				flagUsed[f.Name] = true
			})

			if _, ok := flagUsed["days"]; ok {
				defaults.Days = cli.Days
			}
			if _, ok := flagUsed["past-days"]; ok {
				defaults.PastDays = cli.PastDays
			}
			if _, ok := flagUsed["service"]; ok {
				defaults.Service = cli.Service
			}
			defaults.SetDefault = cli.SetDefault

			if defaults.Days < 1 || defaults.Days > 16 {
				return fmt.Errorf("parameter --days must be in range [1, 16]")
			}

			forecasterFound := false
			for _, s := range config.Services {
				if strings.EqualFold(s.Name, defaults.Service.Name) {
					defaults.Service = s
					forecasterFound = true
				}
			}
			if !forecasterFound {
				return fmt.Errorf("service '%s' not found. Available services:\n%s", defaults.Service.Name, services)
			}

			cached, err := config.LoadLocations()
			if err != nil {
				return err
			}

			coords, ok := cached[defaults.Location]
			if ok && !forceSearchLoc {
				defaults.Coords = coords
				a := app.New(defaults)
				if err := a.Run(); err != nil {
					return err
				}
				return nil
			}

			a := app.NewLocationDialog(defaults, cached)
			if err := a.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	root.Flags().IntVarP(&cli.Days, "days", "d", 7, "Number of forecast days in range [1, 16]")
	root.Flags().IntVarP(&cli.PastDays, "past-days", "p", 0, "Number of past days to include")
	root.Flags().StringVarP(&cli.Service.Name, "service", "s", "OM", "Forecast service:\n"+services)
	root.Flags().BoolVarP(&cli.SetDefault, "default", "", false, "Save given location and settings as default")

	root.Flags().SortFlags = false

	return &root, nil
}
