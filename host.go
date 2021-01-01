package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/labstack/echo/v4"
	"github.com/ziphttpd/zhsig/pkg/zhsig"
	"github.com/ziphttpd/ziphttpd.com/util"
)

// List はノートの記録形式です。
type List struct {
	Hosts []string `json:"hosts"`
}

// 一覧を取得
func getList(c echo.Context) (*List, error) {
	data := List{
		Hosts: []string{},
	}
	if jsonstr, _ := ds.Get(c.Request(), listPRE); jsonstr != "" {
		err := json.Unmarshal([]byte(jsonstr), &data)
		if err != nil {
			util.LogOutput(jsonstr)
			util.LogOutput(err.Error())
			return &data, fmt.Errorf("invalid json : %s", err.Error())
		}
	}
	return &data, nil
}

// ホストを登録
func addHost(c echo.Context, hostname string) error {
	var err error
	key := hostPRE + hostname
	if _, err := ds.Get(c.Request(), key); err == nil {
		// 登録済み
		// TODO: 公開鍵が無くなっていたら登録削除すべきか？
		return nil
	}
	// 未登録
	// 公開鍵署名チェック
	if _, err := zhsig.PublicKeyFetchFromURL(zhsig.NewHost("/", hostname)); err != nil {
		// 公開鍵を確認できなかった
		return err
	}
	// 登録
	if err = ds.Set(c.Request(), key, hostname); err != nil {
		return err
	}
	// リストの追加
	if data, err := getList(c); err == nil {
		// TODO: 重複を考慮する必要は本当に無いのか？
		data.Hosts = append(data.Hosts, hostname)
		sort.Strings(data.Hosts)
		if bytes, err := json.Marshal(data); err == nil {
			// リストの削除
			ds.Delete(c.Request(), listPRE)
			// 再登録
			ds.Set(c.Request(), listPRE, string(bytes))
		}
	}
	return err
}
