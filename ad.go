package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ziphttpd/ziphttpd.com/datastore"
	"github.com/ziphttpd/ziphttpd.com/util"
)

var (
	ad728ds *datastore.DataStore
	ad125ds *datastore.DataStore
)

func init() {
	ad728ds = datastore.NewDataStore("Ad728")
	ad125ds = datastore.NewDataStore("Ad125")
}

// AdkeyList は広告キーの情報です
type AdkeyList struct {
	UTC    string   `json:"utc"`
	Sdkeys []string `json:"sdkeys"`
}

func getAdkey(c echo.Context, ads *datastore.DataStore, szcode string) []string {
	adkeys := AdkeyList{Sdkeys: []string{}}
	cur := time.Now().UTC()
	expr := true
	cachekey := confPRE + "adkey" + szcode
	// 広告キーリストを取得
	adkeyjson, err := ds.Get(c.Request(), cachekey)
	if err == nil {
		if err = json.Unmarshal([]byte(adkeyjson), &adkeys); err == nil {
			// adkey をデコードできた
			if old, err := util.Ymdhms2time(adkeys.UTC); err == nil {
				dur := cur.Sub(old)
				if dur.Hours() < 1 {
					// 満了していない
					expr = false
				}
			}
		}
	}
	if expr {
		util.LogOutput("refresh ad list")
		adkeys.UTC = util.Time2ymdhms(cur)
		// 広告データベースからキーを取得
		adkeys.Sdkeys = ads.Keys()
		if bytes, err := json.Marshal(&adkeys); err == nil {
			// キャッシュ保存
			_ = ds.SetCache(c.Request(), cachekey, string(bytes))
		}
	}
	return adkeys.Sdkeys
}

// AdParamElem 広告テンプレートパラメータ要素
type AdParamElem struct {
	URL string
}

// AdParam 広告テンプレートパラメータ
type AdParam struct {
	Ad728 []AdParamElem
	Ad125 []AdParamElem
}

func getAdParam(c echo.Context, ads *datastore.DataStore, adkeys []string, reqsize int) []AdParamElem {
	ret := []AdParamElem{}
	size := len(adkeys)
	keys := map[string]string{}
	for {
		key := adkeys[rand.Intn(size)%size]
		html, err := ads.Get(c.Request(), key)
		if err != nil {
			util.LogOutput(err.Error())
			return ret
		}
		if html != "" {
			elem := strings.Split(html, "\"")
			if len(elem) > 1 && elem[1] != "" {
				// ...."xxxxxxx".... の xxxxxxx
				keys[key] = elem[1]
			} else {
				util.LogOutput(fmt.Sprintf("invalid ad data %s: %s", key, html))
				continue
			}
		}
		if len(keys) >= reqsize {
			break
		}
	}
	for _, value := range keys {
		// 埋め込みタグからsrc=""の中身
		ret = append(ret, AdParamElem{URL: value})
	}
	return ret
}

// GetAdParam はテンプレートに渡すパラメータ情報を構築します。
func GetAdParam(c echo.Context) AdParam {
	param := AdParam{}
	rand.Seed(time.Now().UnixNano())
	param.Ad728 = getAdParam(c, ad728ds, getAdkey(c, ad728ds, "728"), 1)
	param.Ad125 = getAdParam(c, ad125ds, getAdkey(c, ad125ds, "125"), 5)
	return param
}

// GetAllAdParam はテンプレートに渡すパラメータ情報を構築します。
func GetAllAdParam(c echo.Context) AdParam {
	param := AdParam{}
	rand.Seed(time.Now().UnixNano())
	param.Ad728 = listAll(c, ad728ds, getAdkey(c, ad728ds, "728"))
	param.Ad125 = listAll(c, ad125ds, getAdkey(c, ad125ds, "125"))
	return param
}

func listAll(c echo.Context, ads *datastore.DataStore, adkeys []string) []AdParamElem {
	ret := []AdParamElem{}
	size := len(adkeys)
	for i := 0; i < size; i++ {
		key := adkeys[i]
		html, err := ads.Get(c.Request(), key)
		if err != nil {
			util.LogOutput(err.Error())
			return ret
		}
		if html != "" {
			elem := strings.Split(html, "\"")
			if len(elem) > 1 && elem[1] != "" {
				// ...."xxxxxxx".... の xxxxxxx
				value := elem[1]
				// 埋め込みタグからsrc=""の中身
				ret = append(ret, AdParamElem{URL: value})
			} else {
				util.LogOutput(fmt.Sprintf("invalid ad data %s: %s", key, html))
				continue
			}
		}
	}
	return ret
}
