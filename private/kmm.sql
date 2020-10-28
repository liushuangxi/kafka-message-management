-- MySQL dump 10.13  Distrib 5.7.32, for Linux (x86_64)
--
-- Host: localhost    Database: kmm
-- ------------------------------------------------------
-- Server version	5.7.32

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `rms_backend_user`
--

DROP TABLE IF EXISTS `rms_backend_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rms_backend_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `real_name` varchar(255) NOT NULL DEFAULT '',
  `user_name` varchar(255) NOT NULL DEFAULT '',
  `user_pwd` varchar(255) NOT NULL DEFAULT '',
  `is_super` tinyint(1) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0',
  `mobile` varchar(16) NOT NULL DEFAULT '',
  `email` varchar(256) NOT NULL DEFAULT '',
  `avatar` varchar(256) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rms_backend_user`
--

LOCK TABLES `rms_backend_user` WRITE;
/*!40000 ALTER TABLE `rms_backend_user` DISABLE KEYS */;
INSERT INTO `rms_backend_user` VALUES (1,'Lucky Wolf','admin','e10adc3949ba59abbe56e057f20f883e',1,1,'18701502926','729021170@163.com','/static/upload/kafka.jpg'),(6,'kafka','kafka','541661622185851c248b41bf0cea7ad0',0,1,'','','/static/upload/kafka.jpg');
/*!40000 ALTER TABLE `rms_backend_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rms_backend_user_rms_roles`
--

DROP TABLE IF EXISTS `rms_backend_user_rms_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rms_backend_user_rms_roles` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `rms_backend_user_id` int(11) NOT NULL,
  `rms_role_id` int(11) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rms_backend_user_rms_roles`
--

LOCK TABLES `rms_backend_user_rms_roles` WRITE;
/*!40000 ALTER TABLE `rms_backend_user_rms_roles` DISABLE KEYS */;
/*!40000 ALTER TABLE `rms_backend_user_rms_roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rms_kafka_broker`
--

DROP TABLE IF EXISTS `rms_kafka_broker`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rms_kafka_broker` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `broker` varchar(300) NOT NULL DEFAULT '',
  `alias` varchar(300) NOT NULL DEFAULT '',
  `cluster` varchar(300) NOT NULL DEFAULT '',
  `manager` varchar(300) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_broker` (`broker`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rms_kafka_broker`
--

LOCK TABLES `rms_kafka_broker` WRITE;
/*!40000 ALTER TABLE `rms_kafka_broker` DISABLE KEYS */;
/*!40000 ALTER TABLE `rms_kafka_broker` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rms_kafka_topic_collect`
--

DROP TABLE IF EXISTS `rms_kafka_topic_collect`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rms_kafka_topic_collect` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0',
  `broker` varchar(300) NOT NULL DEFAULT '',
  `topic` varchar(300) NOT NULL DEFAULT '',
  `alias` varchar(300) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rms_kafka_topic_collect`
--

LOCK TABLES `rms_kafka_topic_collect` WRITE;
/*!40000 ALTER TABLE `rms_kafka_topic_collect` DISABLE KEYS */;
/*!40000 ALTER TABLE `rms_kafka_topic_collect` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rms_resource`
--

DROP TABLE IF EXISTS `rms_resource`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rms_resource` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `rtype` int(11) NOT NULL DEFAULT '0',
  `name` varchar(64) NOT NULL DEFAULT '',
  `parent_id` int(11) DEFAULT NULL,
  `seq` int(11) NOT NULL DEFAULT '0',
  `icon` varchar(32) NOT NULL DEFAULT '',
  `url_for` varchar(256) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=45 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rms_resource`
--

LOCK TABLES `rms_resource` WRITE;
/*!40000 ALTER TABLE `rms_resource` DISABLE KEYS */;
INSERT INTO `rms_resource` VALUES (7,1,'权限管理',8,100,'fa fa-balance-scale',''),(8,0,'系统菜单',NULL,200,'',''),(9,1,'资源管理',7,100,'','ResourceController.Index'),(12,1,'角色管理',7,100,'','RoleController.Index'),(13,1,'用户管理',7,100,'','BackendUserController.Index'),(21,0,'业务菜单',NULL,170,'',''),(25,2,'编辑',9,100,'fa fa-pencil','ResourceController.Edit'),(26,2,'编辑',13,100,'fa fa-pencil','BackendUserController.Edit'),(27,2,'删除',9,100,'fa fa-trash','ResourceController.Delete'),(29,2,'删除',13,100,'fa fa-trash','BackendUserController.Delete'),(30,2,'编辑',12,100,'fa fa-pencil','RoleController.Edit'),(31,2,'删除',12,100,'fa fa-trash','RoleController.Delete'),(32,2,'分配资源',12,100,'fa fa-th','RoleController.Allocate'),(35,1,' 首页',NULL,100,'fa fa-dashboard','HomeController.Index'),(39,1,'Kafka Broker',21,100,'','KafkaBrokerController.Index'),(40,1,'Topic Collect',21,100,'','KafkaTopicCollectController.Index'),(41,2,'Kafka Message Permission',39,100,'','KafkaMessageController.Index'),(42,2,'Kafka Topic Permission',39,100,'','KafkaTopicController.Index'),(43,2,'Topic Collect Update',40,100,'','KafkaTopicCollectController.Update'),(44,2,'Topic Collect Delete',40,100,'','KafkaTopicCollectController.Delete');
/*!40000 ALTER TABLE `rms_resource` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rms_role`
--

DROP TABLE IF EXISTS `rms_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rms_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '',
  `seq` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rms_role`
--

LOCK TABLES `rms_role` WRITE;
/*!40000 ALTER TABLE `rms_role` DISABLE KEYS */;
INSERT INTO `rms_role` VALUES (22,'超级管理员',20),(24,'角色管理员',10),(26,'Kafka 消息查看',100);
/*!40000 ALTER TABLE `rms_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rms_role_backenduser_rel`
--

DROP TABLE IF EXISTS `rms_role_backenduser_rel`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rms_role_backenduser_rel` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL,
  `backend_user_id` int(11) NOT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=70 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rms_role_backenduser_rel`
--

LOCK TABLES `rms_role_backenduser_rel` WRITE;
/*!40000 ALTER TABLE `rms_role_backenduser_rel` DISABLE KEYS */;
INSERT INTO `rms_role_backenduser_rel` VALUES (68,22,1,'2019-07-04 23:14:22'),(69,26,6,'2019-07-04 23:29:01');
/*!40000 ALTER TABLE `rms_role_backenduser_rel` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rms_role_resource_rel`
--

DROP TABLE IF EXISTS `rms_role_resource_rel`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rms_role_resource_rel` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL,
  `resource_id` int(11) NOT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=513 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rms_role_resource_rel`
--

LOCK TABLES `rms_role_resource_rel` WRITE;
/*!40000 ALTER TABLE `rms_role_resource_rel` DISABLE KEYS */;
INSERT INTO `rms_role_resource_rel` VALUES (448,24,8,'2017-12-19 06:40:16'),(474,22,35,'2019-06-27 23:09:46'),(475,22,21,'2019-06-27 23:09:46'),(480,22,8,'2019-06-27 23:09:46'),(481,22,7,'2019-06-27 23:09:46'),(482,22,9,'2019-06-27 23:09:46'),(483,22,25,'2019-06-27 23:09:46'),(484,22,27,'2019-06-27 23:09:46'),(485,22,12,'2019-06-27 23:09:46'),(486,22,30,'2019-06-27 23:09:46'),(487,22,31,'2019-06-27 23:09:46'),(488,22,32,'2019-06-27 23:09:46'),(489,22,13,'2019-06-27 23:09:46'),(490,22,26,'2019-06-27 23:09:46'),(491,22,29,'2019-06-27 23:09:46'),(506,26,21,'2019-07-24 12:16:03'),(507,26,39,'2019-07-24 12:16:03'),(508,26,41,'2019-07-24 12:16:03'),(509,26,42,'2019-07-24 12:16:03'),(510,26,40,'2019-07-24 12:16:03'),(511,26,43,'2019-07-24 12:16:03'),(512,26,44,'2019-07-24 12:16:03');
/*!40000 ALTER TABLE `rms_role_resource_rel` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-10-28 18:08:07
