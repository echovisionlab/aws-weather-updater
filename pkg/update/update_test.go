package update

import (
	"context"
	"fmt"
	"github.com/echovisionlab/aws-weather-updater/internal/testutil"
	"github.com/echovisionlab/aws-weather-updater/pkg/database"
	"github.com/echovisionlab/aws-weather-updater/pkg/model"
	"github.com/echovisionlab/aws-weather-updater/pkg/type/fetchresult"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	container := testutil.SetupPostgres(context.Background(), t)
	defer testutil.ShutdownContainer(context.Background(), t, container)

	stations := getStations(10)
	records := getRecords(stations)
	f := fetchresult.New(stations, records)
	db, err := database.NewDatabase()
	assert.NoError(t, err)

	t.Run("context canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		var wg sync.WaitGroup
		wg.Add(1)
		err = Run(ctx, &wg, db, f)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("context timed out", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
		defer cancel()
		var wg sync.WaitGroup
		wg.Add(1)
		err = Run(ctx, &wg, db, f)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("must update", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)
		err = Run(context.Background(), &wg, db, f)
		assert.NoError(t, err)

		var (
			qs []model.Station
			qr []model.Record
		)
		assert.NoError(t, db.SelectContext(context.Background(), &qs, "SELECT * FROM realtime_weather_station"))
		assert.NoError(t, db.SelectContext(context.Background(), &qr, "SELECT * FROM realtime_weather_record"))

		for _, station := range qs {
			assert.Contains(t, stations, station)
		}
		for _, record := range qr {
			record.Time = record.Time.In(time.UTC)
			assert.Contains(t, records, record)
		}
	})
}

func TestContainerUpdate(t *testing.T) {
	container := testutil.SetupPostgres(context.Background(), t)
	defer testutil.ShutdownContainer(context.Background(), t, container)

	db, err := database.NewDatabase()
	t.Cleanup(func() { _ = db.Close() })
	assert.NoError(t, err)

	var stations []model.Station

	t.Run("must update stations", func(t *testing.T) {
		_, err = db.Query("DELETE FROM realtime_weather_station WHERE TRUE")
		assert.NoError(t, err)
		size := rand.Intn(100) + 5
		stations = getStations(size)

		count, err := updateStations(context.Background(), db, stations)
		assert.NoError(t, err)
		assert.Equal(t, int64(size), count)

		var queriedStations []model.Station
		assert.NoError(t, db.SelectContext(context.Background(), &queriedStations, "SELECT * FROM realtime_weather_station"))

		for i := range size {
			assert.Equal(t, stations[i], queriedStations[i])
		}
	})

	t.Run("must update records", func(t *testing.T) {
		_, err = db.Query("DELETE FROM realtime_weather_record WHERE TRUE")
		assert.NoError(t, err)

		stationCnt := len(stations)

		records := getRecords(stations)
		updated, err := updateRecords(context.Background(), db, records)

		assert.NoError(t, err)
		assert.Equal(t, updated, int64(stationCnt))

		var queriedRecords []model.Record
		assert.NoError(t, db.SelectContext(context.Background(), &queriedRecords, "SELECT * FROM realtime_weather_record"))

		for i, qRec := range queriedRecords {
			qRec.Time = qRec.Time.In(time.UTC)
			assert.Equal(t, qRec, records[i])
		}
	})
}

func getRecords(stations []model.Station) []model.Record {
	stationCnt := len(stations)
	records := make([]model.Record, stationCnt)
	for i := 0; i < stationCnt; i++ {
		records[i] = getRecord(stations[i].Id)
		records[i].Time = records[i].Time.In(time.UTC).Truncate(time.Second)
	}
	return records
}

func getStations(size int) []model.Station {
	stations := make([]model.Station, size)
	for i := 0; i < size; i++ {
		stations[i] = getStation(i)
	}
	return stations
}

func getStation(id int) model.Station {
	v := fmt.Sprintf("test station %v", strconv.Itoa(id))
	return model.Station{
		Id:            id,
		Name:          v,
		Altitude:      rand.Intn(100),
		HasRainSensor: rand.Int()%2 > 0,
		Address:       v,
	}
}

func getRecord(stationID int) model.Record {
	return model.Record{
		Id:                      uuid.New(),
		StationID:               stationID,
		RainAcc:                 rand.Float32(),
		RainFifteen:             rand.Float32(),
		RainHour:                rand.Float32(),
		RainThreeHour:           rand.Float32(),
		RainSixHour:             rand.Float32(),
		RainTwelveHour:          rand.Float32(),
		Temperature:             rand.Float32(),
		WindAverageMinute:       rand.Float32(),
		WindAverageMinuteDeg:    rand.Float32(),
		WindAverageTenMinute:    rand.Float32(),
		WindAverageTenMinuteDeg: rand.Float32(),
		Moisture:                rand.Intn(10),
		SeaLevelAirPressure:     rand.Float32(),
		Time:                    time.Now(),
	}
}
