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
    `image` TEXT NOT NULL
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4;