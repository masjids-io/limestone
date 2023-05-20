package test_infra

import (
	"database/sql"

	"github.com/mnadev/limestone/storage"
	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	db.AutoMigrate(&storage.User{})
	db.AutoMigrate(&storage.Masjid{})
	db.AutoMigrate(&storage.Event{})

	suite.StorageManager = &storage.StorageManager{
		DB: db,
	}
}
