package fetch

import (
	"context"

	"log/slog"

	"github.com/brittonhayes/therapy"
	"github.com/brittonhayes/therapy/api"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
)

type Fetcher interface {
	Fetch(config Config) []api.Therapist
}

type fetcher struct {
	ctx    context.Context
	repo   therapy.Repository
	logger *slog.Logger
}

type Config struct {
	CacheDir string
	URL      string
}

func NewFetcher(ctx context.Context, logger *slog.Logger, repo therapy.Repository) Fetcher {
	return &fetcher{
		ctx:    ctx,
		repo:   repo,
		logger: logger,
	}
}

func (s *fetcher) Fetch(config Config) []api.Therapist {

	therapists := []api.Therapist{}

	c := colly.NewCollector(
		colly.AllowedDomains("psychologytoday.com", "www.psychologytoday.com"),
		colly.CacheDir(config.CacheDir),
		colly.ParseHTTPErrorResponse(),
	)

	q, err := queue.New(1, &queue.InMemoryQueueStorage{MaxSize: 10000})
	if err != nil {
		panic(err)
	}

	q.AddURL(config.URL)

	c.OnHTML(".results-row", func(e *colly.HTMLElement) {
		var therapist api.Therapist

		s.logger.DebugContext(s.ctx, "scraping therapist", slog.String("name", e.ChildText(".results-row-name")))

		e.ForEach(".results-row-info", func(i int, e *colly.HTMLElement) {
			therapist.Title = e.ChildText(".profile-title")
			therapist.Credentials = e.ChildText(".profile-subtitle-credentials")
			therapist.Verified = e.ChildText(".verified-badge .profile-subtitle-badge .not-small")
			therapist.Statement = e.ChildText(".statements")
			therapist.Link = e.ChildAttr("a", "href")
		})

		e.ForEach(".profile-features", func(i int, e *colly.HTMLElement) {
			therapist.AcceptingAppointments = e.ChildText(".accepting-appointments")
		})

		e.ForEach(".results-row-contact", func(i int, e *colly.HTMLElement) {
			therapist.Phone = e.ChildText(".results-row-mob")
		})

		therapists = append(therapists, therapist)
	})

	c.OnHTML(".pagination", func(e *colly.HTMLElement) {
		hrefs := e.ChildAttrs("a[href].button-element.page-btn", "href")
		for _, r := range hrefs {
			q.AddURL(r)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		s.logger.DebugContext(s.ctx, "requesting url", slog.String("url", r.URL.String()))
	})

	c.OnError(func(r *colly.Response, err error) {
		s.logger.ErrorContext(s.ctx, "fetcher encountered error", err)
		s.logger.DebugContext(s.ctx, "error at url", slog.String("url", r.Request.URL.String()))
	})

	err = q.Run(c)
	if err != nil {
		panic(err)
	}

	if q.IsEmpty() {
		s.logger.DebugContext(s.ctx, "no more pages to scrape")
	}

	return therapists
}
