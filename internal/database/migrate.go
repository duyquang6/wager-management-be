package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var gormMigrationScripts = []*gormigrate.Migration{
	{
		ID: "00001-CreateWager",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`
CREATE TABLE IF NOT EXISTS wagers
(
    id                    bigint unsigned  not null auto_increment,
    created_at            datetime         not null default CURRENT_TIMESTAMP,
    updated_at            datetime on update CURRENT_TIMESTAMP,
    deleted_at            datetime,
    total_wager_value     int unsigned     not null,
    odds                  int unsigned     not null,
    selling_percentage    tinyint unsigned not null,
    selling_price         decimal(24, 12)   not null,
    current_selling_price decimal(24, 12)   not null,
    percentage_sold       tinyint unsigned,
    amount_sold           decimal(24, 12) unsigned,
    placed_at             datetime         not null,
    KEY idx_wagers_deleted_at (deleted_at) USING BTREE,
    PRIMARY KEY (id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;
				`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("wagers")
		},
	},
	{
		ID: "00002-CreatePurchase",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`
CREATE TABLE IF NOT EXISTS purchases
(
    id           bigint unsigned not null auto_increment,
    created_at   datetime        not null default CURRENT_TIMESTAMP,
    updated_at   datetime on update CURRENT_TIMESTAMP,
    deleted_at   datetime,
    wager_id     bigint unsigned not null,
    buying_price decimal(24, 12)  not null,
    bought_at    datetime        not null,
    KEY idx_wagers_deleted_at (deleted_at) USING BTREE,
    FOREIGN KEY (wager_id) REFERENCES wagers (id),
    PRIMARY KEY (id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;
				`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("purchases")
		},
	},
}

// Migrate migrate schema
func (_db *DB) Migrate(ctx context.Context) error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	db := _db.GetDB().Session(&gorm.Session{Logger: newLogger})

	m := gormigrate.New(db, gormigrate.DefaultOptions, gormMigrationScripts)
	return m.Migrate()
}

// MigrateDown rollback schema
func (_db *DB) MigrateDown(ctx context.Context) error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	db := _db.GetDB().Session(&gorm.Session{Logger: newLogger})

	m := gormigrate.New(db, gormigrate.DefaultOptions, gormMigrationScripts)

	for i := range gormMigrationScripts {
		migration := gormMigrationScripts[len(gormMigrationScripts) - i - 1]
		err := m.RollbackMigration(migration)
		if err != nil {
			newLogger.Error(ctx, "cannot rollback script", migration.ID)
			return err
		}
	}

	return nil
}
