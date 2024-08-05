CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE sensor_log(
    sensor_id INTEGER,
    timestamp TEXT,
    sensor_reading NUMERIC
);
CREATE TABLE sensors(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sensor_name TEXT,
    sensor_type TEXT,
    sensor_device TEXT,
    user_label TEXT,
    user_units TEXT,
    should_log INTEGER
);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('001'),
  ('002'),
  ('003');
