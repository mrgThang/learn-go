// example of use high order function to create middleware of api
package function

import (
	"fmt"
	"testing"
	"time"
)

func Retry(attemps int, delay time.Duration, function func() error) {
	for i := 0; i <= attemps; i++ {
		err := function()
		if err == nil {
			return
		}
		fmt.Printf("Execute function got err: %s\n", err.Error())
		if i == attemps {
			fmt.Printf("Exceed max retries\n")
		} else {
			time.Sleep(delay)
			fmt.Printf("Start retry %d times\n", i+1)
		}
	}
}

func x() error {
	return fmt.Errorf("this func always \n")
}

func TestRetry(t *testing.T) {
	Retry(10, time.Second, x)
}
