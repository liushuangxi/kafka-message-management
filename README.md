# kafka-message-management
## source install
### step 1
go get -v https://github.com/liushuangxi/kafka-message-management

or

git clone https://github.com/liushuangxi/kafka-message-management
### step 2
vim conf/app.conf (set mysql)

mysql -uroot -p123456 kmm < private/kmm.sql

### step 3
go run main.go