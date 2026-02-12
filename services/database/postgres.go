package database

import (
	"api/constant"
	"api/migrations"
	"api/services/cache"
	"fmt"
	"time"

	"github.com/morkid/gocache"
	cache_redis "github.com/morkid/gocache-redis/v8"
	"github.com/morkid/paginate"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DB Main database connection
var DB *gorm.DB

// PG Pagination library
var PG *paginate.Pagination

// InitDatabase initialize database connection
func Connect() {
	if nil == DB {
		db := dbConnect()
		if nil != db {
			DB = db

			var cacheAdapter *gocache.AdapterInterface
			cacheSeconds := viper.GetInt64("CACHE_TTL_SECONDS")

			if nil != cache.Redis && cacheSeconds > 0 {
				cacheAdapter = cache_redis.NewRedisCache(cache_redis.RedisCacheConfig{
					Client:    cache.Redis,
					ExpiresIn: time.Duration(cacheSeconds) * time.Second,
				})
			}

			PG = paginate.New(&paginate.Config{
				CacheAdapter:         cacheAdapter,
				FieldSelectorEnabled: true,
			})
			// only run dbMigrate() while "enable_migration" is set to true
			if viper.GetBool("ENABLE_MIGRATION") {
				dbMigrate()
			}
		}
	}
}

func dbConnect() *gorm.DB {
	logLevel := logger.Info

	switch constant.GetEnvironment() {
	case constant.EnvironmentLocal, constant.EnvironmentDevelopment:
		logLevel = logger.Info
	case constant.EnvironmentStaging:
		logLevel = logger.Error
	case constant.EnvironmentProduction:
		logLevel = logger.Silent
	}

	config := gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   viper.GetString("DB_TABLE_PREFIX"),
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_NAME"),
		viper.GetString("TIMEZONE"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &config)

	if nil != err {
		panic(err)
	}

	if nil != db {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(1)
		sqlDB.SetConnMaxLifetime(time.Second * 5)
	}

	return db
}

func dbMigrate() {
	db := dbConnect()
	if nil != db && len(migrations.ModelMigrations) > 0 {
		err := db.AutoMigrate(migrations.ModelMigrations...)

		if nil != err {
			panic(err)
		}

		seeds := migrations.DataSeeds()
		if len(seeds) > 0 {
			for i := range seeds {
				tx := db.Begin()

				defer func() {
					if r := recover(); r != nil {
						tx.Rollback()
					}
				}()

				if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(seeds[i]).Error; nil != err {
					tx.Rollback()
				}

				if err := tx.Commit().Error; nil != err {
					tx.Rollback()
				}
			}
		}

		db.Migrator().DropTable("schema_migration")

		sqlDB, _ := db.DB()

		defer sqlDB.Close()
	}
}
