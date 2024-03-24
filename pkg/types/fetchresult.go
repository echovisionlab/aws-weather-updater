package types

import "github.com/echovisionlab/aws-weather-updater/pkg/model"

type FetchResult interface {
	Stations() []model.Station
	Records() []model.Record
}

type fetchResult struct {
	records  []model.Record
	stations []model.Station
}

func (f *fetchResult) Stations() []model.Station {
	return f.stations
}

func (f *fetchResult) Records() []model.Record {
	return f.records
}

func NewFetchResult(stations []model.Station, records []model.Record) FetchResult {
	return &fetchResult{records, stations}
}
