CREATE TABLE IF NOT EXISTS settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    value TEXT
);

CREATE TABLE IF NOT EXISTS sensors(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sensor_name TEXT,
    sensor_type TEXT,
    sensor_device TEXT,
    user_label TEXT,
    user_units TEXT,
    should_log INTEGER
);

CREATE TABLE IF NOT EXISTS sensor_log(
    sensor_id INTEGER,
    timestamp TEXT,
    sensor_reading NUMERIC
);