package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strconv"

	_ "github.com/lib/pq" //
)

// Storage implements a PostgreSQL storage backend for colly
type Storage struct {
	URI                string
	VisitedTable       string
	CookiesTable       string
	MaxOpenConnections uint8
	db                 *sql.DB
}

// Init initializes the PostgreSQL storage
func (s *Storage) Init() error {

	var err error

	if s.db, err = sql.Open("postgres", s.URI); err != nil {
		log.Fatal(err)
	}

	if err = s.db.Ping(); err != nil {
		log.Fatal(err)
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (request_id text not null);", s.VisitedTable)

	if _, err = s.db.Exec(query); err != nil {
		log.Fatal(err)
	}

	query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (host text not null, cookies text not null);", s.CookiesTable)

	if _, err = s.db.Exec(query); err != nil {
		log.Fatal(err)
	}

	s.db.SetMaxOpenConns(int(s.MaxOpenConnections))

	return nil

}

// Visited implements colly/storage.Visited()
func (s *Storage) Visited(requestID uint64) error {

	var err error

	query := fmt.Sprintf(`INSERT INTO %s (request_id) VALUES($1);`, s.VisitedTable)

	_, err = s.db.Exec(query, strconv.FormatUint(requestID, 10))

	return err

}

// IsVisited implements colly/storage.IsVisited()
func (s *Storage) IsVisited(requestID uint64) (bool, error) {

	var isVisited bool

	query := fmt.Sprintf(`SELECT EXISTS(SELECT request_id FROM %s WHERE request_id = $1)`, s.VisitedTable)

	err := s.db.QueryRow(query, strconv.FormatUint(requestID, 10)).Scan(&isVisited)

	return isVisited, err

}

// Cookies implements colly/storage.Cookies()
func (s *Storage) Cookies(u *url.URL) string {

	var cookies string

	query := fmt.Sprintf(`SELECT cookies FROM %s WHERE host = $1;`, s.CookiesTable)

	s.db.QueryRow(query, u.Host).Scan(&cookies)

	return cookies

}

// SetCookies implements colly/storage.SetCookies()
func (s *Storage) SetCookies(u *url.URL, cookies string) {

	query := fmt.Sprintf(`INSERT INTO %s (host, cookies) VALUES($1, $2);`, s.CookiesTable)

	s.db.Exec(query, u.Host, cookies)

}
