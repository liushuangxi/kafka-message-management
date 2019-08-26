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

	config := sarama.NewConfig()
	config.Version = sarama.V0_10_0_0

	consumerInstance, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		log("error", fmt.Sprintf("utils.kafkaConsumer.GetKafkaConsumer error %s", err.Error()))
	}

	return consumerInstance
}
