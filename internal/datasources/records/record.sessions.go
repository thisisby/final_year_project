package records

import "time"

type Sessions struct {
	Record
	Activity   Activities `db:"activity"`
	Notes      string     `db:"notes"`
	StartTime  time.Time  `db:"start_time"`
	EndTime    time.Time  `db:"end_time"`
	ActivityID int        `db:"activity_id"`
	OwnerID    int        `db:"owner_id"`
}
