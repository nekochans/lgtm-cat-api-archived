package test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-sql-driver/mysql"
)

type Seeder struct {
	Db      *sql.DB
	DirPath string
}

func (s *Seeder) Execute() error {
	files, err := os.ReadDir(s.DirPath)
	if err != nil {
		return fmt.Errorf("failed to read test data dir: %s", err)
	}

	tx, err := s.Db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %s", err)
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext != ".csv" {
			continue
		}

		table := file.Name()[:len(file.Name())-len(ext)]
		csvFilePath := filepath.Join(s.DirPath, file.Name())

		if _, err := s.loadDataFromCsv(tx, table, csvFilePath); err != nil {
			if err := tx.Rollback(); err != nil {
				return fmt.Errorf("failed to transaction rollback: %s", err)
			}

			return fmt.Errorf("failed to load data from csv: %s", err)
		}
	}

	return tx.Commit()
}

func (s *Seeder) TruncateAllTable() error {
	tx, err := s.Db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %s", err)
	}

	if _, err := tx.Exec("TRUNCATE lgtm_images"); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("failed to transaction rollback: %s", err)
		}

		return fmt.Errorf("failed to exec sql: %s", err)
	}

	return tx.Commit()
}

func (s *Seeder) loadDataFromCsv(tx *sql.Tx, table, filePath string) (sql.Result, error) {
	query := `
		LOAD DATA
			LOCAL INFILE '%s'
		INTO TABLE %s
		FIELDS
			TERMINATED BY ','
		LINES
			TERMINATED BY '\n'
			IGNORE 1 LINES
	`

	mysql.RegisterLocalFile(filePath)

	return tx.Exec(fmt.Sprintf(query, filePath, table))
}
