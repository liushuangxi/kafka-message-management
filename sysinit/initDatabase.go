package sysinit

import (
	"flag"

	_ "github.com/liushuangxi/kafka-message-management/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
)

//初始化数据连接
func InitDatabase() {
	//读取配置文件，设置数据库参数
	//数据库类别
	dbType := beego.AppConfig.String("db_type")
	//连接名称
	dbAlias := beego.AppConfig.String(dbType + "::db_alias")
	//数据库名称
	dbName := beego.AppConfig.String(dbType + "::db_name")
	//数据库连接用户名
	dbUser := beego.AppConfig.String(dbType + "::db_user")
	//数据库连接用户名
	dbPwd := beego.AppConfig.String(dbType + "::db_pwd")
	//数据库IP（域名）
	dbHost := beego.AppConfig.String(dbType + "::db_host")
	//数据库端口
	dbPort := beego.AppConfig.String(dbType + "::db_port")

	//解析命令行参数
	cliHost := flag.String("mysql_host", "", "mysql host")
	cliPort := flag.String("mysql_port", "", "mysql port")
	cliName := flag.String("mysql_db_name", "", "mysql db name")
	cliUser := flag.String("mysql_user", "", "mysql username")
	cliPass := flag.String("mysql_pass", "", "mysql password")
	flag.Parse()

	if *cliHost != "" {
		dbHost = *cliHost
	}
	if *cliPort != "" {
		dbPort = *cliPort
	}
	if *cliName != "" {
		dbName = *cliName
	}
	if *cliUser != "" {
		dbUser = *cliUser
	}
	if *cliPass != "" {
		dbPwd = *cliPass
	}

	switch dbType {
	case "sqlite3":
		orm.RegisterDataBase(dbAlias, dbType, dbName)
	case "mysql":
		dbCharset := beego.AppConfig.String(dbType + "::db_charset")
		orm.RegisterDataBase(dbAlias, dbType, dbUser+":"+dbPwd+"@tcp("+dbHost+":"+
			dbPort+")/"+dbName+"?charset="+dbCharset, 30)
	}
	//如果是开发模式，则显示命令信息
	isDev := (beego.AppConfig.String("runmode") == "dev")
	//自动建表
	orm.RunSyncdb("default", false, isDev)
	if isDev {
		orm.Debug = isDev
	}
}
