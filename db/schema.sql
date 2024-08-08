CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    value TEXT
);
CREATE TABLE sensors(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sensor_name TEXT,
    sensor_type TEXT,
    sensor_device TEXT,
    sensor_source TEXT,
    user_label TEXT,
    user_units TEXT,
    should_log INTEGER,
    sensor_order INTEGER
);
CREATE TABLE sensor_log(
    sensor_id INTEGER,
    timestamp TEXT,
    sensor_reading NUMERIC
);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('001'),
  ('002'),
  ('003');
