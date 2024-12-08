package records

type Activities struct {
	Record
	Name            string `db:"name"`
	ActivityGroupID int    `db:"activity_group_id"`
}
