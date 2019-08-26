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
./private/kafka-message-management -mysql_host='127.0.0.1' -mysql_port=3306 -mysql_db_name='kmm' -mysql_user='root' -mysql_pass='123456'