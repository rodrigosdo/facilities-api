package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/rodrigosdo/facilities-api/internal/config"
	"github.com/rodrigosdo/facilities-api/internal/cursor"
	"github.com/rodrigosdo/facilities-api/internal/postgres"

	"cloud.google.com/go/civil"
	"github.com/stretchr/testify/assert"
)

func TestGetAvailableShifts(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	ctx := context.Background()

	// cfg, err := config.New()
	// assert.NoError(t, err)
	cfg := config.Config{
		Database: config.Database{
			DSN: "postgres://postgres:postgres@localhost:5432/postgres",
		},
	}

	database, err := postgres.NewDatabase(ctx, cfg.Database)
	assert.NoError(t, err)

	if err := database.Ping(ctx); err != nil {
		t.Skipf("no connection could be established at '%s'. skipping postgres tests", cfg.Database.DSN)
	}

	t.Run("should return no available shifts when given an inactive worker", func(t *testing.T) {
		t.Parallel()

		availableShifts, err := database.GetAvailableShifts(
			ctx,
			nil,
			10,
			1,
			civil.Date{},
			civil.Date{},
		)
		assert.NoError(t, err)
		assert.Empty(t, availableShifts)
	})

	t.Run("should return available shifts when given an active worker", func(t *testing.T) {
		t.Parallel()

		availableShifts, err := database.GetAvailableShifts(
			ctx,
			nil,
			10,
			101,
			civil.Date{},
			civil.Date{},
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, availableShifts)
	})

	t.Run("should return available shifts between given dates", func(t *testing.T) {
		t.Parallel()

		date := civil.Date{Year: 2023, Month: time.February, Day: 02}

		availableShifts, err := database.GetAvailableShifts(
			ctx,
			nil,
			10,
			101,
			date,
			date,
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, availableShifts)
		assert.Equal(t, date, civil.DateOf(availableShifts[0].Start))
		assert.Equal(t, date, civil.DateOf(availableShifts[0].End))
	})

	t.Run("should return one available shift when a limit of one is given", func(t *testing.T) {
		t.Parallel()

		availableShifts, err := database.GetAvailableShifts(
			ctx,
			nil,
			1,
			101,
			civil.Date{},
			civil.Date{},
		)
		assert.NoError(t, err)
		assert.Len(t, availableShifts, 1)
	})

	t.Run("should return available shifts when an after cursor is given", func(t *testing.T) {
		t.Parallel()

		availableShifts, err := database.GetAvailableShifts(
			ctx,
			&cursor.Cursor{Direction: cursor.DirectionAfter, Reference: "21"},
			10,
			101,
			civil.Date{},
			civil.Date{},
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, availableShifts)
	})

	t.Run("should return available shifts when an before cursor is given", func(t *testing.T) {
		t.Parallel()

		availableShifts, err := database.GetAvailableShifts(
			ctx,
			&cursor.Cursor{Direction: cursor.DirectionBefore, Reference: "21"},
			10,
			101,
			civil.Date{},
			civil.Date{},
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, availableShifts)
	})

	t.Run("should return an error if an invalid cursor is given", func(t *testing.T) {
		t.Parallel()

		const unknownCursorDirection = cursor.Direction("unknown")

		availableShifts, err := database.GetAvailableShifts(
			ctx,
			&cursor.Cursor{Direction: unknownCursorDirection, Reference: "21"},
			10,
			101,
			civil.Date{},
			civil.Date{},
		)
		assert.Error(t, err)
		assert.Nil(t, availableShifts)
	})
}
