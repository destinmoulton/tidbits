-- migrate:up
CREATE TABLE IF NOT EXISTS sensors(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sensor_name TEXT,
    sensor_type TEXT,
    sensor_device TEXT,
    sensor_source TEXT,
    user_label TEXT,
    user_units TEXT,
    should_log INTEGER,
    order INTEGER
);

-- migrate:down
