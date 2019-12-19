package utils

import (
	"fmt"

	"github.com/Shopify/sarama"
)

var clientInstance sarama.Client
var consumerInstance sarama.Consumer

func GetKafkaClient(broker string) sarama.Client {
	if clientInstance != nil {
		return clientInstance
	}

	log("info", fmt.Sprintf("utils.kafkaConsumer.GetKafkaClient broker %s", broker))

	config := sarama.NewConfig()
	config.Version = sarama.V0_11_0_0

	clientInstance, err := sarama.NewClient([]string{broker}, config)
	if err != nil {
		log("error", fmt.Sprintf("utils.kafkaConsumer.GetKafkaClient broker %s error %s", broker, err.Error()))
	}

	return clientInstance
}

func GetKafkaConsumer(broker string) sarama.Consumer {
	if consumerInstance != nil {
		return consumerInstance
	}

	log("info", fmt.Sprintf("utils.kafkaConsumer.GetKafkaConsumer broker %s", broker))

	config := sarama.NewConfig()
	config.Version = sarama.V0_10_0_0

	consumerInstance, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		log("error", fmt.Sprintf("utils.kafkaConsumer.GetKafkaConsumer broker %s error %s", broker, err.Error()))
	}

	return consumerInstance
}
