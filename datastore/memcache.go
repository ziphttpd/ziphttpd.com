package datastore

import (
	"context"
	"net/http"

	"github.com/ziphttpd/ziphttpd.com/util"
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

// キャッシュに書き込む
func setMemcache(r *http.Request, key string, value []byte) error {
	ctx, cancel := context.WithCancel(appengine.NewContext(r))
	defer cancel()
	item := &memcache.Item{
		Key:   key,
		Value: value,
	}
	return memcache.Set(ctx, item)
}
func setMemcacheMap(r *http.Request, key string, data map[string]string) error {
	bytes, err := util.Map2BYTES(data)
	if err != nil {
		return err
	}
	setMemcache(r, key, bytes)
	return nil
}

// キャッシュから読み込む
func getMemcache(r *http.Request, key string) (string, error) {
	ctx, cancel := context.WithCancel(appengine.NewContext(r))
	defer cancel()
	item, err := memcache.Get(ctx, key)
	if item != nil {
		return string(item.Value), err
	}
	return "", err
}
func getMemcacheMap(r *http.Request, key string) (map[string]string, error) {
	str, err := getMemcache(r, key)
	if err != nil {
		return nil, err
	}
	return util.JSON2Map(str)
}

// キャッシュから削除する
func deleteMemcache(r *http.Request, key string) error {
	ctx, cancel := context.WithCancel(appengine.NewContext(r))
	defer cancel()
	return memcache.Delete(ctx, key)
}
