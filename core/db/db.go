package dbcore

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// New creates new database connection to the database server
func New(dbPsn string, cfg *gorm.Config) (*gorm.DB, error) {
	db := new(gorm.DB)

	db, err := gorm.Open(postgres.Open(dbPsn), cfg)
	if err != nil {
		return nil, err
	}
	return db, nil
}
