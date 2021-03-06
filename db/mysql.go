package db

import (
	"database/sql"
	"fmt"
	"go-graphql/config"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	conn *sql.DB
}

var DB = Database{

}

type MySQLConfig struct {
	Username   string
	Password   string
	Host       string
	Port       int
	UnixSocket string
}

type RowScanner interface {
	Scan(dest ...interface{}) error
}

func newDatabase() (*Database, error) {

	conn, err := sql.Open("mysql", config.MysqlConnectURL)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, err
	}

	db := &Database{
		conn: conn,
	}

	DB = *db

	return db, err
}

func (db *Database) Close() {
	DB.conn.Close()
}

func (db *Database) Query(query string, args interface{}) (*sql.Rows, error) {
	return DB.conn.Query(query, args)
}
func (db *Database) QueryRow(query string, args ...interface{}) (*sql.Row) {

	return DB.conn.QueryRow(query, args...)
}

func (db *Database) Prepare(query string) (*sql.Stmt, error) {
	return DB.conn.Prepare(query)
}

func InitDatabase() (*Database, error) {
	DB, err := newDatabase()
	if err != nil {
		return nil, err
	}

	return DB, err
}

func (db *Database) Insert(query string, args ...interface{}) (int64, error) {

	stmt, _ := DB.Prepare(query)

	r, err := stmt.Exec(args...)
	if err != nil {
		return 0, fmt.Errorf("mysql: could not execute statement: %v", err)
	}
	rowsAffected, err := r.RowsAffected()

	if err != nil {
		return 0, fmt.Errorf("mysql: could not get rows affected: %v", err)
	} else if rowsAffected != 1 {
		return 0, fmt.Errorf("mysql: expected 1 row affected, got %d", rowsAffected)
	}

	lastInsertID, err := r.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
	}
	return lastInsertID, nil

}

func (db *Database) Update(query string, args ...interface{}) (int64, error) {

	stmt, _ := DB.conn.Prepare(query)

	r, err := stmt.Exec(args...)
	if err != nil {
		return 0, fmt.Errorf("mysql: could not execute statement: %v", err)
	}
	rowsAffected, err := r.RowsAffected()

	if err != nil {
		return 0, fmt.Errorf("mysql: could not get rows affected: %v", err)
	} else if rowsAffected != 1 {
		return 0, fmt.Errorf("mysql: expected 1 row affected, got %d", rowsAffected)
	}

	lastInsertID, err := r.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
	}

	return lastInsertID, nil

}

func (db *Database) Count(query string, args ...interface{}) (int, error) {

	var count int
	row := DB.conn.QueryRow(query, args...)

	err := row.Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (db *Database) Get(table string, id int64) (*sql.Row, error) {

	query := "SELECT * FROM " +
		table +
		" WHERE id = ?"
	stmt, _ := DB.conn.Prepare(query)
	row := stmt.QueryRow(id)
	return row, nil

}
func (db *Database) FindOne(query string, args ...interface{}) (*sql.Row, error) {

	stmt, err := DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(args...)
	return row, nil
}

func (db *Database) Delete(query string, args ...interface{}) (int64, error) {

	stmt, _ := DB.conn.Prepare(query)

	r, err := stmt.Exec(args...)
	if err != nil {
		return 0, fmt.Errorf("could not execute statement: %v", err)
	}
	rowsAffected, err := r.RowsAffected()

	if err != nil {
		return 0, fmt.Errorf("could not get rows affected: %v", err)
	} else if rowsAffected != 1 {
		return 0, fmt.Errorf("expected 1 row affected, got %d", rowsAffected)
	}

	lastInsertID, err := r.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("could not get last insert ID: %v", err)
	}

	return lastInsertID, nil

}

func (db *Database) List(query string, args ...interface{}) (*sql.Rows, error) {

	stmt, err := DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(args...)

	return rows, err

}
