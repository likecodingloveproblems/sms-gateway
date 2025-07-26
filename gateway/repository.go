package gateway

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r Repository) GetMessageUnitPrice() uint64 {
	return DefaultMessagePerUnitPrice
}
