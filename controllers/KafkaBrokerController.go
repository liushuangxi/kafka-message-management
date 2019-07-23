package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"

	"github.com/liushuangxi/kafka-message-management/enums"
	"github.com/liushuangxi/kafka-message-management/models"
)

type KafkaBrokerController struct {
	BaseController
}

func (c *KafkaBrokerController) Prepare() {
	c.BaseController.Prepare()
	c.checkAuthor("DataGrid")
}

func (c *KafkaBrokerController) Index() {
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl("kafka_broker/index.html")

	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "kafka_broker/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "kafka_broker/index_footerjs.html"

	c.Data["canEdit"] = c.checkActionAuthor("KafkaBrokerController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("KafkaBrokerController", "Delete")
}

func (c *KafkaBrokerController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.KafkaBrokerQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//获取数据列表和总数
	data, total := models.KafkaBrokerPageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *KafkaBrokerController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := models.KafkaBroker{Id: Id}
	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	c.Data["m"] = m
	c.setTpl("kafka_broker/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "kafka_broker/edit_footerjs.html"
}

func (c *KafkaBrokerController) Save() {
	var err error
	m := models.KafkaBroker{}
	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}
	o := orm.NewOrm()
	if m.Id == 0 {
		if _, err = o.Insert(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
		}

	} else {
		if _, err = o.Update(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}
	c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
}

func (c *KafkaBrokerController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	if num, err := models.KafkaBrokerBatchDelete(ids); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}
