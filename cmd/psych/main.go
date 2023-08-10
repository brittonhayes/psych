package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"sort"

	"os"

	"log/slog"

	"github.com/brittonhayes/therapy"
	"github.com/brittonhayes/therapy/scrape"
	"github.com/brittonhayes/therapy/sqlite"
	"github.com/brittonhayes/therapy/tui"
	"github.com/urfave/cli/v2"
)

func main() {
	var (
		repo therapy.Repository
	)

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	globalFlags := []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug logging",
			Value: false,
		},
		&cli.PathFlag{
			Name:  "config",
			Usage: "Config directory",
			Value: filepath.Join(cfg, "psych"),
		},
		&cli.StringFlag{
			Name:  "db",
			Usage: "Sqlite DB connection string",
			Value: fmt.Sprintf("file:%s", filepath.Join(cfg, "psych", "psych.db")),
		},
	}

	app := &cli.App{
		Name:                 "psych",
		Usage:                "Find therapists on psychologytoday.com",
		Suggest:              true,
		EnableBashCompletion: true,
		Flags:                globalFlags,
		Before: func(c *cli.Context) error {
			repo = sqlite.NewRepository(c.String("db"))
			if c.Bool("debug") {
				logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}))
			}

			err := repo.Init(context.Background())
			if err != nil {
				return err
			}

			if c.String("generate-migration") != "" {
				err = repo.Generate(context.Background(), c.String("generate-migration"))
				if err != nil {
					return err
				}
			}

			return repo.Migrate(context.Background())
		},
		Commands: []*cli.Command{
			{
				Name: "scrape",
				Flags: append(globalFlags,
					&cli.StringFlag{
						Name:  "state",
						Usage: "State to search",
						Value: "",
					},
					&cli.StringFlag{
						Name:  "city",
						Usage: "City to search",
						Value: "",
					},
					&cli.StringFlag{
						Name:  "zip",
						Usage: "Zip code to search",
						Value: "",
					},
					&cli.StringFlag{
						Name:  "county",
						Usage: "County to search",
						Value: "",
					},
				),
				Action: func(c *cli.Context) error {
					s := scrape.NewScraper(c.Context, logger, repo)

					url, err := buildURL(c.String("state"), c.String("county"), c.String("city"), c.String("zip"))
					if err != nil {
						return err
					}

					config := scrape.Config{URL: url, CacheDir: filepath.Join(c.String("config"), "cache/")}
					therapists := s.Scrape(config)
					logger.DebugContext(c.Context, "Found therapists", slog.Int("count", len(therapists)))

					tui.Run(therapists)

					return nil
				},
			},
		}}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func buildURL(state string, county string, city string, zip string) (string, error) {
	base := "https://www.psychologytoday.com/us/therapists/"

	if zip != "" {
		return url.JoinPath(base, zip)
	}

	if state != "" && county != "" {
		return url.JoinPath(base, state, county)
	}

	if state != "" && city != "" {
		return url.JoinPath(base, state, city)
	}

	return "", errors.New("not enough flags provided to generate web scraping URL")
}
