package concurrency

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/mrgThang/learn-go/constant"
	db2 "github.com/mrgThang/learn-go/db"
)

const (
	NumberOfUsers = 10
)

func TestSendEmail(t *testing.T) {
	db, err := db2.InitDb()
	if err != nil {
		println(err.Error())
		panic(err)
	}

	defer func() {
		db2.TruncateDatabases(db, User{}.TableName())
	}()

	users := make([]*User, 0)
	for i := 0; i < NumberOfUsers; i++ {
		email := fmt.Sprintf("user%d@example.com", i)
		phone := "0999999999"
		users = append(users, &User{
			Name:         fmt.Sprintf("user%d", i),
			HashPassword: fmt.Sprintf("pass%d", i),
			Email:        &email,
			Phone:        &phone,
		})
	}

	err = db.CreateInBatches(users, constant.MaxCreateBatchSize).Error
	if err != nil {
		println(err.Error())
		panic(err)
	}

	concurrency(t, db)
	concurrencyWithUnbufferedChannel(t, db)
	goroutineListenToManyChannels(t, db)
}

func concurrency(t *testing.T, db *gorm.DB) {
	var users []*User
	err := db.Model(&User{}).Find(&users).Error
	if err != nil {
		println(err.Error())
		panic(err)
	}

	startTime := time.Now()

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	countEmail := 0
	wg.Add(constant.DefaultNumberOfWorkGroups)
	emailChan := make(chan string, constant.DefaultNumberOfWorkGroups)
	for i := 0; i < constant.DefaultNumberOfWorkGroups; i++ {
		go func() {
			defer wg.Done()
			for email := range emailChan {
				sendEmail(email)
				mutex.Lock()
				countEmail++
				mutex.Unlock()
			}
		}()
	}

	for _, user := range users {
		if user.Email == nil {
			continue
		}
		emailChan <- *user.Email
	}
	close(emailChan)
	wg.Wait()

	assert.Equal(t, NumberOfUsers, countEmail)

	println(fmt.Sprintf("Concurrency solution take %f", time.Since(startTime).Seconds()))
}

func concurrencyWithUnbufferedChannel(t *testing.T, db *gorm.DB) {
	var users []*User
	err := db.Model(&User{}).Find(&users).Error
	if err != nil {
		println(err.Error())
		panic(err)
	}

	startTime := time.Now()
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	countEmail := 0
	wg.Add(constant.DefaultNumberOfWorkGroups)
	emailChan := make(chan string)

	for i := 0; i < constant.DefaultNumberOfWorkGroups; i++ {
		go func() {
			defer wg.Done()
			for email := range emailChan {
				sendEmail(email)
				mutex.Lock()
				countEmail++
				mutex.Unlock()
			}
		}()
	}

	for _, user := range users {
		if user.Email == nil {
			continue
		}
		emailChan <- *user.Email
	}
	close(emailChan)
	wg.Wait()

	assert.Equal(t, NumberOfUsers, countEmail)

	println(fmt.Sprintf("Concurrency with unbuffered channel solution take %f", time.Since(startTime).Seconds()))
}

func goroutineListenToManyChannels(t *testing.T, db *gorm.DB) {
	var users []*User
	err := db.Model(&User{}).Find(&users).Error
	if err != nil {
		println(err.Error())
		panic(err)
	}

	startTime := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(constant.DefaultNumberOfWorkGroups)
	mutex := sync.Mutex{}
	countSent := 0
	zaloChan := make(chan string)
	smsChan := make(chan string)
	for i := 0; i < constant.DefaultNumberOfWorkGroups; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case phone, ok := <-zaloChan:
					if !ok {
						println("zalo channel closed")
						return
					}
					sendZaloMessage(phone)
					mutex.Lock()
					countSent++
					mutex.Unlock()
				case phone, ok := <-smsChan:
					if !ok {
						println("sms channel closed")
						return
					}
					sendSmsMessage(phone)
					mutex.Lock()
					countSent++
					mutex.Unlock()
				}
			}
		}()
	}

	for _, user := range users {
		if user.Phone == nil {
			continue
		}
		zaloChan <- *user.Phone
		smsChan <- *user.Phone
	}
	close(smsChan)
	close(zaloChan)
	wg.Wait()

	assert.Equal(t, NumberOfUsers*2, countSent)

	println(fmt.Sprintf("goroutine listen to many channels solution take %f", time.Since(startTime).Seconds()))
}

func sendEmail(email string) {
	time.Sleep(1 * time.Second)
}

func sendSmsMessage(phone string) {
	time.Sleep(1 * time.Second)
}

func sendZaloMessage(phone string) {
	time.Sleep(1 * time.Second)
}
