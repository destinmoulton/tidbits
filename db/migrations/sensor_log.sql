-- migrate:up
CREATE TABLE IF NOT EXISTS sensor_log(
    sensor_id INTEGER,
    timestamp TEXT,
    sensor_reading NUMERIC
);

-- migrate:down
