create table if not exists public.realtime_weather_station
(
    id              serial not null
        constraint pk__realtime_weather_station__id
            primary key,
    name            varchar(255),
    altitude        integer,
    has_rain_sensor boolean,
    address         varchar(1024)
);

create table if not exists public.realtime_weather_record
(
    id                      uuid default gen_random_uuid() not null
        constraint pk__realtime_weather_record__id
            primary key,
    rain_acc                double precision,
    rain_fifteen            double precision,
    rain_hour               double precision,
    rain_three_hour         double precision,
    rain_six_hour           double precision,
    rain_twelve_hour        double precision,
    temperature             double precision,
    wind_avg_minute         double precision,
    wind_avg_minute_deg     double precision,
    wind_avg_ten_minute     double precision,
    wind_avg_ten_minute_deg double precision,
    moisture                integer,
    sea_level_air_pressure  double precision,
    station_id              integer                        not null
        constraint fk__station_id__station_record_station_id
            references public.realtime_weather_station,
    time                    timestamp                      not null
);

create unique index if not exists uq__realtime_weather_record__time__station_id
    on public.realtime_weather_record (time, station_id);

