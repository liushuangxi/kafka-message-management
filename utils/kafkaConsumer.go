package utils

import (
	"fmt"

	"github.com/Shopify/sarama"
)

var consumerInstance sarama.Consumer

func GetKafkaConsumer(broker string) sarama.Consumer {
	if consumerInstance != nil {
		return consumerInstance
	}

	log("info", fmt.Sprintf("utils.kafkaConsumer.GetKafkaConsumer broker %s", broker))

	consumerInstance, err := sarama.NewConsumer([]string{broker}, nil)
	if err != nil {
		log("error", fmt.Sprintf("utils.kafkaConsumer.GetKafkaConsumer error %s", err.Error()))
	}

	return consumerInstance
}
