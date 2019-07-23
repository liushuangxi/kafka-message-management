package utils

import (
	"sort"
	"strings"
)

const SEARCH_PAGE_SIZE = 1000

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
}

type KafkaMessage struct {
	Offset    int64
	Partition int32
	Data      string
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

		total = 0        //total return record
		totalRecord := 0 //total search record
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

		if int64(params.MaxReturnRecord) < total {
			total = int64(params.MaxReturnRecord)
		}

		messages = messages[0:total]

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
	}

	return messages, total
}
