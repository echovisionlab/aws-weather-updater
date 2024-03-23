package query

const (
	UpsertStationQuery = "INSERT INTO realtime_weather_station (id, name, altitude, has_rain_sensor, address) " +
		"VALUES (:id, :name, :altitude, :has_rain_sensor, :address) " +
		"ON CONFLICT (id) " +
		"DO UPDATE SET name = excluded.name, altitude = excluded.altitude, has_rain_sensor = excluded.has_rain_sensor, address = excluded.address"
	UpsertRecordQuery = "INSERT INTO realtime_weather_record(id, rain_acc, rain_fifteen, rain_hour, rain_three_hour, rain_six_hour, rain_twelve_hour, temperature, wind_avg_minute, wind_avg_minute_deg, wind_avg_ten_minute, wind_avg_ten_minute_deg, moisture, sea_level_air_pressure, station_id, time) " +
		"VALUES (:id, :rain_acc, :rain_fifteen, :rain_hour, :rain_three_hour, :rain_six_hour, :rain_twelve_hour, :temperature, :wind_avg_minute, :wind_avg_minute_deg, :wind_avg_ten_minute, :wind_avg_ten_minute_deg, :moisture, :sea_level_air_pressure, :station_id, :time) " +
		"ON CONFLICT (time, station_id) " +
		"DO UPDATE SET id=excluded.id, rain_acc=excluded.rain_acc, rain_fifteen=excluded.rain_fifteen, rain_hour=excluded.rain_hour, rain_three_hour=excluded.rain_three_hour, rain_six_hour=excluded.rain_six_hour, rain_twelve_hour=excluded.rain_twelve_hour, temperature=excluded.temperature, wind_avg_minute=excluded.wind_avg_minute, wind_avg_minute_deg=excluded.wind_avg_minute_deg, wind_avg_ten_minute=excluded.wind_avg_ten_minute, wind_avg_ten_minute_deg=excluded.wind_avg_ten_minute_deg, moisture=excluded.moisture, sea_level_air_pressure=excluded.sea_level_air_pressure"
)
