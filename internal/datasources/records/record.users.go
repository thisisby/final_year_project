package records

type Users struct {
	Record
	Email    string `db:"email"`
	Password string `db:"password"`
}
