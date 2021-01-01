package datastore

import (
	"fmt"

	"github.com/ziphttpd/ziphttpd.com/util"
)

var (
	mapcache map[string]string
	// NotMatchKey はキーにマッチする値が見つからなかった場合のエラーです。
	NotMatchKey error
)

func init() {
	NotMatchKey = fmt.Errorf("not match key")
}

func setMapcache(key, value string) error {
	if mapcache == nil {
		mapcache = map[string]string{}
	}
	mapcache[key] = value
	return nil
}

func setMapcacheMap(key string, data map[string]string) error {
	bytes, err := util.Map2BYTES(data)
	if err != nil {
		return err
	}
	return setMapcache(key, string(bytes))
}

func getMapcache(key string) (string, error) {
	if mapcache == nil {
		return "", NotMatchKey
	}
	if val, ok := mapcache[key]; ok {
		return val, nil
	}
	return "", NotMatchKey
}

func getMapcacheMap(key string) (map[string]string, error) {
	val, err := getMapcache(key)
	if err != nil {
		return nil, err
	}
	return util.JSON2Map(val)
}

func deleteMapcache(key string) error {
	delete(mapcache, key)
	return nil
}
