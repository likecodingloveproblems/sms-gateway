package gateway

type Service struct {
	repo Repository
}

func (s Service) CalculateMessagePrice(message Message) uint64 {
	// It's assumed that all messages has one price
	return s.repo.GetMessageUnitPrice()
}

func (s Service) ScheduleMessage(message Message) error {
	switch message.Type {
	case NormalMessage:
		Store
	case ExpressMessage:
	default:
	}
	return nil
}
