package concurrency

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var messages []string

func TestChatChannel(t *testing.T) {
	mutex := &sync.RWMutex{}
	wg := sync.WaitGroup{}
	go func() {
		readAction(mutex)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		writeAction(mutex)
	}()
	wg.Wait()
}

func readAction(mutex *sync.RWMutex) {
	readTimes := 0
	for {
		readTimes++
		wg := sync.WaitGroup{}
		wg.Add(10)
		for j := 0; j < 10; j++ {
			go func() {
				defer wg.Done()
				mutex.RLock()
				println(fmt.Sprintf("Read times %d: %d", readTimes, readMessage(messages)))
				mutex.RUnlock()
			}()
		}
		wg.Wait()
		time.Sleep(1 * time.Millisecond)
	}
}

func writeAction(mutex *sync.RWMutex) {
	for i := 0; i < 10; i++ {
		mutex.Lock()
		println("W")
		messages = writeMessage(messages)
		mutex.Unlock()
		time.Sleep(2 * time.Millisecond)
	}
}

func readMessage(messages []string) int {
	return len(messages)
}

func writeMessage(messages []string) []string {
	return append(messages, "hello world")
}
