package test_infra

import (
	"database/sql"

	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/mnadev/limestone/storage"
)

type UnitTestSuite struct {
	suite.Suite
	StorageManager *storage.StorageManager
}

func (suite *UnitTestSuite) BeforeTest(suiteName, testName string) {
	sqlDB, err := sql.Open("ramsql", "Test"+testName)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if !db.Migrator().HasTable(&storage.AdhanFile{}) {
		db.AutoMigrate(&storage.AdhanFile{})
	}
	if !db.Migrator().HasTable(&storage.Event{}) {
		db.AutoMigrate(&storage.Event{})
	}
	if !db.Migrator().HasTable(&storage.Masjid{}) {
		db.AutoMigrate(&storage.Masjid{})
	}
	if !db.Migrator().HasTable(&storage.User{}) {
		db.AutoMigrate(&storage.User{})
	}
	if !db.Migrator().HasTable(&storage.NikkahProfile{}) {
		db.AutoMigrate(&storage.NikkahProfile{})
	}
	if !db.Migrator().HasTable(&storage.NikkahLike{}) {
		db.AutoMigrate(&storage.NikkahLike{})
	}
	if !db.Migrator().HasTable(&storage.NikkahMatch{}) {
		db.AutoMigrate(&storage.NikkahMatch{})
	}
	if !db.Migrator().HasTable(&storage.RevertProfile{}) {
		db.AutoMigrate(&storage.RevertProfile{})
	}
	if !db.Migrator().HasTable(&storage.RevertMatch{}) {
		db.AutoMigrate(&storage.RevertMatch{})
	}

	suite.StorageManager = &storage.StorageManager{
		DB: db,
	}
}
