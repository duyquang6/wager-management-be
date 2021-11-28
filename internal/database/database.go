package database

import (
	"context"
	"errors"

	"github.com/duyquang6/wager-management-be/pkg/logging"
	"gorm.io/gorm"
)

var _ DBFactory = (*DB)(nil)
type DB struct {
	db *gorm.DB
}

type DBFactory interface {
	GetDB() *gorm.DB
	GetDBWithTx() *gorm.DB
	Rollback(tx *gorm.DB)
	Commit(tx *gorm.DB)
}

// IsNotFound determines if an error is a record not found.
func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// Ping attempts a connection and closes it to the database.
func (db *DB) Ping(ctx context.Context) error {
	sqlDB, err := db.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Close releases database connections.
func (db *DB) Close(ctx context.Context) {
	logger := logging.FromContext(ctx)
	logger.Infof("Closing connection pool.")
	_db, err := db.db.DB()
	if err != nil {
		logger.Errorf("Cannot close db connection: %v", err.Error())
	}
	_db.Close()
}

func (db *DB) GetDB() *gorm.DB {
	return db.db
}

func (db *DB) GetDBWithTx() *gorm.DB {
	return db.db.Begin()
}

func (db *DB) Rollback(tx *gorm.DB) {
	tx.Rollback()
}

func (db *DB) Commit(tx *gorm.DB) {
	tx.Commit()
}
