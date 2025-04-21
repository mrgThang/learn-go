// example of using curry function to create validator function
package function

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func makeMinLengthValidator(min int) func(string) bool {
	return func(input string) bool {
		return len(input) >= min
	}
}

var validators = map[string]func(string) bool{
	"username": makeMinLengthValidator(5),
	"password": makeMinLengthValidator(8),
}

func TestValidator(t *testing.T) {
	assert.Equal(t, false, validators["username"]("user"))
	assert.Equal(t, true, validators["username"]("username"))
	assert.Equal(t, true, validators["password"]("secret123"))
}
