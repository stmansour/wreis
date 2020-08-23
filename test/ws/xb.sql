-- MySQL dump 10.13  Distrib 5.7.22, for osx10.12 (x86_64)
--
-- Host: localhost    Database: wreis
-- ------------------------------------------------------
-- Server version	5.7.22

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
-- Table structure for table `Property`
--

DROP TABLE IF EXISTS `Property`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Property` (
  `PRID` bigint(20) NOT NULL AUTO_INCREMENT,
  `Name` varchar(256) NOT NULL DEFAULT '',
  `YearsInBusiness` smallint(6) NOT NULL DEFAULT '0',
  `ParentCompany` varchar(256) NOT NULL DEFAULT '',
  `URL` varchar(1028) NOT NULL DEFAULT '',
  `Symbol` varchar(128) NOT NULL DEFAULT '',
  `Price` decimal(19,4) NOT NULL DEFAULT '0.0000',
  `DownPayment` decimal(19,4) NOT NULL DEFAULT '0.0000',
  `RentableArea` bigint(20) NOT NULL DEFAULT '0',
  `RentableAreaUnits` smallint(6) NOT NULL DEFAULT '0',
  `LotSize` bigint(20) NOT NULL DEFAULT '0',
  `LotSizeUnits` smallint(6) NOT NULL DEFAULT '0',
  `CapRate` float NOT NULL DEFAULT '0',
  `AvgCap` float NOT NULL DEFAULT '0',
  `BuildDate` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `FLAGS` bigint(20) NOT NULL DEFAULT '0',
  `Ownership` smallint(6) NOT NULL DEFAULT '0',
  `TenantTradeName` varchar(256) NOT NULL DEFAULT '',
  `LeaseGuarantor` smallint(6) NOT NULL DEFAULT '0',
  `LeaseType` smallint(6) NOT NULL DEFAULT '0',
  `DeliveryDt` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `OriginalLeaseTerm` bigint(20) NOT NULL DEFAULT '0',
  `LeaseCommencementDt` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `LeaseExpirationDt` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `TermRemainingOnLease` bigint(20) NOT NULL DEFAULT '0',
  `ROLID` bigint(20) NOT NULL DEFAULT '0',
  `RSLID` bigint(20) NOT NULL DEFAULT '0',
  `Address` varchar(100) NOT NULL DEFAULT '',
  `Address2` varchar(100) NOT NULL DEFAULT '',
  `City` varchar(100) NOT NULL DEFAULT '',
  `State` char(25) NOT NULL DEFAULT '',
  `PostalCode` varchar(100) NOT NULL DEFAULT '',
  `Country` varchar(100) NOT NULL DEFAULT '',
  `LLResponsibilities` varchar(2048) NOT NULL DEFAULT '',
  `NOI` decimal(19,4) NOT NULL DEFAULT '0.0000',
  `HQAddress` varchar(100) NOT NULL DEFAULT '',
  `HQAddress2` varchar(100) NOT NULL DEFAULT '',
  `HQCity` varchar(100) NOT NULL DEFAULT '',
  `HQState` char(25) NOT NULL DEFAULT '',
  `HQPostalCode` varchar(100) NOT NULL DEFAULT '',
  `HQCountry` varchar(100) NOT NULL DEFAULT '',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0',
  `CreateTS` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreateBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`PRID`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Property`
--

LOCK TABLES `Property` WRITE;
/*!40000 ALTER TABLE `Property` DISABLE KEYS */;
INSERT INTO `Property` VALUES (1,'Bill\'s Boar Emporium',8,'','http://bbb.com/','BBE',12345.6700,40000.0000,30000,1,40000,1,0.7,0.6,'2020-03-23 00:00:00',0,0,'Bill\'s Boar Emporium',0,0,'2020-03-23 00:00:00',630720000000000000,'2020-03-23 00:00:00','2020-03-23 00:00:00',630720000000000000,4,0,'1234 Elm Street','','Corn Bluff','AK','98765','USA','roof leaks',30000.0000,'1234 Elm Street','','Corn Bluff','AK','98765','USA','2020-07-25 05:52:47',197,'2020-07-16 08:34:53',190),(3,'Bill\'s Bungalo Emporium',5,'','https://bbe.com/','',12345.8900,0.0000,50000,0,60000,0,0.3,0.27,'1975-01-01 00:00:00',5,0,'Bill\'s Bungalo Emporium',0,1,'1975-01-01 00:00:00',30,'2018-06-15 00:00:00','2020-06-15 00:00:00',70,1,0,'1234 Elm Street','','Goober','AK','12345','USA','',25000.0000,'1234 Elm Street','','Goober','AK','12345','USA','2020-07-25 05:52:47',198,'2020-07-16 08:36:02',191),(4,'Sally\'s Sludge Salon',5,'','https://bbe.com/','',12345.8900,0.0000,50000,0,60000,0,0.3,0.27,'1975-01-02 00:00:00',2,0,'Sally\'s Sludge Salon',0,1,'1975-01-02 00:00:00',31,'2018-06-16 00:00:00','2020-06-16 00:00:00',71,0,3,'1235 Elm Street','','Goober','AK','12345','USA','',25000.0000,'1235 Elm Street','','Goober','AK','12345','USA','2020-07-25 05:52:47',199,'2020-07-16 08:36:02',192),(5,'Mungo\'s Mud',5,'','https://bbe.com/','',12345.8900,0.0000,50000,0,60000,0,0.3,0.27,'1975-01-03 00:00:00',5,0,'Mungo\'s Mud',1,1,'1975-01-03 00:00:00',32,'2018-06-17 00:00:00','2020-06-17 00:00:00',72,2,0,'1236 Elm Street','','Goober','AK','12345','USA','',25000.0000,'1236 Elm Street','','Goober','AK','12345','USA','2020-07-25 05:52:47',200,'2020-07-16 08:36:02',193),(6,'Jimbo\'s Junk Yard',5,'','https://bbe.com/','',12345.8900,0.0000,50000,0,60000,0,0.3,0.27,'1975-01-04 00:00:00',2,0,'Jimbo\'s Junk Yard',1,1,'1975-01-04 00:00:00',33,'2018-06-18 00:00:00','2020-06-18 00:00:00',73,0,4,'1237 Elm Street','','Goober','AK','12345','USA','',25000.0000,'1237 Elm Street','','Goober','AK','12345','USA','2020-07-25 05:52:47',201,'2020-07-16 08:36:02',194),(7,'Rosita\'s Taco Town',5,'','https://bbe.com/','',12345.8900,0.0000,50000,0,60000,0,0.3,0.27,'1975-01-05 00:00:00',5,0,'Rosita\'s Taco Town',0,1,'1975-01-05 00:00:00',34,'2018-06-19 00:00:00','2020-06-19 00:00:00',74,0,0,'1238 Elm Street','','Goober','AK','12345','USA','',25000.0000,'1238 Elm Street','','Goober','AK','12345','USA','2020-07-25 05:52:47',202,'2020-07-16 08:36:02',195),(8,'Wings \'n Such',5,'','https://bbe.com/','',12345.8900,0.0000,50000,0,60000,0,0.3,0.27,'1975-01-06 00:00:00',2,0,'Wings \'n Such',0,1,'1975-01-06 00:00:00',35,'2018-06-20 00:00:00','2020-06-20 00:00:00',75,0,0,'1239 Elm Street','','Goober','AK','12345','USA','',25000.0000,'1239 Elm Street','','Goober','AK','12345','USA','2020-07-25 05:52:47',203,'2020-07-16 08:36:02',196);
/*!40000 ALTER TABLE `Property` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `RenewOption`
--

DROP TABLE IF EXISTS `RenewOption`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `RenewOption` (
  `ROID` bigint(20) NOT NULL AUTO_INCREMENT,
  `ROLID` bigint(20) NOT NULL DEFAULT '0',
  `Dt` date NOT NULL DEFAULT '1970-01-01',
  `Count` bigint(20) NOT NULL DEFAULT '0',
  `Opt` bigint(20) NOT NULL DEFAULT '0',
  `Rent` decimal(19,4) NOT NULL DEFAULT '0.0000',
  `FLAGS` bigint(20) NOT NULL DEFAULT '0',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0',
  `CreateTS` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreateBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`ROID`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `RenewOption`
--

LOCK TABLES `RenewOption` WRITE;
/*!40000 ALTER TABLE `RenewOption` DISABLE KEYS */;
INSERT INTO `RenewOption` VALUES (1,1,'2024-07-04',0,1,109709.4500,1,'2020-07-16 08:36:02',0,'2020-07-16 08:36:02',0),(2,1,'2025-07-04',0,2,111903.6300,1,'2020-07-16 08:36:02',0,'2020-07-16 08:36:02',0),(3,1,'2026-07-04',0,3,114141.7100,1,'2020-07-16 08:36:02',0,'2020-07-16 08:36:02',0),(4,2,'0000-00-00',1,1,109709.4500,0,'2020-07-16 08:36:02',0,'2020-07-16 08:36:02',0),(5,2,'0000-00-00',2,2,111903.6300,0,'2020-07-16 08:36:02',0,'2020-07-16 08:36:02',0),(6,2,'0000-00-00',3,3,114141.7100,0,'2020-07-16 08:36:02',0,'2020-07-16 08:36:02',0);
/*!40000 ALTER TABLE `RenewOption` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `RenewOptions`
--

DROP TABLE IF EXISTS `RenewOptions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `RenewOptions` (
  `ROLID` bigint(20) NOT NULL AUTO_INCREMENT,
  `FLAGS` bigint(20) NOT NULL DEFAULT '0',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0',
  `CreateTS` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreateBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`ROLID`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `RenewOptions`
--

LOCK TABLES `RenewOptions` WRITE;
/*!40000 ALTER TABLE `RenewOptions` DISABLE KEYS */;
INSERT INTO `RenewOptions` VALUES (1,1,'2020-07-16 08:36:02',0,'2020-07-16 08:36:02',0),(2,0,'2020-07-16 08:36:02',0,'2020-07-16 08:36:02',0);
/*!40000 ALTER TABLE `RenewOptions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `RentStep`
--

DROP TABLE IF EXISTS `RentStep`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `RentStep` (
  `RSID` bigint(20) NOT NULL AUTO_INCREMENT,
  `RSLID` bigint(20) NOT NULL DEFAULT '0',
  `Dt` date NOT NULL DEFAULT '1970-01-01',
  `Opt` varchar(100) NOT NULL DEFAULT '',
  `Rent` decimal(19,4) NOT NULL DEFAULT '0.0000',
  `FLAGS` bigint(20) NOT NULL DEFAULT '0',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0',
  `CreateTS` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreateBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`RSID`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `RentStep`
--

LOCK TABLES `RentStep` WRITE;
/*!40000 ALTER TABLE `RentStep` DISABLE KEYS */;
INSERT INTO `RentStep` VALUES (1,1,'2020-03-23','0',2760.0000,0,'2020-07-16 08:34:53',0,'2020-07-16 08:34:53',0),(3,3,'2024-07-04','1',2850.0000,1,'2020-08-07 23:05:35',0,'2020-07-16 08:36:02',0),(4,3,'2025-07-04','2',2900.0000,1,'2020-08-07 23:05:35',0,'2020-07-16 08:36:02',0),(5,3,'2026-07-04','3',3000.0000,1,'2020-08-07 23:05:35',0,'2020-07-16 08:36:02',0),(6,4,'2018-01-01','Year 1',3100.0000,0,'2020-08-22 06:54:58',0,'2020-07-16 08:36:02',0),(7,4,'2019-01-01','Year 2',3200.0000,0,'2020-08-22 06:54:58',0,'2020-07-16 08:36:02',0),(8,4,'2020-01-01','Year 3',3300.0000,0,'2020-08-22 06:54:58',0,'2020-07-16 08:36:02',0);
/*!40000 ALTER TABLE `RentStep` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `RentSteps`
--

DROP TABLE IF EXISTS `RentSteps`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `RentSteps` (
  `RSLID` bigint(20) NOT NULL AUTO_INCREMENT,
  `FLAGS` bigint(20) NOT NULL DEFAULT '0',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0',
  `CreateTS` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreateBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`RSLID`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `RentSteps`
--

LOCK TABLES `RentSteps` WRITE;
/*!40000 ALTER TABLE `RentSteps` DISABLE KEYS */;
INSERT INTO `RentSteps` VALUES (1,0,'2020-07-16 08:34:53',0,'2020-07-16 08:34:53',0),(3,1,'2020-07-16 08:36:02',0,'2020-07-16 08:36:02',0),(4,0,'2020-07-16 08:36:02',0,'2020-07-16 08:36:02',0);
/*!40000 ALTER TABLE `RentSteps` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-08-21 23:55:04