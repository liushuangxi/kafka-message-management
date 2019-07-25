package utils

import (
	"fmt"
)

type KafkaTopicQueryParam struct {
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Offset int64  `json:"offset"`
	Limit  int    `json:"limit"`
	Broker string `json:"Broker"`
	Topic  string `json:"Topic"`
}

func GetTopics(broker string) ([]string, error) {
	log("info", fmt.Sprintf("utils.kafkaTopic.GetTopics broker %s\n", broker))

	consumer := GetKafkaConsumer(broker)

	if consumer == nil {
		log("error", fmt.Sprintf("utils.kafkaTopic.GetTopics broker %s", broker))
	}

	defer consumer.Close()

	return consumer.Topics()
}

func GetTopicTotalMessage(broker string, topic string) (total int64, minOffset int64, maxOffset int64) {
	minOffset = int64(-1)
	maxOffset = int64(-1)

	//create consumer

	log("info", fmt.Sprintf("utils.kafkaTopic.GetTopicTotalMessage broker %s topic %s", broker, topic))

	consumer := GetKafkaConsumer(broker)

	if consumer == nil {
		log("error", fmt.Sprintf("utils.kafkaTopic.GetTopicTotalMessage broker %s", broker))
	}

	defer consumer.Close()

	//topic all partitions

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log("error", fmt.Sprintf("utils.kafka.GetTopicTotalMessage error %s", err.Error()))
	}

	for partition := range partitions {
		startOffset, endOffset := GetPartitionOffsets(broker, topic, int32(partition))

		if startOffset >= endOffset {
			continue
		}

		if minOffset == -1 || startOffset < minOffset {
			minOffset = startOffset
		}

		if endOffset > maxOffset {
			maxOffset = endOffset
		}

		total += endOffset - startOffset + 1

		log("info", fmt.Sprintf("utils.kafka.GetTopicTotalMessage broker %s topic %s partition %d startOffset %d endOffset %d",
			broker, topic, partition, startOffset, endOffset))
	}

	log("info", fmt.Sprintf("utils.kafka.GetTopicTotalMessage broker %s topic %s total %d minOffset %d maxOffset %d",
		broker, topic, total, minOffset, maxOffset))

	return total, minOffset, maxOffset
}
