package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func nullUUIDString(v uuid.NullUUID) string {
	if !v.Valid {
		return ""
	}
	return v.UUID.String()
}

func nullTimePtr(v sql.NullTime) *time.Time {
	if !v.Valid {
		return nil
	}
	t := v.Time
	return &t
}
