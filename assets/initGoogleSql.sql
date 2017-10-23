CREATE database if not exists gizsurvey;
USE gizsurvey;

-- CREATE USER if not exists 'gizlinebot'@'127.0.0.1' IDENTIFIED BY 'Mpmc3EzwUU06Pq9hq8T55fEnaN2okglRd5CPS2i4fcA';
-- GRANT ALL PRIVILEGES ON gizsurvey.* TO 'gizlinebot'@'127.0.0.1';
-- FLUSH PRIVILEGES;

-- Stores all the messages sent by a user
-- useful in case we need to re-parse / validate messages
CREATE TABLE IF NOT EXISTS `linebot_raw_events` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `eventtype` varchar(255) NOT NULL,
  `rawevent` text DEFAULT NULL,
  `timestamp` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- Stores all answers sent by a user
CREATE TABLE IF NOT EXISTS `linebot_answers` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `userId` varchar(255) NOT NULL DEFAULT '' COMMENT 'Line userId',
  `questionId` varchar(10) NOT NULL DEFAULT '' COMMENT 'The question Id',
  `answer` text NOT NULL COMMENT 'User entered answer',
  `timestamp` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Stores user profiles. Useful in case we need to contact the users manually
CREATE TABLE IF NOT EXISTS `user_profiles` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `userId` varchar(255) NOT NULL DEFAULT '' COMMENT 'provided by Line during the follow event',
  `displayName` varchar(255) NOT NULL DEFAULT '' COMMENT 'provided by Line during the follow event',
  `timestamp` int(11) NOT NULL COMMENT 'UTC timestamp when this profile was added',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_display_name` (`userId`,`displayName`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
