package models

type TransformerService struct {
	Name string
}

func CreateTransformerService(name string) *TransformerService {
	return &TransformerService{Name: name}
}
