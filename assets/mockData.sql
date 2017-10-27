CREATE database IF NOT EXISTS gizsurvey;
USE gizsurvey;

CREATE TABLE `gizsurvey`.`mock_user_profiles` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `userId` VARCHAR(45) NOT NULL,
  `userLineId` VARCHAR(45) NULL default '100',
  `displayName` VARCHAR(45) NULL default 'A',
  `dateAddedBot` DATETIME NOT NULL,
  `dateRemovedBot` DATETIME NOT NULL,
  `timestamp` DATETIME NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `gizsurvey`.`mock_kpi_info_data` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `questionId` INT(11) NOT NULL,
  `islandName` VARCHAR(45) NOT NULL,
  `job` VARCHAR(45) NOT NULL,
  `dieselPrice` DOUBLE NOT NULL,
  `latitude` DOUBLE NULL,
  `longtitude` DOUBLE NULL,
  `timestamp` DATETIME NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `gizsurvey`.`mock_question_table` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `questionId` INT(11) NOT NULL,
  `userId` VARCHAR(45) NOT NULL,
  `dateSent` DATETIME NOT NULL,
  `dateCompleted` DATETIME NOT NULL,
  `timestamp` DATETIME NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;