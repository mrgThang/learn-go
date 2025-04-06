package db

import (
	"fmt"
	osLog "log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/mrgThang/learn-go/config"
	"github.com/mrgThang/learn-go/constant"
)

func InitDb() (*gorm.DB, error) {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	newLogger := logger.New(
		osLog.New(os.Stderr, "", osLog.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Warn,
			Colorful:      true,
		})

	db, err := gorm.Open(mysql.Open(cfg.Mysql.DSN()), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	})
	if err != nil {
		panic(err)
	}

	return db, err
}

func TruncateDatabases(db *gorm.DB, tableNames ...string) {
	if db == nil {
		panic("db is nil")
	}
	tableNameChan := make(chan string, constant.MaxConcurrenceProcessTruncateTable)
	wg := sync.WaitGroup{}
	wg.Add(constant.MaxConcurrenceProcessTruncateTable)
	for w := 0; w < constant.MaxConcurrenceProcessTruncateTable; w++ {
		go func() {
			defer wg.Done()
			for tableName := range tableNameChan {
				sqlCmd := fmt.Sprintf("TRUNCATE TABLE `%s`", tableName)
				err := db.Exec(sqlCmd).Error
				if err != nil {
					panic(err)
				}
			}
		}()
	}
	for _, tableName := range tableNames {
		tableNameChan <- tableName
	}
	close(tableNameChan)
	wg.Wait()
}
