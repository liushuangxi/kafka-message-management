package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

func GetPartitionOffsets(broker string, topic string, partition int32) (startOffset int64, endOffset int64) {

	// get kafka consumer

	log("info", fmt.Sprintf("utils.kafka.GetPartitionOffsets broker %s topic %s partition %d",
		broker, topic, partition))

	consumer := GetKafkaConsumer(broker)

	if consumer == nil {
		log("error", fmt.Sprintf("utils.kafkaTopic.GetPartitionOffsets broker %s topic %s partition %d",
			broker, topic, partition))
	}

	defer consumer.Close()

	// get kafka partition_consumer

	partition_consumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
	if err != nil {
		log("error", fmt.Sprintf("utils.kafkaTopic.GetPartitionOffsets broker %s topic %s partition %d error %s",
			broker, topic, partition, err.Error()))
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

func GetPartitionOffsetsByTime(broker string, topic string, partition int32, time int64) int64 {
	client := GetKafkaClient(broker)

	offset, err := client.GetOffset(topic, partition, time*1000)
	if err != nil {
		log("error", fmt.Sprintf("utils.kafka.GetPartitionOffsetsByTime broker %s topic %s partition %d time %d error %s",
			broker, topic, partition, time, err.Error()))
	}

	log("info", fmt.Sprintf("utils.kafka.GetPartitionOffsetsByTime broker %s topic %s partition %d time %d offset %d",
		broker, topic, partition, time, offset))

	return offset
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

		if params.StartTime != params.EndTime {
			timeOffset := GetPartitionOffsetsByTime(params.Broker, params.Topic, partition, params.StartTime)
			if timeOffset <= -1 {
				continue
			}

			if timeOffset+params.OriginOffset < endOffset {
				params.Offset = timeOffset + params.OriginOffset
			} else {
				params.Offset = timeOffset
			}
		}

		partition_consumer, err := params.Consumer.ConsumePartition(params.Topic, partition, params.Offset+1)
		if params.Offset+1 < startOffset {
			partition_consumer, err = params.Consumer.ConsumePartition(params.Topic, partition, sarama.OffsetOldest)
		}

		if err != nil {
			log("error", fmt.Sprintf("utils.kafka.GetPartitionMessages broker %s topic %s partition %d error: %s",
				params.Broker, params.Topic, partition, err.Error()))
		}

		defer partition_consumer.Close()

		log("info", fmt.Sprintf("utils.kafka.GetPartitionMessages broker %s topic %s partition %d startOffset %d endOffset %d",
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

				message := KafkaMessage{
					Offset:      msg.Offset,
					Partition:   msg.Partition,
					PublishTime: msg.Timestamp.Format("2006-01-02<br>15:04:05")}

				var out bytes.Buffer
				err := json.Indent(&out, []byte(string(msg.Value)), "", "  ")
				if err != nil {
					message.Data = string(msg.Value)
				} else {
					message.Data = out.String()
				}

				messages = append(messages, &message)

				// log("info", fmt.Sprintf("utils.kafka.GetMessages offset %d, partition %d, timestamp %s, value %s\n",
				//     msg.Offset, msg.Partition, msg.Timestamp.String(), string(msg.Value)))

				// last message
				if msg.Offset == endOffset {
					break LOOP
				}
			case err := <-partition_consumer.Errors():
				log("error", fmt.Sprintf("utils.kafka.GetPartitionMessages broker %s topic %s partition %d error %s",
					params.Broker, params.Topic, partition, err.Error()))
			}
		}
	}

	return messages
}

func GetMessagesBySearch(params KafkaMessageQueryParam) ([]*KafkaMessage, int64) {
	originOffset := params.Offset
	originLimit := params.Limit
	totalReturn := 0 //total return record
	totalSearch := 0 //total search record

	_, minOffset, maxOffset, partitionOffsets := GetTopicTotalMessage(params.Broker, params.Topic)
	params.PartitionOffsets = partitionOffsets

	// set partition start offset

	if params.Sort == "Offset" && params.Order == "asc" {
		if params.Offset == 0 {
			params.Offset = minOffset
		} else {
			params.Offset = minOffset + params.Offset
		}
	} else {
		if params.Offset == 0 {
			params.Offset = maxOffset - int64(params.Limit)
		} else {
			params.Offset = maxOffset - int64(params.Limit) - params.Offset
		}
	}

	messages := make([]*KafkaMessage, 0)

	// set messages limit

	if params.MaxSearchRecord > 5000000 {
		params.Limit = 1000000
	} else if params.MaxSearchRecord > 500000 {
		params.Limit = 100000
	} else {
		params.Limit = 10000
	}

	if params.Order == "asc" {
		params.Offset = minOffset
	}

	for {
		tmps := GetPartitionMessages(params)

		for _, tmp := range tmps {
			totalSearch++

			if strings.Index(tmp.Data, params.Message) <= -1 {
				continue
			}

			messages = append(messages, tmp)

			totalReturn++
		}

		if totalReturn >= params.MaxReturnRecord {
			break
		}

		if totalSearch > params.MaxSearchRecord {
			break
		}

		if params.Offset < minOffset || params.Offset > maxOffset {
			break
		}

		if params.Order == "asc" {
			params.Offset = params.Offset + int64(params.Limit)
		} else {
			params.Offset = params.Offset - int64(params.Limit)
		}
	}

	// sort messages

	if params.Order == "asc" {
		sort.Slice(messages, func(i, j int) bool {
			return messages[i].Offset < messages[j].Offset
		})
	} else {
		sort.Slice(messages, func(i, j int) bool {
			return messages[i].Offset > messages[j].Offset
		})
	}

	if params.MaxReturnRecord < totalReturn {
		totalReturn = params.MaxReturnRecord
		messages = messages[0:totalReturn]
	}

	// set startIndex endIndex

	startIndex := int64(0)
	endIndex := int64(0)
	if int64(len(messages)) >= originOffset {
		startIndex = originOffset
		if int64(len(messages)) <= originOffset+int64(originLimit) {
			endIndex = int64(len(messages))
		} else {
			endIndex = originOffset + int64(originLimit)
		}
	}

	if int64(len(messages)) > 0 {
		messages = messages[startIndex:endIndex]
	}

	return messages, int64(totalReturn)
}
