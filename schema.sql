CREATE TABLE `tasks` (
    `id` INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `content` TEXT NOT NULL,
    `yomi` TEXT NOT NULL,
    `iconUri` TEXT NOT NULL,
    `authorDisplayName` TEXT NOT NULL,
    `grade` TEXT NOT NULL,
    `authorName` TEXT NOT NULL,
    `updatedAt` DATETIME NOT NULL,
    `level` INT(11) NOT NULL,
    `isSensitive` BOOLEAN NOT NULL,
    `citated` TEXT NOT NULL,
    `image` TEXT NOT NULL,
    `messageId` TEXT NOT NULL
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `stamps` (
    `taskId` TEXT NOT NULL,
    `stampId` TEXT NOT NULL,
    `count` int(11) NOT NULL
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `rankings` (
    `id` INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `userName` TEXT NOT NULL,
    `score` DOUBLE NOT NULL,
    `level` INT(11) NOT NULL,
    `timeStamp` DATETIME NOT NULL
 ) ENGINE=INNODB DEFAULT CHARSET=utf8mb4;