package utils

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

func GetPartitionOffsets(broker string, topic string, partition int32) (startOffset int64, endOffset int64) {

	//create consumer

	log("info", fmt.Sprintf("utils.kafka.GetPartitionOffsets broker %s topic %s partition %d",
		broker, topic, partition))

	consumer := GetKafkaConsumer(broker)

	if consumer == nil {
		log("error", fmt.Sprintf("utils.kafkaTopic.GetPartitionOffsets broker %s", broker))
	}

	defer consumer.Close()

	//create partition_consumer

	partition_consumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
	if err != nil {
		log("error", fmt.Sprintf("utils.kafka.GetPartitionOffsets error %s", err.Error()))
	}

	defer partition_consumer.Close()

	//start set timeout
	ch := make(chan bool, 1)

	defer close(ch)

	go func() {
		ch <- true
	}()

	select {
	case startMessage := <-partition_consumer.Messages():
		endOffset = partition_consumer.HighWaterMarkOffset()
		startOffset = startMessage.Offset
	case <-time.After(500 * time.Millisecond):
		endOffset = partition_consumer.HighWaterMarkOffset()
		startOffset = endOffset - 1
	}
	//end set timeout

	log("info", fmt.Sprintf("utils.kafka.GetPartitionOffsets broker %s topic %s partition %d startOffset %d endOffset %d",
		broker, topic, partition, startOffset, endOffset))

	return startOffset, endOffset - 1
}

func GetPartitionMessages(params KafkaMessageQueryParam) []*KafkaMessage {
	messages := make([]*KafkaMessage, 0)

	//create consumer

	log("info", fmt.Sprintf("utils.kafka.GetMessages broker %s topic %s offset %d limit %d order %s sort %s",
		params.Broker, params.Topic, params.Offset, params.Limit, params.Order, params.Sort))

	consumer := GetKafkaConsumer(params.Broker)

	if consumer == nil {
		log("error", fmt.Sprintf("utils.kafkaTopic.GetPartitionMessages broker %s topic %s", params.Broker))
	}

	defer consumer.Close()

	//topic all partitions

	partitions, err := consumer.Partitions(params.Topic)
	if err != nil {
		log("error", fmt.Sprintf("utils.kafka.GetMessages error %s", err.Error()))
	}

	for partition := range partitions {
		//get partition offsets

		startOffset, endOffset := GetPartitionOffsets(params.Broker, params.Topic, int32(partition))
		if startOffset < 0 || endOffset < 0 || params.Offset >= endOffset {
			continue
		}

		//create partition_consumer

		partition_consumer, err := consumer.ConsumePartition(params.Topic, int32(partition), params.Offset+1)
		if params.Offset+1 < startOffset {
			partition_consumer, err = consumer.ConsumePartition(params.Topic, int32(partition), sarama.OffsetOldest)
		}

		if err != nil {
			log("error", fmt.Sprintf("utils.kafka.GetMessages error: %s", err.Error()))
		}

		defer partition_consumer.Close()

		log("info", fmt.Sprintf("utils.kafka.GetMessages broker %s topic %s partition %d startOffset %d endOffset %d",
			params.Broker, params.Topic, partition, startOffset, endOffset))

		//get partition messages
	LOOP:
		for {
			select {
			case msg := <-partition_consumer.Messages():
				//msg.offset out of range offset -- offset + limit
				if msg.Offset > params.Offset+int64(params.Limit) {
					break LOOP
				}

				messages = append(messages, &KafkaMessage{
					Offset:    msg.Offset,
					Partition: msg.Partition,
					Data:      string(msg.Value),
				})

				// log("info", fmt.Sprintf("utils.kafka.GetMessages offset %d, partition %d, timestamp %s, value %s\n",
				//     msg.Offset, msg.Partition, msg.Timestamp.String(), string(msg.Value)))

				//last message
				if msg.Offset == endOffset {
					break LOOP
				}
			case err := <-partition_consumer.Errors():
				log("error", fmt.Sprintf("utils.kafka.` error %s", err.Error()))
			}
		}
	}

	return messages
}
