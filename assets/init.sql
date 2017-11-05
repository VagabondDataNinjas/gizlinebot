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
  `currency` varchar(20) NOT NULL DEFAULT '',
  `location_id` int(11) NOT NULL COMMENT '{locations.id}',
  `time` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `locations` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `thainame` varchar(255) DEFAULT NULL,
  `geolocation` point DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

INSERT INTO `locations` (`id`, `name`, `thainame`, `geolocation`)
VALUES
	(1, 'Ko Chang', 'เกาะช้าง', X'0000000001010000001F85EB51B81E284085EB51B81E955940'),
	(2, 'Ko Kut', 'เกาะกูด', X'00000000010100000052B81E85EB51274014AE47E17AA45940'),
	(3, 'Ko Lanta Yai', 'เกาะลันตาใหญ่', NULL),
	(4, 'Ko Mak', 'เกาะหมาก', X'00000000010100000017D9CEF753A327401F85EB51B89E5940'),
	(5, 'Ko Pha Ngan', 'เกาะพะงัน', NULL),
	(6, 'Ko Phi Phi Lee', 'เกาะพีพีเล', X'000000000101000000B81E85EB51B81E40FED478E926B15840'),
	(7, 'Ko Phi Phi Don', 'เกาะพีพีดอน', X'000000000101000000105839B4C8F61E4052B81E85EBB15840'),
	(8, 'Ko Rang', 'เกาะรัง', X'0000000001010000009A9999999999274054E3A59BC4985940'),
	(9, 'Ko Samet', 'เกาะเสม็ด', X'000000000101000000E17A14AE47212940CDCCCCCCCC5C5940'),
	(10, 'Ko Samui', 'เกาะสมุย', X'00000000010100000000000000000023400000000000005940'),
	(11, 'Ko Tao', 'เกาะเต่า', X'000000000101000000AE47E17A142E24401283C0CAA1F55840'),
	(12, 'Ko Tapu', 'เกาะตะปู', NULL),
	(13, 'Ko Tarutao', 'เกาะตะรุเตา', NULL),
	(14, 'Ko Phuket', 'เกาะภูเก็ต', NULL),
	(15, 'Similan Islands', 'เกาะสิมิลัน', X'000000000101000000CDCCCCCCCC4C21409A99999999695840'),
	(16, 'Ko Sichang', 'เกาะสีชัง', X'0000000001010000000AD7A3703D4A2A4079E9263108345940'),
	(17, 'Ko Kham Yai', 'เกาะขามใหญ่', X'0000000001010000002FDD240681552A40DBF97E6ABC345940'),
	(18, 'Ko Kham Noi', 'เกาะขามน้อย', X'00000000010100000045D8F0F44A592A40371AC05B20355940'),
	(19, 'Ko Ram Dok Mai', 'เกาะร้ามดอกไม้', X'00000000010100000058A835CD3B4E2A40AF94658863355940'),
	(20, 'Ko Khang Khao', 'เกาะค้างคาว', X'000000000101000000EE7C3F355E3A2A40C1CAA145B6335940'),
	(21, 'Ko Prong', 'เกาะปรง', X'00000000010100000031992A1895542A405A643BDF4F355940'),
	(22, 'Ko Yai Thao', '', X'0000000001010000001EA7E8482E3F2A40C7293A92CB335940'),
	(23, 'Ko Thaai Tamuen', '', X'0000000001010000000F9C33A2B4372A401D38674469335940'),
	(24, 'Ko Nok', 'เกาะนอก', X'000000000101000000508D976E12032A401D38674469335940'),
	(25, 'Ko Chun', '', X'0000000001010000005C8FC2F528DC2940A4703D0AD7335940'),
	(26, 'Ko Lan', 'เกาะล้าน', X'000000000101000000D7A3703D0AD7294052B81E85EB315940'),
	(27, 'Ko Sak', 'เกาะสาก', X'0000000001010000007B14AE47E1DA2940A69BC420B0325940'),
	(28, 'Ko Krok', 'เกาะครก', X'000000000101000000C217265305E32940C898BB9690335940'),
	(29, 'Ko Luam', 'เกาะเหลื่อม', X'0000000001010000005F984C158CEA2940C4B12E6EA3295940'),
	(30, 'Ko Luam Noi', 'เกาะเหลื่อมน้อย', X'0000000001010000000BB5A679C7E929404BEA0434112A5940'),
	(31, 'Ko Phai', 'เกาะไผ่', X'00000000010100000091ED7C3F35DE294033333333332B5940'),
	(32, 'Ko Hu Chang', 'เกาะหูช้าง', X'0000000001010000001904560E2DD229408FC2F5285C2B5940'),
	(33, 'Ko Klung Badan', 'เกาะกรุงบาดาล', X'0000000001010000003CBD529621CE29402497FF907E2B5940'),
	(34, 'Ko Man Wichai', 'เกาะมารวิชัย', X'00000000010100000038F8C264AAC02940A5BDC117262B5940'),
	(35, 'Ko Klet Kaeo', 'เกาะเกล็ดแก้ว', X'0000000001010000001283C0CAA1852940AE47E17A14365940'),
	(36, 'Ko Khram', 'เกาะคราม', X'0000000001010000006666666666662940C3F5285C8F325940'),
	(37, 'Ko Khram Noi', 'เกาะครามน้อย', X'0000000001010000004E62105839742940508D976E12335940'),
	(38, 'Ko I Ra', 'เกาะอีร้า', X'000000000101000000273108AC1C5A2940E9263108AC345940'),
	(39, 'Ko Tao Mo', 'เกาะเตาหม้อ', NULL),
	(40, 'Ko Phra', 'เกาะพระ', NULL),
	(41, 'Ko Phra Noi', 'เกาะพระน้อย', NULL),
	(42, 'Ko I Lao', 'เกาะอีเลา', NULL),
	(43, 'Ko Yo', 'เกาะยอ', NULL),
	(44, 'Ko Mu', 'เกาะหมู', NULL),
	(45, 'Ko Maeo', 'เกาะแมว', NULL),
	(46, 'Ko Nang Ram', 'เกาะนางรำ', NULL),
	(47, 'Ko Chorakhe', 'เกาะจระเข้', NULL),
	(48, 'Ko Samae San', 'เกาะแสมสาร', X'0000000001010000004C37894160252940CDCCCCCCCC3C5940'),
	(49, 'Ko Kham', 'เกาะขาม', X'0000000001010000004C37894160252940B29DEFA7C63B5940'),
	(50, 'Ko Raet', 'เกาะแรด', X'000000000101000000EC51B81E852B294004560E2DB23D5940'),
	(51, 'Ko Chang Kluea', '', X'0000000001010000009A99999999192940A01A2FDD243E5940'),
	(52, 'Ko Chuang', 'เกาะช่วง', X'0000000001010000000AD7A3703D0A29403D0AD7A3703D5940'),
	(53, 'Ko Chan', 'เกาะจาน', X'000000000101000000986E1283C00A2940AE47E17A143E5940'),
	(54, 'Ko Rong Nang', 'เกาะโรงหนัง', X'00000000010100000052EDD3F19811294036AB3E575B3D5940'),
	(55, 'Ko Rong Khon', 'เกาะโรงโขน', X'0000000001010000007138F3AB391029404C546F0D6C3D5940'),
	(56, 'Ko Saket', '', NULL),
	(57, 'Ko Samet', 'เกาะเสม็ด', X'000000000101000000E17A14AE47212940CDCCCCCCCC5C5940'),
	(58, 'Ko Platin', '', X'000000000101000000C0EC9E3C2C342940C5FEB27BF2605940'),
	(59, 'Ko Kraui', '', X'00000000010100000072F90FE9B72F2940FFB27BF2B0605940'),
	(60, 'Ko Kudi', '', X'000000000101000000D044D8F0F42A2940226C787AA5605940'),
	(61, 'Ko Khangkhao', '', X'000000000101000000F085C954C128294038F8C264AA605940'),
	(62, 'Ko Kham', 'เกาะขาม', X'000000000101000000AAF1D24D623029407F6ABC7493605940'),
	(63, 'Ko Yangklao', '', X'000000000101000000E86A2BF6971D2940D50968226C645940'),
	(64, 'Ko Thalu', '', X'0000000001010000003FC6DCB5841C29409CC420B072645940'),
	(65, 'Ko Khi Pla', 'เกาะขี้ปลา', NULL),
	(66, 'Ko Man Nai', 'เกาะมันใน', X'0000000001010000007E8CB96B0939294079E92631086C5940'),
	(67, 'Ko Man Khlang', 'เกาะมันกลาง', X'00000000010100000052B81E85EB3129403F355EBA496C5940'),
	(68, 'Ko Man Nok', 'เกาะมันนอก', X'000000000101000000EC2FBB270F2B2940F853E3A59B6C5940'),
	(69, 'Ko Chong Saba', 'เกาะช่องสะบ้า', NULL),
	(70, 'Ko Nom Sao', 'เกาะนมสาว', NULL),
	(71, 'Ko Chula', 'เกาะจุฬา', NULL),
	(72, 'Ko Nu', 'เกาะหนู', NULL),
	(73, 'Ko Nang Ram', 'เกาะนางรำ', NULL),
	(74, 'Ko Kwang', 'เกาะกวาง', NULL),
	(75, 'Ko Chik Nok', 'เกาะจิกนอก', NULL),
	(76, 'Ko Chik Klang', 'เกาะจิกกลาง', NULL),
	(77, 'Ko Chang', 'เกาะช้าง', X'0000000001010000001F85EB51B81E284085EB51B81E955940'),
	(78, 'Ko Ngam', 'เกาะง่าม', X'000000000101000000F4FDD478E9E627403108AC1C5A9C5940'),
	(79, 'Ko Phrao Nai', 'เกาะพร้าวใน', X'000000000101000000105839B4C8F6274054E3A59BC4985940'),
	(80, 'Ko Phrao Nok', 'เกาะพร้าวนอก', X'0000000001010000001904560E2DF227400C022B8716995940'),
	(81, 'Ko Klum', 'เกาะคลุ้ม', X'0000000001010000008716D9CEF7D327402DB29DEFA7965940'),
	(82, 'Mu Ko Lao Ya', 'หมู่เกาะเหลายา', X'000000000101000000FCA9F1D24DE22740355EBA490C9A5940'),
	(83, 'Ko Lao Ya Nai', 'เกาะเหลายาใน', X'000000000101000000FCA9F1D24DE22740355EBA490C9A5940'),
	(84, 'Ko Lao Ya Nok', 'เกาะเหลายานอก', X'000000000101000000E9263108ACDC2740FCA9F1D24D9A5940'),
	(85, 'Ko Wai', 'เกาะหวาย', X'000000000101000000CDCCCCCCCCCC274052B81E85EB995940'),
	(86, 'Mu Ko Mai Si', 'หมู่เกาะไม้ซี้', X'0000000001010000003108AC1C5AE4274066666666669E5940'),
	(87, 'Ko Mai Si Yai', 'เกาะไม้ซี้ใหญ่', X'0000000001010000003108AC1C5AE4274066666666669E5940'),
	(88, 'Ko Mai Si Lek', 'เกาะไม้ซี้เล็ก', X'000000000101000000295C8FC2F5E827402B8716D9CE9F5940'),
	(89, 'Ko Bai Dang', 'เกาะไปแดง', X'000000000101000000B29DEFA7C6CB2740BE9F1A2FDD9C5940'),
	(90, 'Ko Chan', 'เกาะจาน', X'000000000101000000A245B6F3FDD4274083C0CAA1459E5940'),
	(91, 'Ko Kut', 'เกาะกูด', X'00000000010100000052B81E85EB51274014AE47E17AA45940'),
	(92, 'Ko Mak', 'เกาะหมาก', X'00000000010100000017D9CEF753A327401F85EB51B89E5940'),
	(93, 'Ko Rayang Nok', 'เกาะระยั้งนอก', X'0000000001010000007F6ABC7493982740BE9F1A2FDD9C5940'),
	(94, 'Ko Kradat', 'เกาะกระดาด', X'000000000101000000AE47E17A14AE2740F6285C8FC2A15940'),
	(95, 'Mu Ko Rang', 'หมู่เกาะรัง', X'0000000001010000009A9999999999274054E3A59BC4985940'),
	(96, 'Ko Rang', 'เกาะรัง', X'0000000001010000009A9999999999274054E3A59BC4985940'),
	(97, 'Ko Tun', 'เกาะตูน', X'0000000001010000003F355EBA498C2740FED478E926995940'),
	(98, 'Ko Phing', 'เกาะพิง', NULL),
	(99, 'Ko Phang', 'เกาะพัง', NULL),
	(100, 'Ko Sai', 'เกาะทราย', X'0000000001010000009EEFA7C64BF728400000000000005940'),
	(101, 'Ko Sadao', 'เกาะสะเดา', X'0000000001010000003BDF4F8D97EE28400000000000005940'),
	(102, 'Ko Khi Nok', 'เกาะขี้นก', X'0000000001010000009EEFA7C64BF728400000000000005940'),
	(103, 'Ko Lak', 'เกาะหลัก', NULL),
	(104, 'Ko La', 'เกาะละ', NULL),
	(105, 'Ko Rom', 'เกาะร่ม', NULL),
	(106, 'Ko Raet', 'เกาะแรด', NULL),
	(107, 'Ko Chan', 'เกาะจาน', NULL),
	(108, 'Ko Thai Chan', 'เกาะท้ายจาน', NULL),
	(109, 'Ko Nom Sao', 'เกาะนมสาว', NULL),
	(110, 'Ko Kho Ram', 'เกาะโครำ', NULL),
	(111, 'Ko Rawang', 'เกาะระวาง', NULL),
	(112, 'Ko Rawing', 'เกาะระวิง', NULL),
	(113, 'Ko Sattakut', 'เกาะสัตกูต', X'00000000010100000066666666666628401B2FDD2406015940'),
	(114, 'Ko Hua Pin', 'เกาะหัวพิน', NULL),
	(115, 'Ko Thalu', 'เกาะทะลุ', NULL),
	(116, 'Ko Sing', 'เกาะสิงข์', NULL),
	(117, 'Ko Sang', 'เกาะสังข์', NULL),
	(118, 'Ko Ram Ra', 'เกาะร่ำรา', NULL),
	(119, 'Ko Lueam', 'เกาะเหลื่อม', NULL),
	(120, 'Ko Aen', 'เกาะแอ่น', NULL),
	(121, 'Ko Wiang', 'เกาะเวียง', NULL),
	(122, 'Ko Phra', 'เกาะพระ', NULL),
	(123, 'Ko Yo', 'เกาะยอ', NULL),
	(124, 'Ko Khi Nok', 'เกาะขี้นก', NULL),
	(125, 'Ko Si Kong', 'เกาะซีกง', NULL),
	(126, 'Ko Rang', 'เกาะรัง', NULL),
	(127, 'Ko Khai', 'เกาะไข่', NULL),
	(128, 'Ko Chorakhe', 'เกาะจรเข้', NULL),
	(129, 'Ko Ngam Yai', 'เกาะงามใหญ่', NULL),
	(130, 'Ko Ngam Noi', 'เกาะงามน้อย', NULL),
	(131, 'Ko Kalok', 'เกาะกะโหลก', NULL),
	(132, 'Ko Samet', 'เกาะเสม็ด', NULL),
	(133, 'Ko Sak', 'เกาะสาก', NULL),
	(134, 'Ko Maphrao', 'เกาะมะพร้าว', NULL),
	(135, 'Ko Matra', 'เกาะมาตรา', NULL),
	(136, 'Ko Lak Raet', 'เกาะหลักแรด', NULL),
	(137, 'Ko I Raet', 'เกาะอีแรด', NULL),
	(138, 'Ko Lawa', 'เกาะละวะ', NULL),
	(139, 'Ko Ka', 'เกาะกา', NULL),
	(140, 'Ko Thong Lang', 'เกาะทองหลาง', NULL),
	(141, 'Ko Wang Ka Chio', 'เกาะวังกะจิว', NULL),
	(142, 'Ko Kra', '', X'00000000010100000000000000008024400000000000D05840'),
	(143, 'Ko Kula', '', X'00000000010100000000000000008024400000000000D05840'),
	(144, 'Ko Samui', 'เกาะสมุย', X'00000000010100000000000000000023400000000000005940'),
	(145, 'Ko Pha Ngan', 'เกาะพะงัน', NULL),
	(146, 'Ko Nang Yuan', '', X'000000000101000000AE47E17A142E24401283C0CAA1F55840'),
	(147, 'Ko Tao', 'เกาะเต่า', X'000000000101000000AE47E17A142E24401283C0CAA1F55840'),
	(148, 'Ko Similan', 'เกาะสิมิลัน', X'000000000101000000CDCCCCCCCC4C21409A99999999695840'),
	(149, 'Ko Bangu', '', X'000000000101000000022B8716D94E21409A99999999695840'),
	(150, 'Ko Payu', '', X'000000000101000000AE47E17A142E214046B6F3FDD4685840'),
	(151, 'Ko Miang', '', X'000000000101000000333333333333214054E3A59BC4685840'),
	(152, 'Ko Payang', '', X'000000000101000000E5D022DBF9FE2040D34D621058695840'),
	(153, 'Ko Huyang', '', X'0000000001010000006891ED7C3FF520409A99999999695840'),
	(154, 'Ko Nom Sao', 'เกาะนมสาว', NULL),
	(155, 'Ko Tapu', 'เกาะตะปู', NULL),
	(156, 'Ko Pu', 'เกาะปู', NULL),
	(157, 'Ko Siboya', 'เกาะปู', NULL),
	(158, 'Ko Jum', 'เกาะปู', NULL),
	(159, 'Ko Kai, Krabi', 'เกาะไก่', NULL),
	(160, 'Ko Pli', 'เกาะไก่', NULL),
	(161, 'Ko Lek', 'เกาะไก่', NULL),
	(162, 'Ko Lapu Le', 'เกาะไก่', NULL),
	(163, 'Ko A Dang', 'เกาะไก่', NULL),
	(164, 'Ko Poda', '', NULL),
	(165, 'Ko Lanta Yai', 'เกาะลันตาใหญ่', NULL),
	(166, 'Ko Lanta Noi', 'เกาะลันตาน้อย', NULL),
	(167, 'Ko Bu Bu', '', NULL),
	(168, 'Ko Po', '', NULL),
	(169, 'Ko Klang (Krabi)', '', NULL),
	(170, 'Ko Ma', '', NULL),
	(171, 'Ko Ha', '', NULL),
	(172, 'Phuket', 'ภูเก็ต', NULL),
	(173, 'Ko Lon', 'เกาะโหลน', NULL),
	(174, 'Ko Mai Thon', 'เกาะไม้ท่อน', NULL),
	(175, 'Ko Maphrao', 'เกาะมะพร้าว', NULL),
	(176, 'Ko Ngam (Phuket)', 'เกาะงำ', NULL),
	(177, 'Ko Naka Yai', 'เกาะนาคาใหญ่', NULL),
	(178, 'Ko Naka Noi', 'เกาะนาคาน้อย', NULL),
	(179, 'Ko Raet', 'เกาะแรด', NULL),
	(180, 'Ko Siray', 'เกาะสิเหร่', NULL),
	(181, 'Ko Rang Yai', 'เกาะรังใหญ่', NULL),
	(182, 'Ko Rang Noi', 'เกาะรังน้อย', NULL),
	(183, 'Ko Tapao', 'เกาะตะเภา', NULL),
	(184, 'Ko Kaeo Yai', 'เกาะแก้วใหญ่', NULL),
	(185, 'Ko Kaeo', 'เกาะแก้ว', NULL),
	(186, 'Ko Bon', 'เกาะบอน', NULL),
	(187, 'Ko Hae', 'เกาะเฮ', NULL),
	(188, 'Ko Aeo', 'เกาะแอว', NULL),
	(189, 'Ko Nom Sao', 'เกาะนมสาว', NULL),
	(190, 'Ko Tapu', 'เกาะตะปู', NULL),
	(191, 'Ko Pu', 'เกาะปู', NULL),
	(192, 'Ko Siboya', 'เกาะปู', NULL),
	(193, 'Ko Jum', 'เกาะปู', NULL),
	(194, 'Ko Kai, Krabi', 'เกาะไก่', NULL),
	(195, 'Ko Pli', 'เกาะไก่', NULL),
	(196, 'Ko Lek', 'เกาะไก่', NULL),
	(197, 'Ko Lapu Le', 'เกาะไก่', NULL),
	(198, 'Ko A Dang', 'เกาะไก่', NULL),
	(199, 'Ko Poda', '', NULL),
	(200, 'Ko Lanta Yai', 'เกาะลันตาใหญ่', NULL),
	(201, 'Ko Lanta Noi', 'เกาะลันตาน้อย', NULL),
	(202, 'Ko Bu Bu', '', NULL),
	(203, 'Ko Po', '', NULL),
	(204, 'Ko Klang (Krabi)', '', NULL),
	(205, 'Ko Ma', '', NULL),
	(206, 'Ko Ha', '', NULL),
	(207, 'Phuket', 'ภูเก็ต', NULL),
	(208, 'Ko Lon', 'เกาะโหลน', NULL),
	(209, 'Ko Mai Thon', 'เกาะไม้ท่อน', NULL),
	(210, 'Ko Maphrao', 'เกาะมะพร้าว', NULL),
	(211, 'Ko Ngam (Phuket)', 'เกาะงำ', NULL),
	(212, 'Ko Naka Yai', 'เกาะนาคาใหญ่', NULL),
	(213, 'Ko Naka Noi', 'เกาะนาคาน้อย', NULL),
	(214, 'Ko Raet', 'เกาะแรด', NULL),
	(215, 'Ko Siray', 'เกาะสิเหร่', NULL),
	(216, 'Ko Rang Yai', 'เกาะรังใหญ่', NULL),
	(217, 'Ko Rang Noi', 'เกาะรังน้อย', NULL),
	(218, 'Ko Tapao', 'เกาะตะเภา', NULL),
	(219, 'Ko Kaeo Yai', 'เกาะแก้วใหญ่', NULL),
	(220, 'Ko Kaeo', 'เกาะแก้ว', NULL),
	(221, 'Ko Bon', 'เกาะบอน', NULL),
	(222, 'Ko Hae', 'เกาะเฮ', NULL),
	(223, 'Ko Aeo', 'เกาะแอว', NULL),
	(224, 'Ko Ngai', 'เกาะไหง', X'0000000001010000009015FC36C4A81D4089618731E9C95840'),
	(225, 'Ko Kradan', 'เกาะกระดาน', X'000000000101000000F030ED9BFB3B1D40AF43352559D05840'),
	(226, 'Ko Muk', 'เกาะมุก', X'00000000010100000042D13C80457E1D4077DCF0BBE9D25840'),
	(227, 'Ko Racha Yai', 'เกาะราชาใหญ่', X'0000000001010000005E126745D4641E40A872DA5372975840'),
	(228, 'Ko Lamphu', 'เกาะลำพู', X'00000000010100000085EB51B81E4522402506819543D35840'),
	(229, 'Ko Lamphu Rai', 'เกาะลำพูราย', X'00000000010100000066666666666628405A643BDF4FA55940');
