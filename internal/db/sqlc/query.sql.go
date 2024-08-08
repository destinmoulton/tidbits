// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package sqlc

import (
	"context"
	"database/sql"
)

const createSensor = `-- name: CreateSensor :one
INSERT INTO sensors (
    sensor_name,
    sensor_type,
    sensor_device,
    sensor_source,
    user_label
) VALUES (
    ?, ?, ?, ?, ?
)
RETURNING id, sensor_name, sensor_type, sensor_device, sensor_source, user_label, user_units, should_log, sensor_order
`

type CreateSensorParams struct {
	SensorName   sql.NullString
	SensorType   sql.NullString
	SensorDevice sql.NullString
	SensorSource sql.NullString
	UserLabel    sql.NullString
}

func (q *Queries) CreateSensor(ctx context.Context, arg CreateSensorParams) (Sensor, error) {
	row := q.db.QueryRowContext(ctx, createSensor,
		arg.SensorName,
		arg.SensorType,
		arg.SensorDevice,
		arg.SensorSource,
		arg.UserLabel,
	)
	var i Sensor
	err := row.Scan(
		&i.ID,
		&i.SensorName,
		&i.SensorType,
		&i.SensorDevice,
		&i.SensorSource,
		&i.UserLabel,
		&i.UserUnits,
		&i.ShouldLog,
		&i.SensorOrder,
	)
	return i, err
}

const deleteSensor = `-- name: DeleteSensor :exec
DELETE FROM sensors
WHERE id = ?
`

func (q *Queries) DeleteSensor(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteSensor, id)
	return err
}

const getSensor = `-- name: GetSensor :one
SELECT id, sensor_name, sensor_type, sensor_device, sensor_source, user_label, user_units, should_log, sensor_order FROM sensors
WHERE id = ? LIMIT 1
`

func (q *Queries) GetSensor(ctx context.Context, id int64) (Sensor, error) {
	row := q.db.QueryRowContext(ctx, getSensor, id)
	var i Sensor
	err := row.Scan(
		&i.ID,
		&i.SensorName,
		&i.SensorType,
		&i.SensorDevice,
		&i.SensorSource,
		&i.UserLabel,
		&i.UserUnits,
		&i.ShouldLog,
		&i.SensorOrder,
	)
	return i, err
}

const getSensorsBySource = `-- name: GetSensorsBySource :many
SELECT id, sensor_name, sensor_type, sensor_device, sensor_source, user_label, user_units, should_log, sensor_order FROM sensors
WHERE sensor_source = ?
ORDER BY sensor_order ASC
`

func (q *Queries) GetSensorsBySource(ctx context.Context, sensorSource sql.NullString) ([]Sensor, error) {
	rows, err := q.db.QueryContext(ctx, getSensorsBySource, sensorSource)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Sensor
	for rows.Next() {
		var i Sensor
		if err := rows.Scan(
			&i.ID,
			&i.SensorName,
			&i.SensorType,
			&i.SensorDevice,
			&i.SensorSource,
			&i.UserLabel,
			&i.UserUnits,
			&i.ShouldLog,
			&i.SensorOrder,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSensors = `-- name: ListSensors :many
SELECT id, sensor_name, sensor_type, sensor_device, sensor_source, user_label, user_units, should_log, sensor_order FROM sensors
ORDER BY sensor_order ASC
`

func (q *Queries) ListSensors(ctx context.Context) ([]Sensor, error) {
	rows, err := q.db.QueryContext(ctx, listSensors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Sensor
	for rows.Next() {
		var i Sensor
		if err := rows.Scan(
			&i.ID,
			&i.SensorName,
			&i.SensorType,
			&i.SensorDevice,
			&i.SensorSource,
			&i.UserLabel,
			&i.UserUnits,
			&i.ShouldLog,
			&i.SensorOrder,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateSensor = `-- name: UpdateSensor :exec
UPDATE sensors
set user_label = ?,
    user_units = ?,
    should_log = ?,
    sensor_order = ?
WHERE id = ?
`

type UpdateSensorParams struct {
	UserLabel   sql.NullString
	UserUnits   sql.NullString
	ShouldLog   sql.NullInt64
	SensorOrder sql.NullInt64
	ID          int64
}

func (q *Queries) UpdateSensor(ctx context.Context, arg UpdateSensorParams) error {
	_, err := q.db.ExecContext(ctx, updateSensor,
		arg.UserLabel,
		arg.UserUnits,
		arg.ShouldLog,
		arg.SensorOrder,
		arg.ID,
	)
	return err
}
