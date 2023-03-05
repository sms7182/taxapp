package pkg

type KafkaService interface {
	Consumer(msg string, err error)
	Read(id string, callback func(string, error))
}
