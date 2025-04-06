package concurrency

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/mrgThang/learn-go/constant"
	db2 "github.com/mrgThang/learn-go/db"
)

const (
	LenUsers = 100000
)

func TestGetUserInBatches(t *testing.T) {
	db, err := db2.InitDb()
	if err != nil {
		t.Error("init db error", err)
	}
	defer func() {
		db2.TruncateDatabases(db, []string{User{}.TableName()}...)
	}()

	newUsers := make([]*User, 0)
	for i := 0; i < LenUsers; i++ {
		newUsers = append(newUsers, &User{
			Name:         fmt.Sprintf("user%d", i),
			HashPassword: fmt.Sprintf("pass%d", i),
		})
	}

	db.CreateInBatches(newUsers, constant.MaxCreateBatchSize)

	solutionOneGetAll(db, t)

	solutionTwoGetInBatches(db, t)

	solutionThreeGetInBatchesManual(db, t)

	solutionFourGetInBatchesUsingGoroutine(db, t)
}

func solutionOneGetAll(db *gorm.DB, t *testing.T) {
	users := make([]*User, 0)

	start := time.Now()

	err := db.Model(User{}).Find(&users).Error
	if err != nil {
		t.Error("Get users got err", err)
	}

	assert.Equal(t, LenUsers, len(users))

	t.Log(fmt.Sprintf("========== Solution 1 one get all take: %f s", time.Since(start).Seconds()))
}

func solutionTwoGetInBatches(db *gorm.DB, t *testing.T) {
	users := make([]*User, 0)
	var batchUsers []*User

	start := time.Now()

	err := db.Model(User{}).FindInBatches(&batchUsers, constant.MaxGetBatchSize, func(tx *gorm.DB, batch int) error {
		users = append(users, batchUsers...)
		return nil
	}).Error
	if err != nil {
		t.Error("Get users got err", err)
	}

	assert.Equal(t, LenUsers, len(users))

	t.Log(fmt.Sprintf("========== Solution 2 get in batches use lib take: %f s", time.Since(start).Seconds()))
}

func solutionThreeGetInBatchesManual(db *gorm.DB, t *testing.T) {
	users := make([]*User, 0)
	start := time.Now()

	for i := 0; i < LenUsers; i += constant.MaxGetBatchSize {
		currentUsers := make([]*User, 0)
		err := db.Model(User{}).Where("id > ?", i).Limit(constant.MaxGetBatchSize).Find(&currentUsers).Error
		if err != nil {
			t.Error("Get users using goroutine got err", err, zap.Any("lastId", i))
		}
		users = append(users, currentUsers...)
	}

	assert.Equal(t, LenUsers, len(users))

	t.Log(fmt.Sprintf("========== Solution 3 get in batches manual take: %f s", time.Since(start).Seconds()))
}

func solutionFourGetInBatchesUsingGoroutine(db *gorm.DB, t *testing.T) {
	var mu sync.Mutex
	users := make([]*User, 0)
	start := time.Now()

	lastIdChan := make(chan int, constant.MaxConcurrenceProcessTruncateTable)
	wg := sync.WaitGroup{}
	wg.Add(constant.MaxConcurrenceProcessTruncateTable)
	for w := 0; w < constant.MaxConcurrenceProcessTruncateTable; w++ {
		go func() {
			defer wg.Done()
			currentUsers := make([]*User, 0)
			for lastId := range lastIdChan {
				err := db.Model(User{}).Where("id > ?", lastId).Limit(constant.MaxGetBatchSize).Find(&currentUsers).Error
				if err != nil {
					t.Error("Get users using goroutine got err", err, zap.Any("lastId", lastId))
				}
				mu.Lock()
				users = append(users, currentUsers...)
				mu.Unlock()
			}
		}()
	}

	for i := 0; i < LenUsers; i += constant.MaxGetBatchSize {
		lastIdChan <- i
	}
	close(lastIdChan)
	wg.Wait()

	assert.Equal(t, LenUsers, len(users))

	t.Log(fmt.Sprintf("========== Solution 4 get in batches using goroutinne take: %f s", time.Since(start).Seconds()))
}
