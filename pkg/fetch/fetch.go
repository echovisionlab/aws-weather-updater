package fetch

import (
	"context"
	"fmt"
	"github.com/echovisionlab/aws-weather-updater/pkg/model"
	"github.com/echovisionlab/aws-weather-updater/pkg/types"
	"github.com/go-rod/rod"
	"github.com/google/uuid"
	"strings"
	"time"
)

var (
	measurementBaseUrl = "https://www.weather.go.kr/cgi-bin/aws/nph-aws_txt_min_guide_test"
	cssSelector        = "body > table:nth-child(2) > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1) > table:nth-child(1) > tbody > tr"
)

func StationsAndRecords(ctx context.Context, p *rod.Page, t time.Time, retry int) (types.FetchResult, error) {
	strTime := t.Format("200601021504")
	targetUrl := measurementBaseUrl + "?" + strTime
	rch := make(chan types.FetchResult)
	ech := make(chan error)

	go func() {
		count := 0
	reset:
		if res, err := doFetchRecords(p, targetUrl); err != nil {
			if strings.Contains(err.Error(), "net::") && count < retry {
				count++
				goto reset
			} else {
				ech <- err
			}
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

func doFetchRecords(p *rod.Page, targetUrl string) (types.FetchResult, error) {
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

	s, r, err := parse(rows, t)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}

	return types.NewFetchResult(s, r), err
}

func tryNavigate(p *rod.Page, url string) error {
	if err := p.Navigate(url); err != nil {
		return fmt.Errorf("failed to navigate: %w", err)
	}
	if err := p.WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait loading page: %w", err)
	}
	return nil
}

func getObservationTime(p *rod.Page) (time.Time, error) {
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

func parse(rows rod.Elements, at time.Time) ([]model.Station, []model.Record, error) {
	idx := 0
	size := len(rows) - 1
	stations, records := make([]model.Station, size), make([]model.Record, size)
	for _, row := range rows {
		// skips initial row
		if className, err := row.Attribute("class"); err != nil {
			return nil, nil, err
		} else if className == nil || *className == "name" {
			continue
		}

		cols, err := row.Elements("td")
		if err != nil {
			return nil, nil, err
		}

		// skips invalid row
		if len(cols) < 20 {
			continue
		}

		station, record := ParseWeatherData(cols)
		record.Id = uuid.New()
		record.Time = at
		stations[idx] = station
		records[idx] = record
		idx++
	}
	return stations, records, nil
}
