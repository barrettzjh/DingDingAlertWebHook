package loki

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	values sync.Map
	expiry time.Duration
}

func NewCache(expiry time.Duration) *Cache {
	return &Cache{
		expiry: expiry,
	}
}

func (c *Cache) setValue(obj Alert) {
	key := getObjectKey(obj)
	c.values.Store(key, time.Now().Add(c.expiry))
}

func (c *Cache) getValue(obj Alert) (bool, error) {
	key := getObjectKey(obj)
	value, ok := c.values.Load(key)
	if !ok {
		return false, nil
	}

	expiryTime, ok := value.(time.Time)
	if !ok {
		return false, fmt.Errorf("value has incorrect type")
	}
	// 过期删除，返回false
	if time.Now().After(expiryTime) {
		c.values.Delete(key)
		return false, nil
	}

	return true, nil
}

func getObjectKey(obj Alert) string {
	h := md5.New()
	h.Write([]byte(obj.Labels.Body))
	return hex.EncodeToString(h.Sum(nil))
}
