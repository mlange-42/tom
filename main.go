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
			location := strings.ToLower(strings.Join(args, " "))

			forceApi := strings.HasSuffix(location, "?")
			location = strings.TrimSuffix(location, "?")

			cached, err := config.LoadLocations()
			if err != nil {
				return err
			}

			coords, ok := cached[location]
			if ok && !forceApi {
				a := app.New(strings.ToTitle(location), coords)
				a.Run()
				return nil
			}

			a := app.NewLocationDialog(location, cached)
			if err := a.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	root.Flags().SortFlags = false

	return &root
}
