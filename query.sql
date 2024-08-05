-- name: GetSensor :one
SELECT * FROM sensors
WHERE id = ? LIMIT 1;

-- name: ListSensors :many
SELECT * FROM sensors
ORDER BY sensor_name;

-- name: CreateSensor :one
INSERT INTO sensors (
    sensor_name,
    sensor_type,
    sensor_device,
    user_label,
    user_units,
    should_log
) VALUES (
    ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateSensor :exec
UPDATE sensors
set user_label = ?,
    user_units = ?,
    should_log = ?
WHERE id = ?;

-- name: DeleteSensor :exec
DELETE FROM sensors
WHERE id = ?;
