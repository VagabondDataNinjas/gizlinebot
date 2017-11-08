package storage

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"time"

	"github.com/VagabondDataNinjas/gizlinebot/domain"
	"github.com/go-sql-driver/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type Sql struct {
	Db *sql.DB
}

func NewSql(conDsn string) (s *Sql, err error) {
	db, err := sql.Open("mysql", conDsn)
	if err != nil {
		return s, err
	}

	return &Sql{
		Db: db,
	}, nil
}

func (s *Sql) Close() error {
	return s.Db.Close()
}

func (s *Sql) AddRawLineEvent(eventType, rawevent string) error {
	stmt, err := s.Db.Prepare("INSERT INTO linebot_raw_events(eventtype, rawevent, timestamp) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(eventType, rawevent, int32(time.Now().UTC().Unix()))
	if err != nil {
		return err
	}

	return nil
}

// AddUpdateUserProfile adds a user profile
// if the user already exists in the table this method does nothing
func (s *Sql) AddUpdateUserProfile(userID, displayName string) error {
	stmt, err := s.Db.Prepare(`INSERT INTO
		user_profiles(userId, displayName, timestamp, active) VALUES(?, ?, ?, 1)
		ON DUPLICATE KEY UPDATE displayName = ?, timestamp = ?, active = 1
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	ts := int32(time.Now().UTC().Unix())
	_, err = stmt.Exec(userID, displayName, ts, displayName, ts)

	if err != nil {
		if mysqlErr := err.(*mysql.MySQLError); mysqlErr.Number == 1062 {
			// ignore duplicate entry errors for profiles
			return nil
		}
		return err
	}

	return nil
}

func (s *Sql) UpdateUserProfile(p domain.UserProfile) error {
	stmt, err := s.Db.Prepare(`UPDATE user_profiles
		SET displayName = ?, timestamp = ?, active = ?, bot_survey_inited = ? WHERE userId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.DisplayName, p.Timestamp, p.Active, p.SurveyStarted, p.UserId)

	if err != nil {
		return err
	}

	return nil
}

func (s *Sql) ToggleUserSurvey(userId string, started bool) error {
	stmt, err := s.Db.Prepare("UPDATE user_profiles SET bot_survey_inited = ? WHERE userId = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	startedInt := 0
	if started {
		startedInt = 1
	}
	_, err = stmt.Exec(startedInt, userId)

	if err != nil {
		return err
	}

	return nil
}

func (s *Sql) UsersSurveyNotStarted(delaySecs int64) (userIds []string, err error) {
	var (
		userId string
	)

	tsCompare := time.Now().UTC().Unix() - delaySecs
	rows, err := s.Db.Query(`SELECT p.userId FROM user_profiles p
		WHERE p.active = 1 AND p.bot_survey_inited = 0 AND p.timestamp < ?`, tsCompare)
	if err != nil {
		return userIds, err
	}
	defer rows.Close()

	userIds = make([]string, 0)
	for rows.Next() {
		err := rows.Scan(&userId)
		if err != nil {
			return userIds, err
		}
		userIds = append(userIds, userId)
	}

	err = rows.Err()
	if err != nil {
		return userIds, err
	}

	return userIds, nil
}

func (s *Sql) GetUserProfile(userId string) (domain.UserProfile, error) {
	p := domain.UserProfile{
		UserId: userId,
	}
	err := s.Db.QueryRow(`SELECT displayName, timestamp, active, IF(bot_survey_inited = 1, TRUE, FALSE) AS SurveyStarted
		FROM user_profiles where userId = ?`, userId).Scan(&p.DisplayName, &p.Timestamp, &p.Active, &p.SurveyStarted)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, nil
		}
		return p, err
	}

	return p, nil
}

func (s *Sql) UserHasAnswers(userId string) (bool, error) {
	profile, err := s.GetUserProfile(userId)
	if err != nil {
		return false, err
	}

	var hasAnswers int
	// select answers that have been added after the last
	// block/unblock action
	err = s.Db.QueryRow(`SELECT count(id) FROM answers
		WHERE userId = ? AND timestamp > ?`, userId, profile.Timestamp).Scan(&hasAnswers)
	if err != nil {
		return false, err
	}

	if hasAnswers > 0 {
		return true, nil
	}
	return false, nil
}

func (s *Sql) UserGetLastAnswer(uid string) (domain.Answer, error) {
	a := domain.Answer{}
	var ts int64
	profile, err := s.GetUserProfile(uid)
	if err != nil {
		return a, err
	}

	err = s.Db.QueryRow(`SELECT id, userId, questionId, answer, timestamp FROM answers
		WHERE userId = ? AND answer != "" AND timestamp > ?
		ORDER BY timestamp DESC
		LIMIT 0,1
		`, uid, profile.Timestamp).Scan(&a.Id, &a.UserId, &a.QuestionId, &a.Answer, &ts)
	if err != nil {
		return a, err
	}

	a.Timestamp = time.Unix(ts, 0)
	return a, nil
}

func (s *Sql) GetQuestions() (qs *domain.Questions, err error) {
	var (
		id           string
		questionText string
		weight       int
		channel      string
	)
	rows, err := s.Db.Query(`SELECT id, question, weight, channel FROM questions ORDER BY weight ASC`)
	if err != nil {
		return qs, err
	}
	defer rows.Close()

	qs = domain.NewQuestions()
	for rows.Next() {
		err := rows.Scan(&id, &questionText, &weight, &channel)
		if err != nil {
			return qs, err
		}
		err = qs.Add(id, questionText, weight, channel)
		if err != nil {
			return qs, err
		}
	}

	err = rows.Err()
	if err != nil {
		return qs, err
	}

	return qs, nil
}

type WelcomeMsgTplVars struct {
	UserId   string
	Hostname string
}

func (s *Sql) GetWelcomeMsgs(tplVars *WelcomeMsgTplVars) (msgs []string, err error) {
	var (
		msgRaw string
	)
	rows, err := s.Db.Query(`SELECT msg FROM welcome_msgs WHERE channel IN ("line", "both") ORDER BY weight ASC`)
	if err != nil {
		return msgs, err
	}
	defer rows.Close()

	msgs = make([]string, 0)
	for rows.Next() {
		err := rows.Scan(&msgRaw)
		if err != nil {
			return msgs, err
		}
		msg, err := s.applyWelcomeTpl(msgRaw, tplVars)
		if err != nil {
			return msgs, err
		}
		msgs = append(msgs, msg)
	}

	err = rows.Err()
	if err != nil {
		return msgs, err
	}

	return msgs, nil
}

type UserAnswerData struct {
	// @TODO embed domain.Answer
	// domain.Answer
	Id         uint
	UserId     string
	QuestionId string
	Answer     string
	Channel    string
	Timestamp  int
}

type UserGpsAnswerData struct {
	Id        uint
	UserId    string
	Address   string
	Lat       float64
	Lon       float64
	Timestamp int
	Channel   string
}

func (s *Sql) GetGpsAnswerData() (answerGpsData []UserGpsAnswerData, err error) {
	rows, err := s.Db.Query(`SELECT p.id, p.userId, IFNULL(a.address, ""), IFNULL(a.lat, 0.0), IFNULL(a.lon, 0.0), IFNULL(a.channel, ""), IFNULL(a.timestamp, 0) FROM user_profiles p
		LEFT JOIN answers_gps a ON a.userId = p.userId
		ORDER BY a.timestamp ASC
		`)
	if err != nil {
		return answerGpsData, err
	}
	defer rows.Close()

	answerGpsData = make([]UserGpsAnswerData, 0)
	for rows.Next() {
		a := UserGpsAnswerData{}
		err := rows.Scan(&a.Id, &a.UserId, &a.Address, &a.Lat, &a.Lon, &a.Channel, &a.Timestamp)
		if err != nil {
			return answerGpsData, err
		}
		answerGpsData = append(answerGpsData, a)
	}

	err = rows.Err()
	if err != nil {
		return answerGpsData, err
	}

	return answerGpsData, nil

}
func (s *Sql) GetUserAnswerData() (answerData []UserAnswerData, err error) {
	var (
		userId     string
		questionId string
		answer     string
		channel    string
		answerTime int
	)
	rows, err := s.Db.Query(`SELECT p.userId, IFNULL(a.questionId, ""), IFNULL(a.answer, ""), IFNULL(a.channel, ""), IFNULL(a.timestamp, 0) as answerTime FROM user_profiles p
		LEFT JOIN answers a ON a.userId = p.userId
		ORDER BY a.timestamp ASC
		`)
	if err != nil {
		return answerData, err
	}
	defer rows.Close()

	answerData = make([]UserAnswerData, 0)
	for rows.Next() {
		err := rows.Scan(&userId, &questionId, &answer, &channel, &answerTime)
		if err != nil {
			return answerData, err
		}
		answerData = append(answerData, UserAnswerData{
			UserId:     userId,
			QuestionId: questionId,
			Answer:     answer,
			Channel:    channel,
			Timestamp:  answerTime,
		})
	}

	err = rows.Err()
	if err != nil {
		return answerData, err
	}

	return answerData, nil
}

func (s *Sql) GetLocations() (locs []domain.LocationThai, err error) {
	rows, err := s.Db.Query(`SELECT id, name, thainame, IFNULL(latitude, 0), IFNULL(longitude, 0) FROM locations`)
	if err != nil {
		return locs, err
	}
	defer rows.Close()

	locs = make([]domain.LocationThai, 0)
	for rows.Next() {
		loc := domain.LocationThai{}
		err := rows.Scan(&loc.Id, &loc.Name, &loc.NameThai, &loc.Latitude, &loc.Longitude)
		if err != nil {
			return locs, err
		}
		locs = append(locs, loc)
	}

	err = rows.Err()
	if err != nil {
		return locs, err
	}

	return locs, nil
}

func (s *Sql) applyWelcomeTpl(msg string, tplVars *WelcomeMsgTplVars) (string, error) {
	tmpl, err := template.New("welcomeMsg").Parse(msg)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, tplVars)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *Sql) UserAddAnswer(answer domain.Answer) error {
	stmt, err := s.Db.Prepare("INSERT INTO answers(userId, questionId, answer, channel, timestamp) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(answer.UserId, answer.QuestionId, answer.Answer, answer.Channel, int32(time.Now().UTC().Unix()))

	if err != nil {
		return err
	}

	return nil
}

type WebSurveyBtnConfig struct {
	Title     string
	Text      string
	Label     string
	ImageName string
}

func (s *Sql) GetWebSurveyBtnConfig() (WebSurveyBtnConfig, error) {
	cfg := WebSurveyBtnConfig{}
	title, err := s.GetConfigVal("web_survey_btn_title")
	if err != nil {
		return cfg, err
	}
	text, err := s.GetConfigVal("web_survey_btn_text")
	if err != nil {
		return cfg, err
	}
	label, err := s.GetConfigVal("web_survey_btn_label")
	if err != nil {
		return cfg, err
	}
	filename, err := s.GetConfigVal("web_survey_btn_imgname")
	if err != nil {
		return cfg, err
	}

	return WebSurveyBtnConfig{
		Title:     title,
		Text:      text,
		Label:     label,
		ImageName: filename,
	}, nil
}

func (s *Sql) GetConfigVal(key string) (string, error) {
	var val string
	err := s.Db.QueryRow(`
		SELECT value FROM config WHERE k = ?
	`, key).Scan(&val)
	if err != nil {
		return "", err
	}

	return val, nil
}

func (s *Sql) GetPriceTplMsg() (string, error) {
	return s.GetConfigVal("price_tpl")
}

func (s *Sql) GetUserNearbyPrices(userId string) (lp []domain.LocationPrice, err error) {
	locStr, err := s.userLastLocationAnswer(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return lp, nil
		}
		return lp, err
	}

	loc, err := s.findLocation(locStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return lp, nil
		}
		return lp, err
	}

	return s.getNearbyLocations(loc.Latitude, loc.Longitude, 300.0, 3)
}

func (s *Sql) findLocation(locName string) (domain.Location, error) {
	var loc = domain.Location{}
	err := s.Db.QueryRow(`SELECT id, thainame, latitude, longitude
		FROM locations
		WHERE (name = ? OR thainame = ?)
			AND latitude IS NOT NULL AND longitude IS NOT NULL
		LIMIT 0,1
		`, locName, locName).Scan(&loc.Id, &loc.Name, &loc.Latitude, &loc.Longitude)
	if err != nil {
		return loc, err
	}

	return loc, nil
}

func (s *Sql) userLastLocationAnswer(userId string) (string, error) {
	var answer string
	err := s.Db.QueryRow(`SELECT answer FROM answers
		WHERE questionid = 'island' AND userId = ?
			AND answer IS NOT NULL AND answer != ""
		ORDER BY timestamp DESC LIMIT 0,1
		`, userId).Scan(&answer)
	if err != nil {
		return "", err
	}

	return answer, nil
}

func (s *Sql) getNearbyLocations(lat, lon, radius float64, limit int) (locs []domain.LocationPrice, err error) {
	rows, err := s.Db.Query(`SELECT name,
		id, latitude, longitude, price
		FROM (
				SELECT l.id, l.name AS name,
					l.latitude, l.longitude, AVG(pp.price) AS price,
					p.radius,
					p.distance_unit
										* DEGREES(ACOS(COS(RADIANS(p.latpoint))
										* COS(RADIANS(l.latitude))
										* COS(RADIANS(p.longpoint - l.longitude))
										+ SIN(RADIANS(p.latpoint))
										* SIN(RADIANS(l.latitude)))
					) AS distance
				FROM locations AS l
				INNER JOIN pricepoints pp
					ON pp.location_id = l.id
				JOIN (   /* these are the query parameters */
						SELECT ? AS latpoint, ? AS longpoint,
										? AS radius, 111.045 AS distance_unit
					) AS p ON 1=1
				WHERE l.latitude
					BETWEEN p.latpoint  - (p.radius / p.distance_unit)
							AND p.latpoint  + (p.radius / p.distance_unit)
				AND l.longitude
					BETWEEN p.longpoint - (p.radius / (p.distance_unit * COS(RADIANS(p.latpoint))))
							AND p.longpoint + (p.radius / (p.distance_unit * COS(RADIANS(p.latpoint))))
				GROUP BY l.id
				) AS d
		WHERE distance <= radius
		ORDER BY distance ASC
		/* Ignore the first result which is the island of the lat/lon args */
		LIMIT 1, ?;
		`, lat, lon, radius, limit)
	if err != nil {
		return locs, err
	}
	defer rows.Close()

	locs = make([]domain.LocationPrice, 0)
	for rows.Next() {
		var lp = domain.LocationPrice{}
		err := rows.Scan(&lp.Name, &lp.Id, &lp.Latitude, &lp.Longitude, &lp.Price)
		if err != nil {
			return locs, err
		}
		locs = append(locs, lp)
	}

	err = rows.Err()
	if err != nil {
		return locs, err
	}

	return locs, nil
}

func (s *Sql) WipeUser(userId string) error {
	for _, table := range []string{"user_profiles", "answers", "answers_gps"} {
		err := s.deleteFromTableUserId(table, userId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Sql) deleteFromTableUserId(table string, userId string) error {
	// @TODO find out how to use dynamic table name in prepared query
	q := fmt.Sprintf("DELETE FROM %s WHERE userId = ?", table)
	stmt, err := s.Db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sql) UserAddGpsAnswer(answer domain.AnswerGps) error {
	stmt, err := s.Db.Prepare("INSERT INTO answers_gps(userId, lat, lon, address, channel, timestamp) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(answer.UserId, answer.Lat, answer.Lon, answer.Address, answer.Channel, int32(time.Now().UTC().Unix()))

	if err != nil {
		return err
	}

	return nil
}
