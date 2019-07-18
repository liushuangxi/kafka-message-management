package controllers

import (
	"encoding/json"

	"github.com/liushuangxi/kafka-message-management/utils"
)

type KafkaMessageController struct {
	BaseController
}

func (c *KafkaMessageController) Prepare() {
	c.BaseController.Prepare()
	c.checkAuthor("DataGrid")
}

func (c *KafkaMessageController) Index() {
	broker := c.GetString("broker")
	topic := c.GetString("topic")

	c.Data["currentBroker"] = broker
	c.Data["currentTopic"] = topic

	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl("kafka_message/index.html")

	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "kafka_message/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "kafka_message/index_footerjs.html"

	c.Data["canEdit"] = c.checkActionAuthor("KafkaMessageController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("KafkaMessageController", "Delete")
}

func (c *KafkaMessageController) DataGrid() {
	var params utils.KafkaMessageQueryParam

	json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	messages, total := utils.GetMessages(params)

	result := make(map[string]interface{})

	result["total"] = total
	result["rows"] = messages

	c.Data["json"] = result

	c.ServeJSON()
}