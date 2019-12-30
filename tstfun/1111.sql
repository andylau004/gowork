-- MySQL dump 10.13  Distrib 5.7.27, for Linux (x86_64)
--
-- Host: localhost    Database: IM_YUN_FILE
-- ------------------------------------------------------
-- Server version	5.7.24-log

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
-- Table structure for table `IM_FILE`
--

DROP TABLE IF EXISTS `IM_FILE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `IM_FILE` (
  `fileId` bigint(20) NOT NULL DEFAULT '0' COMMENT '文件编号',
  `userId` bigint(20) NOT NULL,
  `fileName` varchar(100) CHARACTER SET utf8mb4 DEFAULT NULL COMMENT '文件名称 不含后缀',
  `fileType` int(11) DEFAULT '0' COMMENT '文件类型 1 文件 2 文件夹 3 团队文件夹 4面向全部的团队文件夹',
  `fileSuffix` varchar(20) COLLATE utf8_bin DEFAULT NULL COMMENT '文件后缀名',
  `filePath` varchar(500) COLLATE utf8_bin DEFAULT NULL COMMENT '文件在文件服务器上的路径',
  `fileSize` bigint(20) DEFAULT '0' COMMENT '文件大小',
  `encryptFileSize` bigint(20) DEFAULT '0' COMMENT '加密后文件大小',
  `secretKey` varchar(50) COLLATE utf8_bin DEFAULT NULL COMMENT '加密密钥',
  `uploaderId` bigint(20) DEFAULT '0' COMMENT '上传者id',
  `uploaderName` varchar(50) COLLATE utf8_bin DEFAULT NULL COMMENT '上传者名称',
  `md5` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `sha2` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `createdAt` bigint(20) DEFAULT '0' COMMENT '创建时间',
  `updatedAt` bigint(20) DEFAULT '0' COMMENT '修改时间',
  `deletedAt` bigint(20) DEFAULT '0' COMMENT '删除时间',
  `fileCode` varchar(250) COLLATE utf8_bin DEFAULT NULL COMMENT '文件夹code: 1001',
  `orderNum` bigint(20) DEFAULT '0' COMMENT '排序',
  `roleId` int(11) DEFAULT '0' COMMENT '权限(角色编号)',
  `members` int(11) DEFAULT '0' COMMENT '成员数',
  `pdfPath` varchar(500) COLLATE utf8_bin DEFAULT NULL COMMENT '预览文件地址',
  `validity` int(11) DEFAULT '0' COMMENT '有效期',
  `validityAt` bigint(20) DEFAULT '0' COMMENT '文件到期时间',
  `realFileName` varchar(500) COLLATE utf8_bin DEFAULT NULL,
  `timeStamp` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`fileId`,`userId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='文件表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `IM_FILE_MEMBER`
--

DROP TABLE IF EXISTS `IM_FILE_MEMBER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `IM_FILE_MEMBER` (
  `fileId` bigint(20) NOT NULL DEFAULT '0' COMMENT '文件编号',
  `userId` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户编号',
  `roleId` bigint(20) DEFAULT '0' COMMENT '权限(角色编号)',
  `createdAt` bigint(20) DEFAULT '0' COMMENT '创建时间',
  `updatedAt` bigint(20) DEFAULT '0' COMMENT '修改时间',
  `userName` varchar(50) COLLATE utf8_bin DEFAULT NULL COMMENT '用户名称',
  `headImg` varchar(500) COLLATE utf8_bin DEFAULT NULL COMMENT '用户头像',
  `orgInfo` varchar(255) COLLATE utf8_bin DEFAULT NULL COMMENT ' ',
  PRIMARY KEY (`fileId`,`userId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='文件成员表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `IM_FILE_STAT`
--

DROP TABLE IF EXISTS `IM_FILE_STAT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `IM_FILE_STAT` (
  `timeAt` bigint(20) NOT NULL DEFAULT '0' COMMENT '时间',
  `uploadCount` bigint(20) DEFAULT '0' COMMENT '上传次数',
  `downloadCount` bigint(20) DEFAULT '0' COMMENT '下载次数',
  PRIMARY KEY (`timeAt`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='统计表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `IM_FILE_STAT2`
--

DROP TABLE IF EXISTS `IM_FILE_STAT2`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `IM_FILE_STAT2` (
  `timeAt` bigint(20) NOT NULL DEFAULT '0' COMMENT '时间',
  `fileSuffix` varchar(20) COLLATE utf8_bin DEFAULT NULL COMMENT '文件后缀',
  `uploadCount` bigint(20) DEFAULT '0' COMMENT '上传次数',
  `downloadCount` bigint(20) DEFAULT '0' COMMENT '下载次数',
  PRIMARY KEY (`timeAt`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='文件统计表2';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `IM_ROLE`
--

DROP TABLE IF EXISTS `IM_ROLE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `IM_ROLE` (
  `roleID` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '角色编号 (上传,下载,创建团队,删除,重命名)',
  `name` varchar(50) COLLATE utf8_bin DEFAULT NULL COMMENT '角色名称',
  `permission` varchar(50) COLLATE utf8_bin DEFAULT NULL COMMENT '角色权限',
  `uploadSize` bigint(20) DEFAULT '0' COMMENT '上传大小限制  -1表示无上限',
  `totalSize` bigint(20) DEFAULT '0' COMMENT '文件夹总空间限制  -1表示无上限',
  `createdAt` bigint(20) DEFAULT '0' COMMENT '创建时间',
  `updatedAt` bigint(20) DEFAULT '0' COMMENT '修改时间',
  `orderNum` bigint(20) DEFAULT '0' COMMENT '排序参数',
  `isShow` int(11) DEFAULT '0' COMMENT '是否展示在前端 1是 0否',
  `roleType` int(11) DEFAULT '2' COMMENT '角色类型  1管理员  2普通成员  3其他角色',
  `status` int(11) DEFAULT '0' COMMENT '状态  1 正常  2删除',
  PRIMARY KEY (`roleID`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='角色表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `IM_TMP_FILE_0`
--

DROP TABLE IF EXISTS `IM_TMP_FILE_0`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `IM_TMP_FILE_0` (
  `fileId` bigint(20) NOT NULL COMMENT '文件编号',
  `pid` bigint(20) NOT NULL DEFAULT '0' COMMENT '上级目录编号',
  `ownerId` bigint(20) NOT NULL DEFAULT '0' COMMENT '所有者编号',
  `fileName` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件名称 不含后缀',
  `fileType` int(11) NOT NULL DEFAULT '0' COMMENT '文件类型 1 文件 2 文件夹 3 团队文件夹 4面向全部的团队文件夹',
  `fileSuffix` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件后缀名',
  `filePath` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件在文件服务器上的路径',
  `fileSize` bigint(20) DEFAULT NULL COMMENT '文件大小',
  `encryptFileSize` bigint(20) DEFAULT NULL COMMENT '加密后文件大小',
  `secretKey` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '加密密钥',
  `uploaderId` bigint(20) DEFAULT NULL COMMENT '上传者id',
  `uploaderName` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '上传者名称',
  `md5Hash` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
  `sha1Hash` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
  `createdAt` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updatedAt` bigint(20) DEFAULT NULL COMMENT '修改时间',
  `deletedAt` bigint(20) DEFAULT NULL COMMENT '删除时间',
  `fileCode` varchar(2000) COLLATE utf8mb4_bin DEFAULT NULL,
  `orderNum` bigint(20) DEFAULT NULL COMMENT '排序',
  `roleId` int(11) DEFAULT NULL COMMENT '权限(角色编号)',
  `members` int(11) DEFAULT '0' COMMENT '成员数',
  `fileReMark` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `pdfPath` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '预览文件地址',
  `validity` int(11) DEFAULT '0' COMMENT '有效期',
  `fileClass` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件分类',
  `validityAt` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`fileId`,`pid`,`ownerId`),
  KEY `pOwnerId` (`pid`,`ownerId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='临时文件表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `IM_TMP_FILE_1`
--

DROP TABLE IF EXISTS `IM_TMP_FILE_1`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `IM_TMP_FILE_1` (
  `fileId` bigint(20) NOT NULL COMMENT '文件编号',
  `pid` bigint(20) NOT NULL DEFAULT '0' COMMENT '上级目录编号',
  `ownerId` bigint(20) NOT NULL DEFAULT '0' COMMENT '所有者编号',
  `fileName` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件名称 不含后缀',
  `fileType` int(11) NOT NULL DEFAULT '0' COMMENT '文件类型 1 文件 2 文件夹 3 团队文件夹 4面向全部的团队文件夹',
  `fileSuffix` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件后缀名',
  `filePath` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件在文件服务器上的路径',
  `fileSize` bigint(20) DEFAULT NULL COMMENT '文件大小',
  `encryptFileSize` bigint(20) DEFAULT NULL COMMENT '加密后文件大小',
  `secretKey` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '加密密钥',
  `uploaderId` bigint(20) DEFAULT NULL COMMENT '上传者id',
  `uploaderName` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '上传者名称',
  `md5Hash` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
  `sha1Hash` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
  `createdAt` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updatedAt` bigint(20) DEFAULT NULL COMMENT '修改时间',
  `deletedAt` bigint(20) DEFAULT NULL COMMENT '删除时间',
  `fileCode` varchar(2000) COLLATE utf8mb4_bin DEFAULT NULL,
  `orderNum` bigint(20) DEFAULT NULL COMMENT '排序',
  `roleId` int(11) DEFAULT NULL COMMENT '权限(角色编号)',
  `members` int(11) DEFAULT '0' COMMENT '成员数',
  `fileReMark` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `pdfPath` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '预览文件地址',
  `validity` int(11) DEFAULT '0' COMMENT '有效期',
  `fileClass` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件分类',
  `validityAt` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`fileId`,`pid`,`ownerId`),
  KEY `pOwnerId` (`pid`,`ownerId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='临时文件表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `IM_TMP_FILE_2`
--

DROP TABLE IF EXISTS `IM_TMP_FILE_2`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `IM_TMP_FILE_2` (
  `fileId` bigint(20) NOT NULL COMMENT '文件编号',
  `pid` bigint(20) NOT NULL DEFAULT '0' COMMENT '上级目录编号',
  `ownerId` bigint(20) NOT NULL DEFAULT '0' COMMENT '所有者编号',
  `fileName` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件名称 不含后缀',
  `fileType` int(11) NOT NULL DEFAULT '0' COMMENT '文件类型 1 文件 2 文件夹 3 团队文件夹 4面向全部的团队文件夹',
  `fileSuffix` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件后缀名',
  `filePath` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件在文件服务器上的路径',
  `fileSize` bigint(20) DEFAULT NULL COMMENT '文件大小',
  `encryptFileSize` bigint(20) DEFAULT NULL COMMENT '加密后文件大小',
  `secretKey` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '加密密钥',
  `uploaderId` bigint(20) DEFAULT NULL COMMENT '上传者id',
  `uploaderName` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '上传者名称',
  `md5Hash` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
  `sha1Hash` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
  `createdAt` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updatedAt` bigint(20) DEFAULT NULL COMMENT '修改时间',
  `deletedAt` bigint(20) DEFAULT NULL COMMENT '删除时间',
  `fileCode` varchar(2000) COLLATE utf8mb4_bin DEFAULT NULL,
  `orderNum` bigint(20) DEFAULT NULL COMMENT '排序',
  `roleId` int(11) DEFAULT NULL COMMENT '权限(角色编号)',
  `members` int(11) DEFAULT '0' COMMENT '成员数',
  `fileReMark` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `pdfPath` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '预览文件地址',
  `validity` int(11) DEFAULT '0' COMMENT '有效期',
  `fileClass` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件分类',
  `validityAt` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`fileId`,`pid`,`ownerId`),
  KEY `pOwnerId` (`pid`,`ownerId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='临时文件表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `IM_TMP_FILE_3`
--

DROP TABLE IF EXISTS `IM_TMP_FILE_3`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `IM_TMP_FILE_3` (
  `fileId` bigint(20) NOT NULL COMMENT '文件编号',
  `pid` bigint(20) NOT NULL DEFAULT '0' COMMENT '上级目录编号',
  `ownerId` bigint(20) NOT NULL DEFAULT '0' COMMENT '所有者编号',
  `fileName` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件名称 不含后缀',
  `fileType` int(11) NOT NULL DEFAULT '0' COMMENT '文件类型 1 文件 2 文件夹 3 团队文件夹 4面向全部的团队文件夹',
  `fileSuffix` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件后缀名',
  `filePath` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件在文件服务器上的路径',
  `fileSize` bigint(20) DEFAULT NULL COMMENT '文件大小',
  `encryptFileSize` bigint(20) DEFAULT NULL COMMENT '加密后文件大小',
  `secretKey` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '加密密钥',
  `uploaderId` bigint(20) DEFAULT NULL COMMENT '上传者id',
  `uploaderName` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '上传者名称',
  `md5Hash` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
  `sha1Hash` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
  `createdAt` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updatedAt` bigint(20) DEFAULT NULL COMMENT '修改时间',
  `deletedAt` bigint(20) DEFAULT NULL COMMENT '删除时间',
  `fileCode` varchar(2000) COLLATE utf8mb4_bin DEFAULT NULL,
  `orderNum` bigint(20) DEFAULT NULL COMMENT '排序',
  `roleId` int(11) DEFAULT NULL COMMENT '权限(角色编号)',
  `members` int(11) DEFAULT '0' COMMENT '成员数',
  `fileReMark` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `pdfPath` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '预览文件地址',
  `validity` int(11) DEFAULT '0' COMMENT '有效期',
  `fileClass` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件分类',
  `validityAt` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`fileId`,`pid`,`ownerId`),
  KEY `pOwnerId` (`pid`,`ownerId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='临时文件表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `IM_TMP_FILE_4`
--

DROP TABLE IF EXISTS `IM_TMP_FILE_4`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `IM_TMP_FILE_4` (
  `fileId` bigint(20) NOT NULL COMMENT '文件编号',
  `pid` bigint(20) NOT NULL DEFAULT '0' COMMENT '上级目录编号',
  `ownerId` bigint(20) NOT NULL DEFAULT '0' COMMENT '所有者编号',
  `fileName` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件名称 不含后缀',
  `fileType` int(11) NOT NULL DEFAULT '0' COMMENT '文件类型 1 文件 2 文件夹 3 团队文件夹 4面向全部的团队文件夹',
  `fileSuffix` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件后缀名',
  `filePath` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件在文件服务器上的路径',
  `fileSize` bigint(20) DEFAULT NULL COMMENT '文件大小',
  `encryptFileSize` bigint(20) DEFAULT NULL COMMENT '加密后文件大小',
  `secretKey` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '加密密钥',
  `uploaderId` bigint(20) DEFAULT NULL COMMENT '上传者id',
  `uploaderName` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '上传者名称',
  `md5Hash` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
  `sha1Hash` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
  `createdAt` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updatedAt` bigint(20) DEFAULT NULL COMMENT '修改时间',
  `deletedAt` bigint(20) DEFAULT NULL COMMENT '删除时间',
  `fileCode` varchar(2000) COLLATE utf8mb4_bin DEFAULT NULL,
  `orderNum` bigint(20) DEFAULT NULL COMMENT '排序',
  `roleId` int(11) DEFAULT NULL COMMENT '权限(角色编号)',
  `members` int(11) DEFAULT '0' COMMENT '成员数',
  `fileReMark` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `pdfPath` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '预览文件地址',
  `validity` int(11) DEFAULT '0' COMMENT '有效期',
  `fileClass` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件分类',
  `validityAt` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`fileId`,`pid`,`ownerId`),
  KEY `pOwnerId` (`pid`,`ownerId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='临时文件表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `IM_USER`
--

DROP TABLE IF EXISTS `IM_USER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `IM_USER` (
  `userId` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `createdAt` bigint(20) DEFAULT '0' COMMENT '创建时间',
  `updatedAt` bigint(20) DEFAULT '0' COMMENT '修改时间',
  `useSpace` bigint(255) DEFAULT '0' COMMENT '已使用容量(Bit)',
  `status` int(10) DEFAULT '1' COMMENT '状态:1未升级,2升级中,3正常',
  PRIMARY KEY (`userId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='用户信息表';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2019-09-05 21:15:52
