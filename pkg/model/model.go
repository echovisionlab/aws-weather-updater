package model

import (
	"github.com/google/uuid"
	"time"
)

type Station struct {
	Id            int    `db:"id"`
	Name          string `db:"name"`
	Altitude      int    `db:"altitude"`
	HasRainSensor bool   `db:"has_rain_sensor"`
	Address       string `db:"address"`
}

type Record struct {
	Id                      uuid.UUID `db:"id"`
	StationID               int       `db:"station_id"`
	RainAcc                 float32   `db:"rain_acc"`
	RainFifteen             float32   `db:"rain_fifteen"`
	RainHour                float32   `db:"rain_hour"`
	RainThreeHour           float32   `db:"rain_three_hour"`
	RainSixHour             float32   `db:"rain_six_hour"`
	RainTwelveHour          float32   `db:"rain_twelve_hour"`
	Temperature             float32   `db:"temperature"`
	WindAverageMinute       float32   `db:"wind_avg_minute"`
	WindAverageMinuteDeg    float32   `db:"wind_avg_minute_deg"`
	WindAverageTenMinute    float32   `db:"wind_avg_ten_minute"`
	WindAverageTenMinuteDeg float32   `db:"wind_avg_ten_minute_deg"`
	Moisture                int       `db:"moisture"`
	SeaLevelAirPressure     float32   `db:"sea_level_air_pressure"`
	Time                    time.Time `db:"time"`
}
