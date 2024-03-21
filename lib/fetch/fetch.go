package fetch

import (
	"context"
	"fmt"
	"github.com/echovisionlab/aws-weather-updater/lib/model"
	"github.com/echovisionlab/aws-weather-updater/lib/parser"
	"github.com/echovisionlab/aws-weather-updater/lib/types/fetchresult"
	"github.com/echovisionlab/chromium"
	"github.com/go-rod/rod"
	"strings"
	"time"
)

const (
	measurementBaseUrl = "https://www.weather.go.kr/cgi-bin/aws/nph-aws_txt_min_guide_test"
	cssSelector        = "body > table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > table:nth-child(1) > tbody > tr"
)

func browser() *chromium.Browser {
	return chromium.NewBrowser(chromium.NewBrowserOption(1, 0, nil))
}

func StationsAndRecords(ctx context.Context, t time.Time) (fetchresult.FetchResult, error) {
	b := browser()
	p := b.GetPage()

	defer func() {
		b.PutPage(p)
		b.CleanUp()
	}()

	strTime := t.Format("200601021504")
	targetUrl := measurementBaseUrl + "?" + strTime
	rch := make(chan fetchresult.FetchResult)
	ech := make(chan error)

	go func() {
		if res, err := doFetchRecords(p, targetUrl); err != nil {
			ech <- err
		} else {
			rch <- res
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-rch:
		return res, nil
	case err := <-ech:
		return nil, err
	}
}

func doFetchRecords(p *chromium.Page, targetUrl string) (fetchresult.FetchResult, error) {
	if err := tryNavigate(p, targetUrl); err != nil {
		return nil, fmt.Errorf("failed to fetch records: %w", err)
	}

	t, err := getObservationTime(p)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch records: %w", err)
	}

	rows, err := p.Elements(cssSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch records: %w", err)
	}

	return fetchresult.New(parse(rows, t)), nil
}

func tryNavigate(p *chromium.Page, url string) error {
	if err := p.TryNavigate(url); err != nil {
		return fmt.Errorf("failed to navigate: %w", err)
	}
	if err := p.WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait loading page: %w", err)
	}
	return nil
}

func getObservationTime(p *chromium.Page) (time.Time, error) {
	elem, err := p.Element(".ehead")
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to fetch observation time: %w", err)
	}
	text, err := elem.Text()
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to fetch observation time: %w", err)
	}
	dateStr := strings.Replace(text, "[ 매분관측자료 ] ", "", 1)
	return time.Parse("2006.01.02.15:04", dateStr)
}

func parse(rows rod.Elements, at time.Time) ([]model.Station, []model.Record) {
	stations, records := make([]model.Station, 0), make([]model.Record, 0)
	for _, row := range rows {
		// skips initial row
		if className := row.MustAttribute("class"); className == nil || *className == "name" {
			continue
		}

		cols := row.MustElements("td")
		// skips invalid row
		if len(cols) < 20 {
			continue
		}

		station, record := parser.ParseWeatherData(cols)
		record.Time = at
		stations = append(stations, station)
		records = append(records, record)
	}
	return stations, records
}
