package database

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var once sync.Once
var dba *gorm.DB

// GetInstance - Returns a DB instance
func GetInstance() *gorm.DB {
	once.Do(func() {
		//documentation  - https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL
		// user
		user := viper.GetString("database.user")
		// password
		password := viper.GetString("database.password")
		// hots
		host := viper.GetString("database.host")
		// port
		port := viper.GetString("database.port")
		// database name
		dbname := viper.GetString("database.dbname")
		// sslmode
		sslmode := viper.GetString("database.sslmode")
		//TimeZone
		timezone := viper.GetString("database.timezone")
		//Schema
		shcema := viper.GetString("database.schema")

		//dsn
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s search_path=%s", host, user, password, dbname, port, sslmode, timezone, shcema)
		//database configuration
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			Logger:                 logger.Default.LogMode(logger.Info),
		})
		// instance of db
		dba = db
		if err != nil {
			logrus.Panicf("Error connecting to the database at %s:%s/%s", host, port, dbname)
		}
		//print log when database connection is established
		logrus.Infof("Successfully Established Connection to -- %s:%s/%s", host, port, dbname)
	})
	//return
	return dba
}

// Migrate - AutoMigrate models
func Migrate(models ...interface{}) {
	if dba == nil {
		GetInstance()
	}
	if err := dba.AutoMigrate(models...); err != nil {
		logrus.Errorf("Migration failed: %s", err)
	} else {
		logrus.Info("Database Migration Completed Successfully")
	}
}
