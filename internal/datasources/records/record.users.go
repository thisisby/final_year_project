package records

type Users struct {
	Record
	Email    string `db:"email"`
	Password string `db:"password"`
	Username string `db:"username"`
	Bio      string `db:"bio"`
	Avatar   string `db:"avatar"`
	CardPAN  string `db:"card_pan"`
}
