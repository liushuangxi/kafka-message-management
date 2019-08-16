package utils

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Shopify/sarama"
)

type KafkaMessageQueryParam struct {
	Sort            string `json:"sort"`
	Order           string `json:"order"`
	Offset          int64  `json:"offset"`
	Limit           int    `json:"limit"`
	Broker          string `json:"Broker"`
	Topic           string `json:"Topic"`
	Message         string `json:"Message"`
	MaxReturnRecord int    `json:"MaxReturnRecord"`
	MaxSearchRecord int    `json:"MaxSearchRecord"`

	// internal used
	PartitionOffsets map[int32]map[string]int64
	Consumer         sarama.Consumer
}

type KafkaMessage struct {
	Offset    int64
	Partition int32
	Data      string
}

func GetMessages(params KafkaMessageQueryParam) ([]*KafkaMessage, int64) {
	// get kafka consumer

	params.Consumer = GetKafkaConsumer(params.Broker)
	if params.Consumer == nil {
		log("error", fmt.Sprintf("utils.kafkaTopic.GetPartitionMessages broker %s topic %s", params.Broker, params.Topic))
	}

	defer params.Consumer.Close()

	// search topic message

	if len(params.Message) > 0 {
		return GetMessagesBySearch(params)
	}

	total, minOffset, maxOffset, partitionOffsets := GetTopicTotalMessage(params.Broker, params.Topic)
	params.PartitionOffsets = partitionOffsets

	// set partition's start offset

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

	messages := GetPartitionMessages(params)

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

	return messages, total
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
