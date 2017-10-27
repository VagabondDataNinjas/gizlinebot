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

CREATE TABLE `answers` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `userId` varchar(255) NOT NULL DEFAULT '' COMMENT 'Line userId',
  `questionId` varchar(10) NOT NULL DEFAULT '' COMMENT 'The question Id',
  `answer` text NOT NULL COMMENT 'User entered answer',
  `timestamp` int(11) NOT NULL,
  `channel` varchar(10) NOT NULL DEFAULT '' COMMENT '"line" or "web" - the channel where the answer was receieved from',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

CREATE TABLE `answers_gps` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `userId` varchar(255) NOT NULL DEFAULT '' COMMENT 'Line userId',
  `lat` decimal(10,8) NOT NULL COMMENT 'Latitude',
  `lon` decimal(11,8) NOT NULL COMMENT 'Longitude',
  `timestamp` int(11) NOT NULL,
  `channel` varchar(10) NOT NULL DEFAULT '' COMMENT '"line" or "web" - the channel where the answer was receieved from',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Stores user profiles
CREATE TABLE `user_profiles` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `userId` varchar(255) NOT NULL DEFAULT '' COMMENT 'provided by Line during the follow event',
  `displayName` blob NOT NULL COMMENT 'Line field that allows users to set emojy''s and other non UTF-8 content',
  `timestamp` int(11) NOT NULL COMMENT 'UTC timestamp when this profile was added',
  `bot_survey_inited` tinyint(4) NOT NULL DEFAULT 0 COMMENT 'If set to 0: the line bot has not yet sent the initial question to the user. If non-0 linebot already sent the first question of the survey.',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`userId`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

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
	('job', 'What is your occupation?', -9, 'both'),
	('gps', 'What is your location', -8, 'web'),
	('island', 'What is the name of your island where you live?', -8, 'both'),
	('price', 'How much do you pay for diesel in your area?', -5, 'both'),
	('lineid', 'What is your line id?', -1, 'both'),
	('thank_you', 'Thank you for all your help! We might ask you more questions in the future', 0, 'both');

CREATE TABLE `welcome_msgs` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `msg` text NOT NULL,
  `weight` int(11) NOT NULL,
  `channel` varchar(10) NOT NULL DEFAULT '' COMMENT '"both", "web" or "line"',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

INSERT INTO `welcome_msgs` (`id`, `msg`, `weight`, `channel`)
VALUES
	(1, 'Thank you for following us!\nYou can find out more about us: https://www.youtube.com/watch?v=Vec5DML9yp4', -4, 'line'),
	(2, 'Please fill in the following survey in order to help our cause https://survey.delta9.link/?uid={{.UserId}}', -1, 'line');
