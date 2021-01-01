package datastore

import (
	"net/http"

	"github.com/ziphttpd/ziphttpd.com/firebase"
	"github.com/ziphttpd/ziphttpd.com/util"
)

const (
	dataKey = "data"
)

// DataStore はデータを管理します。
type DataStore struct {
	name string
	fs   *firebase.Collection
}

// NewDataStore は DataStore を生成します。
func NewDataStore(collectionname string) *DataStore {
	return &DataStore{name: collectionname + "/", fs: firebase.NewCollection(collectionname)}
}

// SetCacheMap はキーに対するマップを記録します。
func (d *DataStore) SetCacheMap(r *http.Request, key string, data map[string]string) error {
	cachekey := d.name + key
	bytes, err := util.Map2BYTES(data)
	if err != nil {
		util.LogOutput(err.Error())
		return err
	}

	err = setMemcache(r, cachekey, bytes)
	if err != nil {
		util.LogOutput(err.Error())
		return err
	}

	setMapcache(cachekey, string(bytes))

	return nil
}

// SetCache はキーに対する値を記録します。
func (d *DataStore) SetCache(r *http.Request, key, value string) error {
	data := map[string]string{dataKey: value}
	return d.SetCacheMap(r, key, data)
}

// SetMap はキーに対するマップを記録します。
func (d *DataStore) SetMap(r *http.Request, key string, data map[string]string) error {
	// Firestore に書き込む
	err := d.fs.SetMap(key, data)
	if err != nil {
		util.LogOutput(err.Error())
		return err
	}
	return d.SetCacheMap(r, key, data)
}

// Set はキーに対する値を記録します。
func (d *DataStore) Set(r *http.Request, key, value string) error {
	data := map[string]string{dataKey: value}
	return d.SetMap(r, key, data)
}

// GetMap はキーに対してのマップを取得します。
func (d *DataStore) GetMap(r *http.Request, key string) (map[string]string, error) {
	cachekey := d.name + key
	data, err := getMapcacheMap(cachekey)
	if len(data) != 0 {
		// キャッシュにあったので返す
		return data, err
	}

	// キャッシュを検索
	// logOutput("memcache.Get")
	data, err = getMemcacheMap(r, cachekey)
	if len(data) != 0 {
		// キャッシュにあったので返す
		setMapcacheMap(cachekey, data)
		return data, err
	}

	// キャッシュにないので Firestore から取得する
	// logOutput("getFirestore")
	data, err = d.fs.GetMap(key)
	if err != nil {
		// Firestore にもなかった
		util.LogOutput(err.Error())
		return nil, err
	}

	// キャッシュに書き込む
	setMapcacheMap(cachekey, data)
	// logOutput("memcache.Set")
	err = setMemcacheMap(r, cachekey, data)
	if err != nil {
		util.LogOutput(err.Error())
		return nil, err
	}

	return data, err
}

// Get はキーに対しての値を取得します。
func (d *DataStore) Get(r *http.Request, key string) (string, error) {
	data, err := d.GetMap(r, key)
	if err != nil {
		return "", err
	}
	if val, ok := data[dataKey]; ok {
		return val, nil
	}
	return "", NotMatchKey
}

// Delete はキーのデータを削除します。
func (d *DataStore) Delete(r *http.Request, key string) error {
	cachekey := d.name + key
	deleteMapcache(cachekey)
	deleteMemcache(r, cachekey)
	d.fs.Delete(key)
	return nil
}

// Keys はキーの一覧を返します。
func (d *DataStore) Keys() []string {
	return d.fs.Keys()
}
