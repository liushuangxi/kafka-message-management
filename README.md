# kafka-message-management
## source code install
### step 1
go get -v https://github.com/liushuangxi/kafka-message-management

or

git clone https://github.com/liushuangxi/kafka-message-management
### step 2
vim conf/app.conf (set mysql)

mysql -uroot -p123456 kmm < private/kmm.sql

### step 3
go run main.go

## direct execute
### step 1
go get -v https://github.com/liushuangxi/kafka-message-management

or

git clone https://github.com/liushuangxi/kafka-message-management
### step 2

mysql -uroot -p123456 kmm < private/kmm.sql

### step 3
https://github.com/liushuangxi/resource/blob/master/kafka-message-management

./kafka-message-management -mysql_host='127.0.0.1' -mysql_port=3306 -mysql_db_name='kmm' -mysql_user='root' -mysql_pass='123456'

## View
![avatar](https://static.studygolang.com/201013/c885de3f5a8e70e675883020f2a71cb5.png)
<br>
![avatar](https://static.studygolang.com/201013/8ce5ac721733721c91f9992561da6f1a.png)
