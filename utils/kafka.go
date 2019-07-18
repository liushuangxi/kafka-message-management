package utils

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

const SEARCH_PAGE_SIZE = 1000;

type KafkaTopicQueryParam struct {
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Offset int64  `json:"offset"`
	Limit  int    `json:"limit"`
	Broker string `json:"Broker"`
	Topic  string `json:"Topic"`
}

type KafkaMessageQueryParam struct {
	Sort            string `json:"sort"`
	Order           string `json:"order"`
	Offset          int64  `json:"offset"`
	Limit           int    `json:"limit"`
	Broker          string `json:"Broker"`
	Topic           string `json:"Topic"`
	Message         string `json:"Message"`
	MaxReturnRecord int `json:"MaxReturnRecord"`
	MaxSearchRecord int `json:"MaxSearchRecord"`
}

type KafkaMessage struct {
	Offset    int64
	Partition int32
	Data      string
}

func GetTopics(broker string) ([]string, error) {
	log("info", fmt.Sprintf("utils.kafka.GetTopics broker %s\n", broker))

	consumer, err := sarama.NewConsumer([]string{broker}, nil)
	if err != nil {
		log("error", fmt.Sprintf("utils.kafka.GetTopics error %s\n", err.Error()))
	}

	defer consumer.Close()

	return consumer.Topics()
}

func GetMessages(params KafkaMessageQueryParam) ([]*KafkaMessage, int64) {
	originOffset := params.Offset
	originLimit := params.Limit

	total, minOffset, maxOffset := GetTopicTotalMessage(params.Broker, params.Topic)

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

	if len(params.Message) <= 0 {
		messages = GetPartitionMessages(params)

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
	} else {
		// search messages

		total = 0         //total return record
		totalRecord := 0  //total search record
		params.Limit = SEARCH_PAGE_SIZE
		if params.Order == "asc" {
			params.Offset = minOffset
		}
		for {
			tmps := GetPartitionMessages(params)

			for _, tmp := range tmps {
				totalRecord++

				if strings.Index(tmp.Data, params.Message) <= -1 {
					continue
				}

				messages = append(messages, tmp)

				total++
			}

			if total >= int64(params.MaxReturnRecord) {
				break
			}

			if totalRecord > params.MaxSearchRecord {
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

		total = int64(params.MaxReturnRecord)

		messages = messages[0:total]

		// set startIndex endIndex

		startIndex := int64(0)
		endIndex := int64(0)
		if int64(len(messages)) >= originOffset {
			startIndex = originOffset
			if int64(len(messages)) <= originOffset + int64(originLimit) {
				endIndex = int64(len(messages))
			} else {
				endIndex = originOffset + int64(originLimit)
			}
		}

		if int64(len(messages)) > 0 {
			messages = messages[startIndex:endIndex]
		}
	}

	return messages, total
}

func GetPartitionMessages(params KafkaMessageQueryParam) []*KafkaMessage {
	messages := make([]*KafkaMessage, 0)

	//create consumer

	log("info", fmt.Sprintf("utils.kafka.GetMessages broker %s topic %s offset %d limit %d order %s sort %s",
		params.Broker, params.Topic, params.Offset, params.Limit, params.Order, params.Sort))

	consumer, err := sarama.NewConsumer([]string{params.Broker}, nil)
	if err != nil {
		log("error", fmt.Sprintf("utils.kafka.GetMessages error %s", err.Error()))
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

		partition_consumer, err := consumer.ConsumePartition(params.Topic, int32(partition), params.Offset + 1)
		if params.Offset + 1 < startOffset {
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
				if msg.Offset > params.Offset + int64(params.Limit) {
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
				log("error", fmt.Sprintf("utils.kafka.GetMessages error %s", err.Error()))
			}
		}
	}

	return messages
}

func GetTopicTotalMessage(broker string, topic string) (total int64, minOffset int64, maxOffset int64) {
	minOffset = int64(-1)
	maxOffset = int64(-1)

	//create consumer

	log("info", fmt.Sprintf("utils.kafka.GetTopicTotalMessage broker %s topic %s",
		broker, topic))

	consumer, err := sarama.NewConsumer([]string{broker}, nil)
	if err != nil {
		log("error", fmt.Sprintf("utils.kafka.GetTopicTotalMessage error %s", err.Error()))
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

func GetPartitionOffsets(broker string, topic string, partition int32) (startOffset int64, endOffset int64) {

	//create consumer

	log("info", fmt.Sprintf("utils.kafka.GetPartitionOffsets broker %s topic %s partition %d",
		broker, topic, partition))

	consumer, err := sarama.NewConsumer([]string{broker}, nil)
	if err != nil {
		log("error", fmt.Sprintf("utils.kafka.GetPartitionOffsets error %s", err.Error()))
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
