package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_SetCache(t *testing.T) {
	var (
		key   = "foo"
		value = "bar"
	)
	Set(key, value, time.Millisecond*100)
	v, err := Get(key)
	assert.Nil(t, err)
	assert.Equal(t, value, v)

	time.Sleep(time.Millisecond * 200) //after expired -> get not found -> error
	v, err = Get(key)
	assert.Equal(t, NilCacheError, err)
	assert.Nil(t, v)
}

func Test_DelCache(t *testing.T) {
	var (
		key   = "foo"
		value = "bar"
	)
	Set(key, value, time.Millisecond*100)
	v, err := Get(key)
	assert.Nil(t, err)
	assert.Equal(t, value, v)
	existed := Delete(key)
	assert.True(t, existed)

	time.Sleep(time.Millisecond * 200) //after expired -> deleted -> not existed
	existed = Delete(key)
	assert.False(t, existed)
}

func Test_GetCache(t *testing.T) {
	var (
		key   = "foo"
		value = "bar"
	)
	Set(key, value, time.Millisecond*100)
	v, err := Get(key)
	assert.Nil(t, err)
	assert.Equal(t, value, v)

	time.Sleep(time.Millisecond * 50)
	v, err = Get(key)
	assert.Nil(t, err)
	assert.Equal(t, value, v)

	time.Sleep(time.Millisecond * 200)
	v, err = Get(key)
	assert.Equal(t, NilCacheError, err)
	assert.Nil(t, v)
	existed := Delete(key)
	assert.False(t, existed)
}
