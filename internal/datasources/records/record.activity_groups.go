package records

type ActivityGroups struct {
	Record
	Name        string `db:"name"`
	Description string `db:"description"`
}
