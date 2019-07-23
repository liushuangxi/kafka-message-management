package models

import (
	"github.com/astaxie/beego/orm"
)

// TableName 设置KafkaTopicCollect表名
func (a *KafkaTopicCollect) TableName() string {
	return KafkaTopicCollectTBName()
}

// KafkaTopicCollectQueryParam 用于查询的类
type KafkaTopicCollectQueryParam struct {
	BaseQueryParam
	UserId int64
	Broker string
	Topic  string
}

// KafkaTopicCollect 实体类
type KafkaTopicCollect struct {
	Id     int
	UserId int    `orm:"size(11)"`
	Broker string `orm:"size(300)"`
	Topic  string `orm:"size(300)"`
}

// KafkaTopicCollectPageList 获取分页数据
func KafkaTopicCollectPageList(params *KafkaTopicCollectQueryParam) ([]*KafkaTopicCollect, int64) {
	query := orm.NewOrm().QueryTable(KafkaTopicCollectTBName())
	data := make([]*KafkaTopicCollect, 0)
	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Broker":
		sortorder = "Broker"
	case "Topic":
		sortorder = "Topic"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}

	if params.UserId > 0 {
		query = query.Filter("user_id", params.UserId)
	}

	query = query.Filter("broker__icontains", params.Broker)
	query = query.Filter("topic__icontains", params.Topic)

	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

func KafkaTopicCollectTotal() int64 {
	query := orm.NewOrm().QueryTable(KafkaTopicCollectTBName())

	total, _ := query.Count()

	return total
}

// KafkaTopicCollectOne 根据id获取单条
func KafkaTopicCollectOne(id int) (*KafkaTopicCollect, error) {
	o := orm.NewOrm()
	m := KafkaTopicCollect{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
