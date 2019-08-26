package utils

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

func GetPartitionOffsets(broker string, topic string, partition int32) (startOffset int64, endOffset int64) {

	// get kafka consumer

	log("info", fmt.Sprintf("utils.kafka.GetPartitionOffsets broker %s topic %s partition %d",
		broker, topic, partition))

	consumer := GetKafkaConsumer(broker)

	if consumer == nil {
		log("error", fmt.Sprintf("utils.kafkaTopic.GetPartitionOffsets broker %s", broker))
	}

	defer consumer.Close()

	// get kafka partition_consumer

	partition_consumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
	if err != nil {
		log("error", fmt.Sprintf("utils.kafka.GetPartitionOffsets error %s", err.Error()))
	}

	defer partition_consumer.Close()

	// start set timeout
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
	// end set timeout

	log("info", fmt.Sprintf("utils.kafka.GetPartitionOffsets broker %s topic %s partition %d startOffset %d endOffset %d",
		broker, topic, partition, startOffset, endOffset))

	return startOffset, endOffset - 1
}

func GetPartitionMessages(params KafkaMessageQueryParam) []*KafkaMessage {
	messages := make([]*KafkaMessage, 0)

	log("info", fmt.Sprintf("utils.kafka.GetMessages broker %s topic %s offset %d limit %d order %s sort %s",
		params.Broker, params.Topic, params.Offset, params.Limit, params.Order, params.Sort))

	// iterator topic's all partitions

	for partition, partitionOffset := range params.PartitionOffsets {
		startOffset := partitionOffset["start"]
		endOffset := partitionOffset["end"]

		if startOffset < 0 || endOffset < 0 || params.Offset >= endOffset {
			continue
		}

		partition_consumer, err := params.Consumer.ConsumePartition(params.Topic, partition, params.Offset+1)
		if params.Offset+1 < startOffset {
			partition_consumer, err = params.Consumer.ConsumePartition(params.Topic, partition, sarama.OffsetOldest)
		}

		if err != nil {
			log("error", fmt.Sprintf("utils.kafka.GetMessages error: %s", err.Error()))
		}

		defer partition_consumer.Close()

		log("info", fmt.Sprintf("utils.kafka.GetMessages broker %s topic %s partition %d startOffset %d endOffset %d",
			params.Broker, params.Topic, partition, startOffset, endOffset))

		//get partition's messages
	LOOP:
		for {
			select {
			case msg := <-partition_consumer.Messages():
				// msg.offset out of range offset --> offset + limit
				if msg.Offset > params.Offset+int64(params.Limit) {
					break LOOP
				}

				messages = append(messages, &KafkaMessage{
					Offset:      msg.Offset,
					Partition:   msg.Partition,
					PublishTime: msg.Timestamp.Format("2006-01-02 15:04:05"),
					Data:        string(msg.Value),
				})

				// log("info", fmt.Sprintf("utils.kafka.GetMessages offset %d, partition %d, timestamp %s, value %s\n",
				//     msg.Offset, msg.Partition, msg.Timestamp.String(), string(msg.Value)))

				// last message
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
