package records

type SessionDetails struct {
	Record
	SessionID int    `db:"session_id"`
	Name      string `db:"name"`
	Value     string `db:"value"`
}
