package kafka

import (
	"encoding/json"
	"fmt"
	"tax-management/external"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
)

type SyncConsumer struct {
	Conn *kafka.Consumer
}

func (s *SyncConsumer) StartConsuming(topics []string, msgProcessor func(topicName string, transaction external.RawTransaction) error) {
	defer s.Conn.Close()
	err := s.Conn.SubscribeTopics(topics, nil)
	if err != nil {
		msg := fmt.Sprintf("failed to subscribe to topics: %+v msg, err:%+v", topics, err)
		log.Error(msg)
		panic(msg)
	}

	for true {
		ev := s.Conn.Poll(1000)
		switch e := ev.(type) {
		case *kafka.Message:
			log.Info(fmt.Sprintf("value: %+v", string(e.Value)))
			log.Info(fmt.Sprintf("whole message: %+v", e))
			var rawData external.RawTransaction
			err = json.Unmarshal(e.Value, &rawData)
			if err != nil {
				msg := fmt.Sprintf("failed to unmarshal msg, err:%+v", err)
				log.Error(msg)
				panic(msg)
			}
			if e := msgProcessor(*e.TopicPartition.Topic, rawData); e != nil {
				panic(fmt.Sprintf("failed to process message, rawData: %+v, err: %s", rawData, e.Error()))
			}
			if _, cErr := s.Conn.CommitMessage(e); cErr != nil {
				panic(fmt.Sprintf("failed to process message, rawData: %+v, err: %+v", rawData, cErr))
			}
		case kafka.Error:
			msg := fmt.Sprintf("kafka err:%+v", e.Error())
			log.Error(msg)
			panic(msg)
		default:

		}
	}
}
