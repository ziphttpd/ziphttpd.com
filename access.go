package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/ziphttpd/ziphttpd.com/util"
)

var accmu sync.Mutex

func addAccess(r *http.Request) int {
	accmu.Lock()
	defer accmu.Unlock()

	// 現在時間
	UTC := time.Now().UTC()
	UTCstr := util.Time2ymdhms(UTC)

	// Access はアクセスカウンタです
	type Access struct {
		Count int    `json:"count"`
		UTC   string `json:"utc"`
	}
	access := Access{
		Count: 0,
		UTC:   UTCstr,
	}

	needsave := false
	// アクセスカウント取得
	const key = confPRE + "Access"
	if jsonstr, err := ds.Get(r, key); err == nil {
		if err := json.Unmarshal([]byte(jsonstr), &access); err == nil {
			if old, err := util.Ymdhms2time(access.UTC); err == nil {
				dur := UTC.Sub(old)
				if dur.Minutes() > 10 {
					// 10分経過したのでセーブ予定
					needsave = true
				}
			}
			access.UTC = UTCstr
		}
	}
	access.Count++

	// バイト変換
	bytes, err := json.Marshal(&access)
	if err != nil {
		util.LogOutput(err.Error())
		return access.Count
	}

	// キャッシュに保存
	jsonstr := string(bytes)
	err = ds.SetCache(r, key, jsonstr)
	if err != nil {
		util.LogOutput(err.Error())
		return access.Count
	}

	if needsave {
		// 10分経過していたのでFirestoreに記録
		err := ds.Set(r, key, jsonstr)
		if err != nil {
			util.LogOutput(jsonstr)
			util.LogOutput(err.Error())
		}
	}
	return access.Count
}
