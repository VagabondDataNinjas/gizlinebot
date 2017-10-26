CREATE database if not exists gizsurvey;
USE gizsurvey;


CREATE USER if not exists 'gizlinebot'@'127.0.0.1' IDENTIFIED BY 'Mpmc3EzwUU06Pq9hq8T55fEnaN2okglRd5CPS2i4fcA';
GRANT ALL PRIVILEGES ON gizsurvey.* TO 'gizlinebot'@'127.0.0.1';
FLUSH PRIVILEGES;

-- Stores all the messages sent by a user
-- useful in case we need to re-parse / validate messages
CREATE TABLE `linebot_raw_events` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `eventtype` varchar(255) NOT NULL,
  `rawevent` text DEFAULT NULL,
  `timestamp` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- Stores all answers sent by a user
CREATE TABLE `answers` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `userId` varchar(255) NOT NULL DEFAULT '' COMMENT 'Line userId',
  `questionId` varchar(10) NOT NULL DEFAULT '' COMMENT 'The question Id',
  `answer` text NOT NULL COMMENT 'User entered answer',
  `timestamp` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `answers_gps` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `userId` varchar(255) NOT NULL DEFAULT '' COMMENT 'Line userId',
  `lat` decimal(10,8) NOT NULL COMMENT 'Latitude',
  `lon` decimal(11,8) NOT NULL COMMENT 'Longitude',
  `timestamp` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Stores user profiles. Useful in case we need to contact the users manually
CREATE TABLE `user_profiles` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `userId` varchar(255) NOT NULL DEFAULT '' COMMENT 'provided by Line during the follow event',
  `displayName` varchar(255) NOT NULL DEFAULT '' COMMENT 'provided by Line during the follow event',
  `timestamp` int(11) NOT NULL COMMENT 'UTC timestamp when this profile was added',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_display_name` (`userId`,`displayName`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

-- list of questions for the users
CREATE TABLE `questions` (
  `id` varchar(255) NOT NULL DEFAULT '',
  `question` text NOT NULL,
  `weight` int(11) NOT NULL,
  `channel` varchar(10) NOT NULL DEFAULT '' COMMENT '"both", "web" or "line"',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `questions` (`id`, `question`, `weight`, `channel`)
VALUES
	('welcome', 'Thank you for following us!\nYou can find out more about us: https://www.youtube.com/watch?v=Vec5DML9yp4\nPlease fill in the following survey in order to help our cause https://survey.delta9.link?uid={{.UserId}}', -10, 'both'),
	('job', 'What is your occupation?', -9, 'both'),
	('gps', 'What is your location', -8, 'web'),
	('island', 'What is the name of your island where you live?', -8, 'both'),
	('price', 'How much do you pay for diesel in your area?', -5, 'both'),
	('lineid', 'What is your line id?', -1, 'both'),
	('thank_you', 'Thank you for all your help! We might ask you more questions in the future', 0, 'both');
