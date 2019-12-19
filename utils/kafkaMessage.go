package utils

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

type KafkaMessageQueryParam struct {
	Sort            string `json:"sort"`
	Order           string `json:"order"`
	Offset          int64  `json:"offset"`
	OriginOffset    int64  `json:"originOffset"`
	Limit           int    `json:"limit"`
	Broker          string `json:"Broker"`
	Topic           string `json:"Topic"`
	TimeStr         string `json:"TimeStr"`
	Message         string `json:"Message"`
	MaxReturnRecord int    `json:"MaxReturnRecord"`
	MaxSearchRecord int    `json:"MaxSearchRecord"`

	// internal param
	StartTime int64
	EndTime   int64

	// internal used
	PartitionOffsets map[int32]map[string]int64
	Consumer         sarama.Consumer
	Client           sarama.Client
}

type KafkaMessage struct {
	Offset      int64
	Partition   int32
	PublishTime string
	Data        string
}

func GetMessages(params KafkaMessageQueryParam) ([]*KafkaMessage, int64) {
	// get kafka consumer

	params.Consumer = GetKafkaConsumer(params.Broker)
	if params.Consumer == nil {
		log("error", fmt.Sprintf("utils.kafkaTopic.GetPartitionMessages broker %s topic %s", params.Broker, params.Topic))
	}

	defer params.Consumer.Close()

	params.StartTime, params.EndTime = ParseTimeStr(params.TimeStr)

	// search topic message

	if len(params.Message) > 0 {
		return GetMessagesBySearch(params)
	}

	total, minOffset, maxOffset, partitionOffsets := GetTopicTotalMessage(params.Broker, params.Topic)
	params.PartitionOffsets = partitionOffsets

	// set partition's start offset

	params.OriginOffset = params.Offset

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

func ParseTimeStr(TimeStr string) (int64, int64) {
	if TimeStr == "" {
		return 0, 0
	}

	TimeArr := strings.Split(TimeStr, " - ")

	timeFormat := "2006-01-02 15:04:05"

	location, _ := time.LoadLocation("Local")

	startTime, _ := time.ParseInLocation(timeFormat, TimeArr[0], location)
	endTime, _ := time.ParseInLocation(timeFormat, TimeArr[1], location)

	return startTime.Unix(), endTime.Unix()
}
