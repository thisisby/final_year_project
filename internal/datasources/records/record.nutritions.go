package records

type Nutritions struct {
	Record
	OwnerID int    `db:"owner_id"`
	Name    string `db:"name"`
	Value   string `db:"value"`
}
