package storage

import (
	"database/sql"
	"time"

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
	_, err = stmt.Exec(eventType, rawevent, int32(time.Now().UTC().Unix()))
	if err != nil {
		return err
	}

	return nil
}

// AddUserProfile adds a user profile
// if the user already exists in the table this method does nothing
func (s *Sql) AddUserProfile(userID, displayName string) error {
	stmt, err := s.Db.Prepare("INSERT INTO user_profiles(userId, displayName, timestamp) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userID, displayName, int32(time.Now().UTC().Unix()))

	if err != nil {
		if mysqlErr := err.(*mysql.MySQLError); mysqlErr.Number == 1062 {
			// ignore duplicate entry errors for profiles
			return nil
		}
		return err
	}

	return nil
}
