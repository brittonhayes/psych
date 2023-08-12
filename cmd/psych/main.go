package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
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
		repo   therapy.Repository
		logger *slog.Logger
	)

	output := os.Stderr

	level := new(slog.LevelVar)
	level.Set(slog.LevelInfo)
	logger = slog.New(slog.NewTextHandler(output, &slog.HandlerOptions{
		Level: level,
	}))

	cfg, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	globalFlags := []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "Enable verbose logging",
		},
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Config directory path",
			Value:   filepath.Join(cfg, "psych/"),
		},
		&cli.StringFlag{
			Name:  "db",
			Usage: "Sqlite DB file path",
			Value: fmt.Sprintf("file:%s", filepath.Join(cfg, "psych/", "psych.db")),
		},
	}

	app := &cli.App{
		Name:                   "psych",
		Usage:                  "Find therapists on psychologytoday.com",
		Suggest:                true,
		EnableBashCompletion:   true,
		UseShortOptionHandling: true,
		Flags:                  globalFlags,
		Before: func(c *cli.Context) error {
			if c.Bool("verbose") {
				level.Set(slog.LevelDebug)
			}

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "refresh",
				Usage: "Clear the cache and sqlite database",
				Flags: globalFlags,
				Action: func(c *cli.Context) error {
					if _, err := os.Stat(c.String("config")); err == nil {
						logger.DebugContext(c.Context, "deleting config directory", slog.String("path", c.String("config")))
						err := os.RemoveAll(c.String("config"))
						if err != nil {
							return err
						}
					}

					return nil
				},
			},
			{
				Name: "scrape",
				Flags: append(globalFlags,
					&cli.StringFlag{
						Name:     "state",
						Usage:    "State to search",
						Value:    "",
						Category: "Scraping",
					},
					&cli.StringFlag{
						Name:     "city",
						Usage:    "City to search",
						Value:    "",
						Category: "Scraping",
					},
					&cli.StringFlag{
						Name:     "zip",
						Usage:    "Zip code to search",
						Value:    "",
						Category: "Scraping",
					},
					&cli.StringFlag{
						Name:     "county",
						Usage:    "County to search",
						Value:    "",
						Category: "Scraping",
					},
				),
				Before: func(c *cli.Context) error {
					if _, err := os.Stat(c.String("config")); err != nil {
						err := os.MkdirAll(c.String("config"), fs.ModePerm)
						if err != nil {
							return err
						}
					}

					repo = sqlite.NewRepository(c.String("db"), logger)

					err := repo.Init(context.Background())
					if err != nil {
						return err
					}

					err = repo.Migrate(context.Background())
					if err != nil {
						return err
					}

					return nil
				},
				Action: func(c *cli.Context) error {

					url, err := buildURL(c.String("state"), c.String("county"), c.String("city"), c.String("zip"))
					if err != nil {
						return err
					}

					config := scrape.Config{URL: url, CacheDir: filepath.Join(c.String("config"), "cache/")}

					logger.InfoContext(c.Context, "Scraping psychologytoday.com for therapists")
					s := scrape.NewScraper(c.Context, logger, repo)
					therapists := s.Scrape(config)

					uniqueTherapists := map[string]therapy.Therapist{}
					for _, therapist := range therapists {
						uniqueTherapists[therapist.Title] = therapist
					}

					logger.InfoContext(c.Context, "Saving therapists to database")
					for k, v := range uniqueTherapists {
						logger.DebugContext(c.Context, "saving therapist", slog.String("title", k))
						err := repo.Save(c.Context, v)
						if err != nil {
							return err
						}
					}

					logger.InfoContext(c.Context, "Saved therapists to database", slog.Int("count", len(uniqueTherapists)))
					return nil
				},
			},
			{
				Name:  "browse",
				Usage: "Browse therapists in the terminal",
				Action: func(c *cli.Context) error {
					if repo == nil {
						repo = sqlite.NewRepository(c.String("db"), logger)
					}

					therapists, err := repo.List(c.Context)
					if err != nil {
						return err
					}

					if len(therapists) == 0 {
						return errors.New("no therapists found - please run the scrape command first")
					}

					return tui.Run(therapists)
				},
			},
		}}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		logger.Error(err.Error())
		os.
			Exit(1)
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
