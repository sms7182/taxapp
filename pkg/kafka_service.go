package pkg

type KafkaService interface {
	Consumer(msg string, err error)
	Publish(msg string) error
	Read(id string, callback func(string, error))
}
