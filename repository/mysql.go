package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
}

func NewMySQLDB(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitDB(db *sql.DB) error {
	// Начало транзакции
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Таблица активностей
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS activity (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		sid VARCHAR(255) NOT NULL,
		score INT NOT NULL
	);
	`
	_, err = tx.Exec(createTableSQL)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create table: %v", err)
	}
	// Таблица активностей
	createTableSQL = `
	CREATE TABLE IF NOT EXISTS players (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		score INT NOT NULL
	);
	`
	_, err = tx.Exec(createTableSQL)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create table: %v", err)
	}

	// Таблица активностей
	createTableSQL = `
	CREATE TABLE IF NOT EXISTS rel_activity_players (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		activity_id INT UNSIGNED NOT NULL,
		player_id INT UNSIGNED NOT NULL,
		score INT NOT NULL,
		time INT NOT NULL,
		FOREIGN KEY (activity_id) REFERENCES activity(id),
		FOREIGN KEY (player_id) REFERENCES players(id)
	);
	`
	_, err = tx.Exec(createTableSQL)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create table: %v", err)
	}

	// Завершение транзакции
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	fmt.Println("Database and table created successfully")
	return nil
}
