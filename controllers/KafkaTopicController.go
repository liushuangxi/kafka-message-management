package controllers

import (
	"encoding/json"
	"log"
	"sort"
	"strings"

	"github.com/liushuangxi/kafka-message-management/utils"
	"github.com/liushuangxi/kafka-message-management/models"
)

type KafkaTopicController struct {
	BaseController
}

func (c *KafkaTopicController) Prepare() {
	c.BaseController.Prepare()
	c.checkAuthor("DataGrid")
}

func (c *KafkaTopicController) Index() {
	broker := c.GetString("broker")

	c.Data["currentBroker"] = broker

	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl("kafka_topic/index.html")

	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "kafka_topic/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "kafka_topic/index_footerjs.html"

	c.Data["canEdit"] = c.checkActionAuthor("KafkaTopicController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("KafkaTopicController", "Delete")
}

func (c *KafkaTopicController) DataGrid() {
	var params utils.KafkaTopicQueryParam

	json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	offset := int64(params.Offset)
	limit := int64(params.Limit)

	topics, _ := utils.GetTopics(params.Broker)

	sort.Strings(topics)

	total := int64(len(topics))

	topicAll := make([]interface{}, total)

	log.Printf("%s", params.Topic)

	paramsCollect := models.KafkaTopicCollectQueryParam{}
	paramsCollect.UserId = int64(c.curUser.Id)
	collectList, _ := models.KafkaTopicCollectPageList(&paramsCollect)

	var collectMap map[string]int
	collectMap = make(map[string]int)
	for _, collect := range collectList {
		collectMap[collect.Topic] = 1
	}

	total = 0
	for i := 0; i < len(topics); i++ {
		if strings.Index(topics[i], params.Topic) < 0 {
			continue
		}

		topic := make(map[string]string)
		topic["Topic"] = topics[i]

		_, ok := collectMap[topics[i]]

		if ok {
			topic["Collect"] = "1"
		} else {
			topic["Collect"] = "0"
		}

		topicAll[total] = topic
		total++
	}

	maxOffset := offset + limit
	if maxOffset >= total {
		maxOffset = total
	}

	topicRet := make([]interface{}, limit)

	if offset < total {
		topicRet = topicAll[offset:maxOffset]
	}

	result := make(map[string]interface{})

	result["total"] = total
	result["rows"] = topicRet

	c.Data["json"] = result

	c.ServeJSON()
}
