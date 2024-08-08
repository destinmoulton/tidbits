-- name: GetSensor :one
SELECT * FROM sensors
WHERE id = ? LIMIT 1;

-- name: GetSensorsBySource :many
SELECT * FROM sensors
WHERE sensor_source = ?
ORDER BY sensor_order ASC;

-- name: ListSensors :many
SELECT * FROM sensors
ORDER BY sensor_order ASC;

-- name: CreateSensor :one
INSERT INTO sensors (
    sensor_name,
    sensor_type,
    sensor_device,
    sensor_source,
    user_label
) VALUES (
    ?, ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateSensor :exec
UPDATE sensors
set user_label = ?,
    user_units = ?,
    should_log = ?,
    order = ?
WHERE id = ?;

-- name: DeleteSensor :exec
DELETE FROM sensors
WHERE id = ?;
