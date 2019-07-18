package models

import (
	"github.com/astaxie/beego/orm"

	"github.com/liushuangxi/kafka-message-management/utils"
)

// TableName 设置KafkaBroker表名
func (a *KafkaBroker) TableName() string {
	return KafkaBrokerTBName()
}

// KafkaBrokerQueryParam 用于查询的类
type KafkaBrokerQueryParam struct {
	BaseQueryParam
	Broker string
	Alias  string
}

// KafkaBroker 实体类
type KafkaBroker struct {
	Id     int
	Broker string `orm:"size(300)"`
	Alias  string `orm:"size(300)"`
}

// KafkaBrokerPageList 获取分页数据
func KafkaBrokerPageList(params *KafkaBrokerQueryParam) ([]*KafkaBroker, int64) {
	query := orm.NewOrm().QueryTable(KafkaBrokerTBName())
	data := make([]*KafkaBroker, 0)
	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}

	query = query.Filter("broker__icontains", params.Broker)
	query = query.Filter("alias__icontains", params.Alias)

	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

func KafkaBrokerTotal() (int64) {
	query := orm.NewOrm().QueryTable(KafkaBrokerTBName())

	total, _ := query.Count()

	return total
}

func KafkaTopicTotal() (int64) {
	query := orm.NewOrm().QueryTable(KafkaBrokerTBName())
	data := make([]*KafkaBroker, 0)
	query.Limit(100, 0).All(&data)

	total := 0;
	for _, broker := range data {
		topics, _ := utils.GetTopics(broker.Broker)

		total += len(topics)
	}

	return int64(total)
}

// KafkaBrokerOne 根据id获取单条
func KafkaBrokerOne(id int) (*KafkaBroker, error) {
	o := orm.NewOrm()
	m := KafkaBroker{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func KafkaBrokerBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(KafkaBrokerTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}
