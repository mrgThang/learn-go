package concurrency

import (
	"sync"
	"testing"
	"time"
)

var messages []string

func TestChatChannel(t *testing.T) {
	mutex := &sync.RWMutex{}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		readAction(mutex)
	}()
	go func() {
		defer wg.Done()
		writeAction(mutex)
	}()

	wg.Wait()
}

func readAction(mutex *sync.RWMutex) {
	for i := 0; i < 10; i++ {
		wg := sync.WaitGroup{}
		wg.Add(10)
		for j := 0; j < 10; j++ {
			go func() {
				defer wg.Done()
				mutex.RLock()
				print(readMessage(messages))
				mutex.RUnlock()
				time.Sleep(1 * time.Millisecond)
			}()
		}
		wg.Wait()
	}
}

func writeAction(mutex *sync.RWMutex) {
	for i := 0; i < 10; i++ {
		mutex.Lock()
		print("W")
		messages = writeMessage(messages)
		time.Sleep(1 * time.Millisecond)
		mutex.Unlock()
	}
}

func readMessage(messages []string) int {
	return len(messages)
}

func writeMessage(messages []string) []string {
	return append(messages, "hello world")
}
