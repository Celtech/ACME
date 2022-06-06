package service

import (
	"baker-acme/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSecretKey(t *testing.T) {
	c := config.Init("testing")
	c.SetDefault("secret", "abcd1234")
	res := getSecretKey()
	assert.Equal(t, "abcd1234", res)
}

func TestGetSecretKeyNoDefault(t *testing.T) {
	config.Init("testing")
	res := getSecretKey()
	assert.Equal(t, "correct-horse-battery-staple", res)
}
