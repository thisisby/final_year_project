package records

type Tokens struct {
	Record
	Token  string `db:"token"`
	UserID int    `db:"user_id"`
}
