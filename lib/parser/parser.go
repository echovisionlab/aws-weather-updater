package parser

import (
	"github.com/echovisionlab/aws-weather-updater/lib/model"
	"github.com/go-rod/rod"
	"strconv"
	"strings"
)

// ParseWeatherData parses single row into the station and record.
func ParseWeatherData(columns rod.Elements) (l model.Station, r model.Record) {
	hasSensor := strings.TrimSpace(columns[3].MustText()) == "○"
	l = model.Station{
		Id:            parseInt(columns[0].MustText()),
		Name:          parseStr(columns[1].MustText()),
		Altitude:      parseInt(strings.ReplaceAll(columns[2].MustText(), "m", "")),
		HasRainSensor: hasSensor,
		Address:       parseStr(columns[19].MustText()),
	}
	r = model.Record{
		StationID:               l.Id,
		RainAcc:                 parseFloat32(columns[9].MustText()),
		RainFifteen:             parseFloat32(columns[4].MustText()),
		RainHour:                parseFloat32(columns[5].MustText()),
		RainThreeHour:           parseFloat32(columns[6].MustText()),
		RainSixHour:             parseFloat32(columns[7].MustText()),
		RainTwelveHour:          parseFloat32(columns[8].MustText()),
		Temperature:             parseFloat32(columns[10].MustText()),
		WindAverageMinute:       parseFloat32(columns[13].MustText()),
		WindAverageMinuteDeg:    parseFloat32(columns[11].MustText()),
		WindAverageTenMinute:    parseFloat32(columns[16].MustText()),
		WindAverageTenMinuteDeg: parseFloat32(columns[14].MustText()),
		Moisture:                parseInt(columns[17].MustText()),
		SeaLevelAirPressure:     parseFloat32(columns[18].MustText()),
	}
	return
}

func parseFloat32(s string) float32 {
	s = strings.TrimSpace(s)
	if s == "." {
		return -1
	}
	tmp, _ := strconv.ParseFloat(s, 32)
	return float32(tmp)
}

func parseStr(s string) string {
	s = strings.TrimSpace(s)
	if s == "." {
		return "None"
	} else {
		return s
	}
}

func parseInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "." {
		return -1
	}
	v, _ := strconv.Atoi(s)
	return v
}
