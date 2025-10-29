package ttt

type z struct {
	ID   string
	Name string
}

func NewZ(id string, name string) *z {
	return &z{
		ID:   id,
		Name: name,
	}
}
