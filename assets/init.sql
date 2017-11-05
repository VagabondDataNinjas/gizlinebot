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
  `address` text NOT NULL COMMENT 'Certain channels (eg. line) provide an address field in the location message',
  `lat` decimal(10,8) NOT NULL COMMENT 'Latitude',
  `lon` decimal(11,8) NOT NULL COMMENT 'Longitude',
  `timestamp` int(11) NOT NULL,
  `channel` varchar(10) NOT NULL DEFAULT '' COMMENT '"line" or "web" - the channel where the answer was receieved from',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

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
	('job', 'คุณเป็นใคร?', -9, 'both'),
	('island', 'คุณอาศัยอยู่ที่เกาะอะไร', -8, 'both'),
	('price', 'น้ำมันดีเซล (โซลาร์) บนเกาะของคุณราคาลิตรละกี่บาท (ราคาที่ขายหรือซื้อจากร้านค้าหรือตัวแทนบนเกาะ)', -5, 'both'),
	('thank_you', 'ขอบคุณอีกครั้งที่ช่วยแบ่งปันข้อมูลให้เรา เราอยากเป็นเพื่อนกับคนที่อยู่บนเกาะให้ได้มากที่สุดจากหลายๆเกาะ โปรดช่วยส่งเราต่อไปให้เพื่อนคนอื่นๆในไลน์ของคุณด้วยวิธีง่ายๆเพียง\n- คลิก “v” ตรงมุมขวาบนของหน้านี้\n- คลิก “แนะนำ”\n- เลือกเพื่อนของคุณที่อยู่บนเกาะเดียวกันนี้หรือบนเกาะอื่นและคลิก “ตกลง”', 0, 'both');

CREATE TABLE `welcome_msgs` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `msg` text NOT NULL,
  `weight` int(11) NOT NULL,
  `channel` varchar(10) NOT NULL DEFAULT '' COMMENT '"both", "web" or "line"',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

INSERT INTO `welcome_msgs` (`id`, `msg`, `weight`, `channel`)
VALUES
	(1, 'สวัสดีจ้า ขอบคุณมากที่มาเป็นเพื่อนกัน เราชื่อกรู๊ดส์ เรากำลังทดลองเก็บข้อมูลราคาน้ำมันดีเซล (โซลาร์) จากคนบนเกาะในประเทศไทย ลองมาทำความรู้จักว่ากรู๊ดส์เป็นใครมาจากไหน อยากรู้ราคาน้ำมันไปทำอะไรและจะเป็นประโยชน์กับเกาะของคุณยังไงกันเลย (วิดีโอ)', -4, 'line'),
	(4, '{{.Hostname}}/media/groots_th.mp4|{{.Hostname}}/media/groots.png', -3, 'line'),
	(3, 'รู้อย่างงี้แล้วมาเริ่มแบ่งปันข้อมูลให้เรากันเลย ( {{.Hostname}}/?uid={{.UserId}} )', 0, 'line'),
	(7, 'Message shown only in test environment: you can remove all your profile and answer data at any time by accessing the following link: {{.Hostname}}/api/user/wipe/{{.UserId}}', 1, 'line');

CREATE TABLE `pricepoints` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL COMMENT '{user_profiles.id}',
  `price` float NOT NULL,
  `currency` varchar(20) NOT NULL DEFAULT 'thb',
  `location_id` int(11) NOT NULL COMMENT '{locations.id}',
  `time` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
INSERT INTO `pricepoints` (`id`, `user_id`, `price`, `currency`, `location_id`, `time`)
VALUES
	(1, 1, 35, 'thb', 230, 1509904925),
	(2, 1, 28, 'thb', 231, 1509904925),
	(3, 1, 27, 'thb', 232, 1509904925),
	(4, 1, 33, 'thb', 233, 1509904925),
	(5, 1, 38, 'thb', 234, 1509904925),
	(6, 1, 26, 'thb', 235, 1509904925),
	(7, 1, 31, 'thb', 236, 1509904925),
	(8, 1, 33, 'thb', 237, 1509904925);



CREATE TABLE `locations` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `thainame` varchar(255) DEFAULT NULL,
  `latitude` decimal(10,8) DEFAULT NULL,
  `longitude` decimal(11,8) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

INSERT INTO `locations` (`id`, `name`, `thainame`, `latitude`, `longitude`)
VALUES
	(1, 'Ko Chang', 'เกาะช้าง', 12.06000000, 102.33000000),
	(2, 'Ko Kut', 'เกาะกูด', 11.66000000, 102.57000000),
	(3, 'Ko Lanta Yai', 'เกาะลันตาใหญ่', NULL, NULL),
	(4, 'Ko Mak', 'เกาะหมาก', 11.81900000, 102.48000000),
	(5, 'Ko Pha Ngan', 'เกาะพะงัน', NULL, NULL),
	(6, 'Ko Phi Phi Lee', 'เกาะพีพีเล', 7.68000000, 98.76800000),
	(7, 'Ko Phi Phi Don', 'เกาะพีพีดอน', 7.74100000, 98.78000000),
	(8, 'Ko Rang', 'เกาะรัง', 11.80000000, 102.38700000),
	(9, 'Ko Samet', 'เกาะเสม็ด', 12.56500000, 101.45000000),
	(10, 'Ko Samui', 'เกาะสมุย', 9.50000000, 100.00000000),
	(11, 'Ko Tao', 'เกาะเต่า', 10.09000000, 99.83800000),
	(12, 'Ko Tapu', 'เกาะตะปู', NULL, NULL),
	(13, 'Ko Tarutao', 'เกาะตะรุเตา', NULL, NULL),
	(14, 'Ko Phuket', 'เกาะภูเก็ต', NULL, NULL),
	(15, 'Similan Islands', 'เกาะสิมิลัน', 8.65000000, 97.65000000),
	(16, 'Ko Sichang', 'เกาะสีชัง', 13.14500000, 100.81300000),
	(17, 'Ko Kham Yai', 'เกาะขามใหญ่', 13.16700000, 100.82400000),
	(18, 'Ko Kham Noi', 'เกาะขามน้อย', 13.17440000, 100.83010000),
	(19, 'Ko Ram Dok Mai', 'เกาะร้ามดอกไม้', 13.15280000, 100.83420000),
	(20, 'Ko Khang Khao', 'เกาะค้างคาว', 13.11400000, 100.80800000),
	(21, 'Ko Prong', 'เกาะปรง', 13.16520000, 100.83300000),
	(22, 'Ko Yai Thao', '', 13.12340000, 100.80930000),
	(23, 'Ko Thaai Tamuen', '', 13.10880000, 100.80330000),
	(24, 'Ko Nok', 'เกาะนอก', 13.00600000, 100.80330000),
	(25, 'Ko Chun', '', 12.93000000, 100.81000000),
	(26, 'Ko Lan', 'เกาะล้าน', 12.92000000, 100.78000000),
	(27, 'Ko Sak', 'เกาะสาก', 12.92750000, 100.79200000),
	(28, 'Ko Krok', 'เกาะครก', 12.94340000, 100.80570000),
	(29, 'Ko Luam', 'เกาะเหลื่อม', 12.95810000, 100.65060000),
	(30, 'Ko Luam Noi', 'เกาะเหลื่อมน้อย', 12.95660000, 100.65730000),
	(31, 'Ko Phai', 'เกาะไผ่', 12.93400000, 100.67500000),
	(32, 'Ko Hu Chang', 'เกาะหูช้าง', 12.91050000, 100.67750000),
	(33, 'Ko Klung Badan', 'เกาะกรุงบาดาล', 12.90260000, 100.67960000),
	(34, 'Ko Man Wichai', 'เกาะมารวิชัย', 12.87630000, 100.67420000),
	(35, 'Ko Klet Kaeo', 'เกาะเกล็ดแก้ว', 12.76100000, 100.84500000),
	(36, 'Ko Khram', 'เกาะคราม', 12.70000000, 100.79000000),
	(37, 'Ko Khram Noi', 'เกาะครามน้อย', 12.72700000, 100.79800000),
	(38, 'Ko I Ra', 'เกาะอีร้า', 12.67600000, 100.82300000),
	(39, 'Ko Tao Mo', 'เกาะเตาหม้อ', NULL, NULL),
	(40, 'Ko Phra', 'เกาะพระ', NULL, NULL),
	(41, 'Ko Phra Noi', 'เกาะพระน้อย', NULL, NULL),
	(42, 'Ko I Lao', 'เกาะอีเลา', NULL, NULL),
	(43, 'Ko Yo', 'เกาะยอ', NULL, NULL),
	(44, 'Ko Mu', 'เกาะหมู', NULL, NULL),
	(45, 'Ko Maeo', 'เกาะแมว', NULL, NULL),
	(46, 'Ko Nang Ram', 'เกาะนางรำ', NULL, NULL),
	(47, 'Ko Chorakhe', 'เกาะจระเข้', NULL, NULL),
	(48, 'Ko Samae San', 'เกาะแสมสาร', 12.57300000, 100.95000000),
	(49, 'Ko Kham', 'เกาะขาม', 12.57300000, 100.93400000),
	(50, 'Ko Raet', 'เกาะแรด', 12.58500000, 100.96400000),
	(51, 'Ko Chang Kluea', '', 12.55000000, 100.97100000),
	(52, 'Ko Chuang', 'เกาะช่วง', 12.52000000, 100.96000000),
	(53, 'Ko Chan', 'เกาะจาน', 12.52100000, 100.97000000),
	(54, 'Ko Rong Nang', 'เกาะโรงหนัง', 12.53437000, 100.95870000),
	(55, 'Ko Rong Khon', 'เกาะโรงโขน', 12.53169000, 100.95972000),
	(56, 'Ko Saket', '', NULL, NULL),
	(58, 'Ko Platin', '', 12.60190000, 101.51480000),
	(59, 'Ko Kraui', '', 12.59320000, 101.51080000),
	(60, 'Ko Kudi', '', 12.58390000, 101.51010000),
	(61, 'Ko Khangkhao', '', 12.57960000, 101.51040000),
	(63, 'Ko Yangklao', '', 12.55780000, 101.56910000),
	(64, 'Ko Thalu', '', 12.55570000, 101.56950000),
	(65, 'Ko Khi Pla', 'เกาะขี้ปลา', NULL, NULL),
	(66, 'Ko Man Nai', 'เกาะมันใน', 12.61140000, 101.68800000),
	(67, 'Ko Man Khlang', 'เกาะมันกลาง', 12.59750000, 101.69200000),
	(68, 'Ko Man Nok', 'เกาะมันนอก', 12.58410000, 101.69700000),
	(69, 'Ko Chong Saba', 'เกาะช่องสะบ้า', NULL, NULL),
	(70, 'Ko Nom Sao', 'เกาะนมสาว', NULL, NULL),
	(71, 'Ko Chula', 'เกาะจุฬา', NULL, NULL),
	(72, 'Ko Nu', 'เกาะหนู', NULL, NULL),
	(74, 'Ko Kwang', 'เกาะกวาง', NULL, NULL),
	(75, 'Ko Chik Nok', 'เกาะจิกนอก', NULL, NULL),
	(76, 'Ko Chik Klang', 'เกาะจิกกลาง', NULL, NULL),
	(78, 'Ko Ngam', 'เกาะง่าม', 11.95100000, 102.44300000),
	(79, 'Ko Phrao Nai', 'เกาะพร้าวใน', 11.98200000, 102.38700000),
	(80, 'Ko Phrao Nok', 'เกาะพร้าวนอก', 11.97300000, 102.39200000),
	(81, 'Ko Klum', 'เกาะคลุ้ม', 11.91400000, 102.35400000),
	(82, 'Mu Ko Lao Ya', 'หมู่เกาะเหลายา', 11.94200000, 102.40700000),
	(83, 'Ko Lao Ya Nai', 'เกาะเหลายาใน', 11.94200000, 102.40700000),
	(84, 'Ko Lao Ya Nok', 'เกาะเหลายานอก', 11.93100000, 102.41100000),
	(85, 'Ko Wai', 'เกาะหวาย', 11.90000000, 102.40500000),
	(86, 'Mu Ko Mai Si', 'หมู่เกาะไม้ซี้', 11.94600000, 102.47500000),
	(87, 'Ko Mai Si Yai', 'เกาะไม้ซี้ใหญ่', 11.94600000, 102.47500000),
	(88, 'Ko Mai Si Lek', 'เกาะไม้ซี้เล็ก', 11.95500000, 102.49700000),
	(89, 'Ko Bai Dang', 'เกาะไปแดง', 11.89800000, 102.45100000),
	(93, 'Ko Rayang Nok', 'เกาะระยั้งนอก', 11.79800000, 102.45100000),
	(94, 'Ko Kradat', 'เกาะกระดาด', 11.84000000, 102.52750000),
	(95, 'Mu Ko Rang', 'หมู่เกาะรัง', 11.80000000, 102.38700000),
	(97, 'Ko Tun', 'เกาะตูน', 11.77400000, 102.39300000),
	(98, 'Ko Phing', 'เกาะพิง', NULL, NULL),
	(99, 'Ko Phang', 'เกาะพัง', NULL, NULL),
	(100, 'Ko Sai', 'เกาะทราย', 12.48300000, 100.00000000),
	(101, 'Ko Sadao', 'เกาะสะเดา', 12.46600000, 100.00000000),
	(102, 'Ko Khi Nok', 'เกาะขี้นก', 12.48300000, 100.00000000),
	(103, 'Ko Lak', 'เกาะหลัก', NULL, NULL),
	(104, 'Ko La', 'เกาะละ', NULL, NULL),
	(105, 'Ko Rom', 'เกาะร่ม', NULL, NULL),
	(108, 'Ko Thai Chan', 'เกาะท้ายจาน', NULL, NULL),
	(110, 'Ko Kho Ram', 'เกาะโครำ', NULL, NULL),
	(111, 'Ko Rawang', 'เกาะระวาง', NULL, NULL),
	(112, 'Ko Rawing', 'เกาะระวิง', NULL, NULL),
	(113, 'Ko Sattakut', 'เกาะสัตกูต', 12.20000000, 100.01600000),
	(114, 'Ko Hua Pin', 'เกาะหัวพิน', NULL, NULL),
	(116, 'Ko Sing', 'เกาะสิงข์', NULL, NULL),
	(117, 'Ko Sang', 'เกาะสังข์', NULL, NULL),
	(118, 'Ko Ram Ra', 'เกาะร่ำรา', NULL, NULL),
	(119, 'Ko Lueam', 'เกาะเหลื่อม', NULL, NULL),
	(120, 'Ko Aen', 'เกาะแอ่น', NULL, NULL),
	(121, 'Ko Wiang', 'เกาะเวียง', NULL, NULL),
	(125, 'Ko Si Kong', 'เกาะซีกง', NULL, NULL),
	(127, 'Ko Khai', 'เกาะไข่', NULL, NULL),
	(129, 'Ko Ngam Yai', 'เกาะงามใหญ่', NULL, NULL),
	(130, 'Ko Ngam Noi', 'เกาะงามน้อย', NULL, NULL),
	(131, 'Ko Kalok', 'เกาะกะโหลก', NULL, NULL),
	(134, 'Ko Maphrao', 'เกาะมะพร้าว', NULL, NULL),
	(135, 'Ko Matra', 'เกาะมาตรา', NULL, NULL),
	(136, 'Ko Lak Raet', 'เกาะหลักแรด', NULL, NULL),
	(137, 'Ko I Raet', 'เกาะอีแรด', NULL, NULL),
	(138, 'Ko Lawa', 'เกาะละวะ', NULL, NULL),
	(139, 'Ko Ka', 'เกาะกา', NULL, NULL),
	(140, 'Ko Thong Lang', 'เกาะทองหลาง', NULL, NULL),
	(141, 'Ko Wang Ka Chio', 'เกาะวังกะจิว', NULL, NULL),
	(142, 'Ko Kra', '', 10.25000000, 99.25000000),
	(143, 'Ko Kula', '', 10.25000000, 99.25000000),
	(146, 'Ko Nang Yuan', '', 10.09000000, 99.83800000),
	(148, 'Ko Similan', 'เกาะสิมิลัน', 8.65000000, 97.65000000),
	(149, 'Ko Bangu', '', 8.65400000, 97.65000000),
	(150, 'Ko Payu', '', 8.59000000, 97.63800000),
	(151, 'Ko Miang', '', 8.60000000, 97.63700000),
	(152, 'Ko Payang', '', 8.49800000, 97.64600000),
	(153, 'Ko Huyang', '', 8.47900000, 97.65000000),
	(156, 'Ko Pu', 'เกาะปู', NULL, NULL),
	(157, 'Ko Siboya', 'เกาะปู', NULL, NULL),
	(158, 'Ko Jum', 'เกาะปู', NULL, NULL),
	(159, 'Ko Kai, Krabi', 'เกาะไก่', NULL, NULL),
	(160, 'Ko Pli', 'เกาะไก่', NULL, NULL),
	(161, 'Ko Lek', 'เกาะไก่', NULL, NULL),
	(162, 'Ko Lapu Le', 'เกาะไก่', NULL, NULL),
	(163, 'Ko A Dang', 'เกาะไก่', NULL, NULL),
	(164, 'Ko Poda', '', NULL, NULL),
	(166, 'Ko Lanta Noi', 'เกาะลันตาน้อย', NULL, NULL),
	(167, 'Ko Bu Bu', '', NULL, NULL),
	(168, 'Ko Po', '', NULL, NULL),
	(169, 'Ko Klang (Krabi)', '', NULL, NULL),
	(170, 'Ko Ma', '', NULL, NULL),
	(171, 'Ko Ha', '', NULL, NULL),
	(172, 'Phuket', 'ภูเก็ต', NULL, NULL),
	(173, 'Ko Lon', 'เกาะโหลน', NULL, NULL),
	(174, 'Ko Mai Thon', 'เกาะไม้ท่อน', NULL, NULL),
	(176, 'Ko Ngam (Phuket)', 'เกาะงำ', NULL, NULL),
	(177, 'Ko Naka Yai', 'เกาะนาคาใหญ่', NULL, NULL),
	(178, 'Ko Naka Noi', 'เกาะนาคาน้อย', NULL, NULL),
	(180, 'Ko Siray', 'เกาะสิเหร่', NULL, NULL),
	(181, 'Ko Rang Yai', 'เกาะรังใหญ่', NULL, NULL),
	(182, 'Ko Rang Noi', 'เกาะรังน้อย', NULL, NULL),
	(183, 'Ko Tapao', 'เกาะตะเภา', NULL, NULL),
	(184, 'Ko Kaeo Yai', 'เกาะแก้วใหญ่', NULL, NULL),
	(185, 'Ko Kaeo', 'เกาะแก้ว', NULL, NULL),
	(186, 'Ko Bon', 'เกาะบอน', NULL, NULL),
	(187, 'Ko Hae', 'เกาะเฮ', NULL, NULL),
	(188, 'Ko Aeo', 'เกาะแอว', NULL, NULL),
	(224, 'Ko Ngai', 'เกาะไหง', 7.41481100, 99.15485800),
	(225, 'Ko Kradan', 'เกาะกระดาน', 7.30857700, 99.25544100),
	(226, 'Ko Muk', 'เกาะมุก', 7.37331200, 99.29551600),
	(227, 'Ko Racha Yai', 'เกาะราชาใหญ่', 7.59846600, 98.36635300),
	(228, 'Ko Lamphu', 'เกาะลำพู', 9.13500000, 99.30100000),
	(229, 'Ko Lamphu Rai', 'เกาะลำพูราย', 12.20000000, 102.58300000),
	(230, 'Ko Panyi', 'เกาะปันหยี', 8.33543000, 98.50453000),
	(231, 'Ko Lon', 'เกาะโหลน', 7.78651700, 98.37184600),
	(232, 'Ko Naga Noi', 'เกาะ นาคา น้อย', 8.02819100, 98.46050400),
	(233, 'Ko Phra Thong', 'เกาะพระทอง', 9.09319600, 98.29515300),
	(234, 'Ko Por', 'เกาะปอ', 7.53622900, 99.12782200),
	(235, 'Ko Bulon Don', 'เกะดอน', 6.85621000, 99.59314000),
	(236, 'Ko Mak Noi', 'เกาะหมากน้อย', 8.28506160, 98.59035210),
	(237, 'Ko Mai Phai', 'เกาะไม้ไผ่', 7.81670000, 98.79546000);
