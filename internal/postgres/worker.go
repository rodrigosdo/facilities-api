package postgres

import (
	"context"
	"errors"

	"github.com/rodrigosdo/facilities-api/internal/cursor"
	"github.com/rodrigosdo/facilities-api/internal/domain"

	"cloud.google.com/go/civil"
	"github.com/jackc/pgx/v5"
)

func (d *Database) GetAvailableShifts(ctx context.Context, queryCursor *cursor.Cursor, limit int, workerID int64, start civil.Date, end civil.Date) (domain.Shifts, error) {
	shiftsBuilder := d.sq.
		Select(
			"s.start AS rounded_start",
			"s.end AS rounded_end",
			"s.id AS shift_id",
			"f.id AS facility_id",
			"f.name AS facility_name",
		).
		From("\"Shift\" s").
		InnerJoin("\"Facility\" f ON f.id = s.facility_id").
		InnerJoin("\"Worker\" w ON w.profession = s.profession").
		InnerJoin("facility_documents fd ON fd.facility_id = f.id").
		InnerJoin("worker_documents wd ON wd.worker_id = w.id AND wd.worker_documents @> fd.required_documents").
		Where("f.is_active = TRUE").
		Where("s.is_deleted = FALSE").
		Where("w.is_active = TRUE").
		Where("s.worker_id IS NULL").
		Where("w.id = ?", workerID).
		OrderBy("rounded_start", "rounded_end", "facility_id")

	if !start.IsZero() && !end.IsZero() {
		shiftsBuilder = shiftsBuilder.
			Where("DATE_TRUNC('day', s.start) BETWEEN ? AND ?", start, end).
			Where("DATE_TRUNC('day', s.end) BETWEEN ? AND ?", start, end)
	}

	sqlBuilder := d.sq.
		Select(
			"facility_id",
			"facility_name",
			"shift_id",
			"rounded_start",
			"rounded_end",
		).
		PrefixExpr(
			d.sq.Select(
				"facility_id",
				"ARRAY_AGG(document_id) AS required_documents",
			).
				From("\"FacilityRequirement\"").
				GroupBy("facility_id").
				Prefix("WITH facility_documents AS (").Suffix("),"),
		).
		PrefixExpr(
			d.sq.Select(
				"worker_id",
				"ARRAY_AGG(document_id) AS worker_documents",
			).
				From("\"DocumentWorker\"").
				Where("worker_id = ?", workerID).
				GroupBy("worker_id").
				Prefix("worker_documents AS (").Suffix("),"),
		).
		PrefixExpr(
			shiftsBuilder.Prefix("rounded_shifts AS (").Suffix(")"),
		).
		From("rounded_shifts").
		Limit(uint64(limit))

	switch {
	case queryCursor == nil:
		sqlBuilder = sqlBuilder.OrderBy("shift_id ASC")
	case queryCursor.Direction == cursor.DirectionBefore:
		sqlBuilder = sqlBuilder.Where("shift_id < ? ", queryCursor.Reference).OrderBy("shift_id DESC")
	case queryCursor.Direction == cursor.DirectionAfter:
		sqlBuilder = sqlBuilder.Where("shift_id > ? ", queryCursor.Reference).OrderBy("shift_id ASC")
	default:
		return nil, errors.New("invalid queryCursor direction")
	}

	sql, args, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := d.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	shifts, err := deserializeWorkerAvailableShifts(rows)
	if err != nil {
		return nil, err
	}

	return shifts, nil
}

func deserializeWorkerAvailableShifts(rows pgx.Rows) (domain.Shifts, error) {
	defer rows.Close()

	shifts := domain.Shifts{}
	for rows.Next() {
		s := domain.Shift{}

		if err := rows.Scan(
			&s.Facility.ID,
			&s.Facility.Name,
			&s.ID,
			&s.Start,
			&s.End,
		); err != nil {
			if err == pgx.ErrNoRows {
				continue
			}

			return nil, err
		}

		shifts = append(shifts, s)
	}

	return shifts, nil
}
