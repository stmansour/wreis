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
-- Table structure for table `DataUpdate`
--

DROP TABLE IF EXISTS `DataUpdate`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `DataUpdate` (
  `DUID` bigint(20) NOT NULL AUTO_INCREMENT,
  `GID` bigint(20) NOT NULL DEFAULT '0',
  `DtStart` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `DtStop` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`DUID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `DataUpdate`
--

LOCK TABLES `DataUpdate` WRITE;
/*!40000 ALTER TABLE `DataUpdate` DISABLE KEYS */;
/*!40000 ALTER TABLE `DataUpdate` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `EGroup`
--

DROP TABLE IF EXISTS `EGroup`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `EGroup` (
  `GID` bigint(20) NOT NULL AUTO_INCREMENT,
  `GroupName` varchar(50) NOT NULL DEFAULT '',
  `GroupDescription` varchar(1000) NOT NULL DEFAULT '',
  `DtStart` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `DtStop` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`GID`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `EGroup`
--

LOCK TABLES `EGroup` WRITE;
/*!40000 ALTER TABLE `EGroup` DISABLE KEYS */;
INSERT INTO `EGroup` VALUES (1,'smanmusic','','2020-08-26 06:58:12','2020-08-26 06:58:12','2020-08-26 06:58:11',0);
/*!40000 ALTER TABLE `EGroup` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `PGroup`
--

DROP TABLE IF EXISTS `PGroup`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `PGroup` (
  `PID` bigint(20) NOT NULL DEFAULT '0',
  `GID` bigint(20) NOT NULL DEFAULT '0',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `PGroup`
--

LOCK TABLES `PGroup` WRITE;
/*!40000 ALTER TABLE `PGroup` DISABLE KEYS */;
INSERT INTO `PGroup` VALUES (1,1,'2020-07-05 18:29:18',0),(2,1,'2020-07-05 18:29:18',0);
/*!40000 ALTER TABLE `PGroup` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `People`
--

DROP TABLE IF EXISTS `People`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `People` (
  `PID` bigint(20) NOT NULL AUTO_INCREMENT,
  `FirstName` varchar(100) DEFAULT '',
  `MiddleName` varchar(100) DEFAULT '',
  `LastName` varchar(100) DEFAULT '',
  `PreferredName` varchar(100) DEFAULT '',
  `JobTitle` varchar(100) DEFAULT '',
  `OfficePhone` varchar(100) DEFAULT '',
  `OfficeFax` varchar(100) DEFAULT '',
  `Email1` varchar(50) DEFAULT '',
  `Email2` varchar(50) NOT NULL DEFAULT '',
  `MailAddress` varchar(50) DEFAULT '',
  `MailAddress2` varchar(50) DEFAULT '',
  `MailCity` varchar(100) DEFAULT '',
  `MailState` varchar(50) DEFAULT '',
  `MailPostalCode` varchar(50) DEFAULT '',
  `MailCountry` varchar(50) DEFAULT '',
  `RoomNumber` varchar(50) DEFAULT '',
  `MailStop` varchar(100) DEFAULT '',
  `Status` smallint(6) DEFAULT '0',
  `OptOutDate` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`PID`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `People`
--

LOCK TABLES `People` WRITE;
/*!40000 ALTER TABLE `People` DISABLE KEYS */;
INSERT INTO `People` VALUES (1,'Shannon','CornDog','Kodiak','','','','','shannonkodiak1964@gmail.com','','','','','','','','','',0,'0000-00-00 00:00:00','2020-08-26 06:58:11',0),(2,'Debbie','','Van Compernolle  Conway','','','','','2cdeb650@gmail.com','','30101 East Hanna','','Buckner','MO','64016','','','',0,'0000-00-00 00:00:00','2020-07-05 18:29:18',0);
/*!40000 ALTER TABLE `People` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Property`
--

DROP TABLE IF EXISTS `Property`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Property` (
  `PRID` bigint(20) NOT NULL AUTO_INCREMENT,
  `Name` varchar(256) NOT NULL DEFAULT '',
  `YearFounded` smallint(6) NOT NULL DEFAULT '0',
  `ParentCompany` varchar(256) NOT NULL DEFAULT '',
  `URL` varchar(1028) NOT NULL DEFAULT '',
  `Symbol` varchar(128) NOT NULL DEFAULT '',
  `Price` decimal(19,4) NOT NULL DEFAULT '0.0000',
  `DownPayment` decimal(19,4) NOT NULL DEFAULT '0.0000',
  `RentableArea` bigint(20) NOT NULL DEFAULT '0',
  `RentableAreaUnits` smallint(6) NOT NULL DEFAULT '0',
  `LotSize` decimal(19,4) NOT NULL DEFAULT '0.0000',
  `LotSizeUnits` smallint(6) NOT NULL DEFAULT '0',
  `CapRate` float NOT NULL DEFAULT '0',
  `AvgCap` float NOT NULL DEFAULT '0',
  `BuildYear` smallint(6) NOT NULL DEFAULT '0',
  `FLAGS` bigint(20) NOT NULL DEFAULT '0',
  `Ownership` smallint(6) NOT NULL DEFAULT '0',
  `TenantTradeName` varchar(256) NOT NULL DEFAULT '',
  `LeaseGuarantor` smallint(6) NOT NULL DEFAULT '0',
  `LeaseType` smallint(6) NOT NULL DEFAULT '0',
  `OriginalLeaseTerm` bigint(20) NOT NULL DEFAULT '0',
  `RentCommencementDt` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `LeaseExpirationDt` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `TermRemainingOnLease` bigint(20) NOT NULL DEFAULT '0',
  `TermRemainingOnLeaseUnits` smallint(6) NOT NULL DEFAULT '0',
  `ROLID` bigint(20) NOT NULL DEFAULT '0',
  `RSLID` bigint(20) NOT NULL DEFAULT '0',
  `FlowState` bigint(20) NOT NULL DEFAULT '0',
  `Address` varchar(100) NOT NULL DEFAULT '',
  `Address2` varchar(100) NOT NULL DEFAULT '',
  `City` varchar(100) NOT NULL DEFAULT '',
  `State` char(25) NOT NULL DEFAULT '',
  `PostalCode` varchar(100) NOT NULL DEFAULT '',
  `Country` varchar(100) NOT NULL DEFAULT '',
  `LLResponsibilities` varchar(2048) NOT NULL DEFAULT '',
  `NOI` decimal(19,4) NOT NULL DEFAULT '0.0000',
  `HQCity` varchar(100) NOT NULL DEFAULT '',
  `HQState` char(25) NOT NULL DEFAULT '',
  `Img1` varchar(2048) NOT NULL DEFAULT '',
  `Img2` varchar(2048) NOT NULL DEFAULT '',
  `Img3` varchar(2048) NOT NULL DEFAULT '',
  `Img4` varchar(2048) NOT NULL DEFAULT '',
  `Img5` varchar(2048) NOT NULL DEFAULT '',
  `Img6` varchar(2048) NOT NULL DEFAULT '',
  `Img7` varchar(2048) NOT NULL DEFAULT '',
  `Img8` varchar(2048) NOT NULL DEFAULT '',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0',
  `CreateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreateBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`PRID`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Property`
--

LOCK TABLES `Property` WRITE;
/*!40000 ALTER TABLE `Property` DISABLE KEYS */;
INSERT INTO `Property` VALUES (1,'Bill\'s Boar Emporium',8,'','http://bbb.com/','BBE',12345.6700,40000.0000,8790,0,40.0000,1,0.7,0.6,32767,0,0,'Bill\'s Boar Emporium',0,0,20,'2020-03-23 00:00:00','2020-03-23 00:00:00',25,0,0,0,1,'1234 Elm Street','','Corn Bluff','AK','98765','USA','roof leaks',30000.0000,'Corn Bluff','AK','','','','','','','','','2020-09-28 20:02:15',197,'2020-07-16 08:34:53',190),(3,'Bill\'s Bungalo Emporium',5,'','https://bbe.com/','',12345.8900,0.0000,50000,0,60000.0000,0,0.3,0.27,32767,5,1,'Bill\'s Bungalo Emporium',0,1,30,'2018-06-15 00:00:00','2020-06-15 00:00:00',7,1,1,0,2,'1234 Elm Street','','Kalamazoo','MI','12345','USA','',25000.0000,'Goober','AK','','','','','','','','','2020-09-28 20:08:29',198,'2020-07-16 08:36:02',191),(4,'Sally\'s Sludge Salon',5,'','https://bbe.com/','',12345.8900,510000.0000,16000,0,26.0000,1,0.3,0.27,32767,2,0,'Sally\'s Sludge Salon',0,2,31,'2018-06-16 00:00:00','2020-06-16 00:00:00',71,0,2,3,3,'1235 Elm Street','','Suck-egg Hollow','TN','12345','USA','',25000.0000,'Goober','AK','','','','','','','','','2020-09-28 20:08:29',199,'2020-07-16 08:36:02',192),(5,'Mungo\'s Mud',5,'','https://bbe.com/','',12345.8900,0.0000,24171,0,60000.0000,0,0.3,0.27,32767,5,1,'Mungo\'s Mud',1,3,32,'2018-06-17 00:00:00','2020-06-17 00:00:00',5,1,0,0,4,'1236 Elm Street','','Rabbit Hash','KY','12345','USA','',25000.0000,'Goober','AK','','','','','','','','','2020-09-28 20:08:29',200,'2020-07-16 08:36:02',193),(6,'Jimbo\'s Junk Yard',5,'','https://bbe.com/','',12345.8900,29000.0000,18500,0,17.0000,1,0.3,0.27,32767,2,0,'Jimbo\'s Junk Yard',1,1,33,'2018-06-18 00:00:00','2020-06-18 00:00:00',73,0,0,4,5,'1237 Elm Street','','Gumlog','GA','12345','USA','',25000.0000,'Goober','AK','','','','','','','','','2020-09-28 20:08:29',201,'2020-07-16 08:36:02',194),(7,'Rosita\'s Taco Town',5,'','https://bbe.com/','',12345.8900,0.0000,50000,0,60000.0000,0,0.3,0.27,32767,5,1,'Rosita\'s Taco Town',0,2,34,'2018-06-19 00:00:00','2020-06-19 00:00:00',15,1,0,0,6,'1238 Elm Street','','Frog Eye','AL','12345','USA','',25000.0000,'Goober','AK','','','','','','','','','2020-09-28 20:08:29',202,'2020-07-16 08:36:02',195),(8,'Wings \'n Such',5,'','https://bbe.com/','',12345.8900,400000.0000,6500,0,9.0000,1,0.3,0.27,32767,2,0,'Wings \'n Such',0,3,35,'2018-06-20 00:00:00','2020-06-20 00:00:00',75,0,0,0,7,'1239 Elm Street','','Nimrod','MN','12345','USA','',25000.0000,'Goober','AK','','','','','','','','','2020-09-28 20:08:29',203,'2020-07-16 08:36:02',196);
/*!40000 ALTER TABLE `Property` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Query`
--

DROP TABLE IF EXISTS `Query`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Query` (
  `QID` bigint(20) NOT NULL AUTO_INCREMENT,
  `QueryName` varchar(50) DEFAULT '',
  `QueryDescr` varchar(1000) DEFAULT '',
  `QueryJSON` varchar(3000) DEFAULT '',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`QID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Query`
--

LOCK TABLES `Query` WRITE;
/*!40000 ALTER TABLE `Query` DISABLE KEYS */;
/*!40000 ALTER TABLE `Query` ENABLE KEYS */;
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
  `Opt` varchar(100) NOT NULL DEFAULT '',
  `Rent` decimal(19,4) NOT NULL DEFAULT '0.0000',
  `FLAGS` bigint(20) NOT NULL DEFAULT '0',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0',
  `CreateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreateBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`ROID`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `RenewOption`
--

LOCK TABLES `RenewOption` WRITE;
/*!40000 ALTER TABLE `RenewOption` DISABLE KEYS */;
INSERT INTO `RenewOption` VALUES (1,1,'2024-07-04','1',109709.4500,1,'2020-07-16 08:36:02',0,'2020-07-16 08:36:02',0),(2,1,'2025-07-04','2',111903.6300,1,'2020-07-16 08:36:02',0,'2020-07-16 08:36:02',0),(3,1,'2026-07-04','3',114141.7100,1,'2020-07-16 08:36:02',0,'2020-07-16 08:36:02',0),(4,2,'0000-00-00','Year 1',109709.4500,0,'2020-08-28 19:56:22',0,'2020-07-16 08:36:02',0),(5,2,'0000-00-00','Double Year 2',111903.6300,0,'2020-08-28 19:56:22',0,'2020-07-16 08:36:02',0),(6,2,'0000-00-00','Triple Year 3',114141.7100,0,'2020-08-28 19:56:22',0,'2020-07-16 08:36:02',0);
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
  `CreateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
  `CreateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreateBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`RSID`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `RentStep`
--

LOCK TABLES `RentStep` WRITE;
/*!40000 ALTER TABLE `RentStep` DISABLE KEYS */;
INSERT INTO `RentStep` VALUES (1,1,'2020-03-23','0',2760.0000,0,'2020-07-16 08:34:53',0,'2020-07-16 08:34:53',0),(3,3,'2024-07-04','Sally\'s 1',12850.0000,1,'2020-08-25 21:03:59',0,'2020-07-16 08:36:02',0),(4,3,'2025-07-04','Sally\'s 2',12900.0000,1,'2020-08-25 21:03:59',0,'2020-07-16 08:36:02',0),(5,3,'2026-07-04','Sally\'s 3',13000.0000,1,'2020-08-25 21:03:59',0,'2020-07-16 08:36:02',0),(6,4,'2018-01-01','Year 1',3100.0000,0,'2020-08-22 06:54:58',0,'2020-07-16 08:36:02',0),(7,4,'2019-01-01','Year 2',3200.0000,0,'2020-08-22 06:54:58',0,'2020-07-16 08:36:02',0),(8,4,'2020-01-01','Year 3',3300.0000,0,'2020-08-22 06:54:58',0,'2020-07-16 08:36:02',0);
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
  `CreateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
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

--
-- Table structure for table `StateInfo`
--

DROP TABLE IF EXISTS `StateInfo`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `StateInfo` (
  `SIID` bigint(20) NOT NULL AUTO_INCREMENT,
  `PRID` bigint(20) NOT NULL DEFAULT '0',
  `FLAGS` bigint(20) NOT NULL DEFAULT '0',
  `FlowState` bigint(20) NOT NULL DEFAULT '0',
  `OwnerUID` bigint(20) NOT NULL DEFAULT '0',
  `OwnerDt` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `ApproverUID` bigint(20) NOT NULL DEFAULT '0',
  `ApproverDt` date NOT NULL DEFAULT '1970-01-01',
  `Reason` varchar(256) NOT NULL DEFAULT '',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0',
  `CreateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreateBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`SIID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `StateInfo`
--

LOCK TABLES `StateInfo` WRITE;
/*!40000 ALTER TABLE `StateInfo` DISABLE KEYS */;
/*!40000 ALTER TABLE `StateInfo` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Traffic`
--

DROP TABLE IF EXISTS `Traffic`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Traffic` (
  `TID` bigint(20) NOT NULL AUTO_INCREMENT,
  `PRID` bigint(20) NOT NULL DEFAULT '0',
  `FLAGS` bigint(20) NOT NULL DEFAULT '0',
  `Count` bigint(20) NOT NULL DEFAULT '0',
  `Description` varchar(128) NOT NULL DEFAULT '',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0',
  `CreateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreateBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`TID`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Traffic`
--

LOCK TABLES `Traffic` WRITE;
/*!40000 ALTER TABLE `Traffic` DISABLE KEYS */;
INSERT INTO `Traffic` VALUES (1,1,0,725,'Vehicles per day on Main street','2020-09-01 05:24:01',-99998,'2020-09-01 05:24:01',0),(2,1,0,1400,'Elm Street','2020-09-01 05:24:01',-99998,'2020-09-01 05:24:01',0);
/*!40000 ALTER TABLE `Traffic` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-01-20 17:09:23
