package person

type Person struct {
	Name string `db:"name" json:"name"`
	Age  int    `db:"age" json:"age"`
}
