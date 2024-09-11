package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	acc, err := NewAccount("johndoe", "securepassword123")
	assert.Nil(t, err)
	fmt.Println(acc)
}
