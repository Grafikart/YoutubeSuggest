package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

type Cache interface {
	Get(key string) Subscriptions
	Set(key string, s Subscriptions) error
}

type FileCache struct {
	dir string
}

func NewFileCache(dir string) FileCache {
	return FileCache{dir}
}

func (c FileCache) Get(key string) Subscriptions {
	content, err := ioutil.ReadFile(path.Join(c.dir, key))
	if err != nil {
		return nil
	}
	s := &Subscriptions{}
	err = json.Unmarshal(content, s)
	if err != nil {
		return nil
	}
	return *s
}

func (c FileCache) Set(key string, s Subscriptions) error {
	b, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("can't convert %v to json, %v", s, err)
	}
	p := path.Join(c.dir, key)
	err = ioutil.WriteFile(p, b, 0644)
	if err != nil {
		return fmt.Errorf("can't write content to file %q, %v", p, err)
	}
	return nil
}
