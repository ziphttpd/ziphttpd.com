package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ziphttpd/ziphttpd.com/datastore"
	"github.com/ziphttpd/ziphttpd.com/util"
	"google.golang.org/appengine"
)

const (
	listPRE = "l." // リストのプリフィックス
	hostPRE = "h." // ホストのプリフィックス
	confPRE = "c." // 設定のキーのプリフィックス
)

var (
	proccesing map[string]bool
	mu         sync.Mutex
	ds         *datastore.DataStore
)

func main() {
	proccesing = map[string]bool{}
	ds = datastore.NewDataStore("Data")

	echoInst := echo.New()
	// CORS対応 (XMLHttpRequest で GET したときに、単純リクエストなのになぜかクロスサイトアクセスが必要)
	echoInst.Use(middleware.CORS())

	// トップページ
	echoInst.GET("/", topPage)
	// 広告一覧ページ
	echoInst.GET("/ad", adPage)
	// 取得 API
	echoInst.GET("/api/v1/list", getAPI)
	// 登録 API
	echoInst.GET("/api/v1/regist/:host", putAPI)

	http.Handle("/", echoInst)
	appengine.Main()
}

// メインページ
func topPage(c echo.Context) error {
	// パニックをキャッチ
	defer func() {
		r := recover()
		if r != nil {
			util.LogOutput(fmt.Sprintf("Run() return %+v", r))
		}
	}()
	conf, err := getConfig(c.Request())
	if err != nil {
		conf = &Config{}
	}
	util.LogOutput(fmt.Sprintf("%+v", conf))
	type TopParam struct {
		Acc     string
		Ad      AdParam
		URL     string
		Version string
		Date    string
	}
	param := TopParam{
		// total access
		Acc: util.Comma(addAccess(c.Request())),
		// 広告
		Ad: GetAdParam(c),
		// ファイル
		URL: conf.URL,
		// ファイルバージョン
		Version: conf.Version,
		// 日付
		Date: conf.Date,
	}
	return templateRender(http.StatusOK, "index", param, c)
}

// 広告確認ページ
func adPage(c echo.Context) error {
	type SignParam struct {
		Ad AdParam
	}
	param := SignParam{
		// 広告
		Ad: GetAllAdParam(c),
	}
	return templateRender(http.StatusOK, "ad", param, c)
}

// 取得 API
func getAPI(c echo.Context) error {
	list, err := getList(c)
	if err != nil {
		util.LogOutput(err.Error())
	}
	// 応答
	return c.JSON(http.StatusOK, list)
}

// 登録 API
func putAPI(c echo.Context) error {
	mu.Lock()
	// リクエストしてきたリモート
	requestIP := c.RealIP()
	if b, ok := proccesing[requestIP]; ok && b {
		mu.Unlock()
		// 処理中にまた来た
		// 不正なリクエストとしてエラー
		mes := "processing"
		util.LogOutput(mes)
		return c.String(http.StatusBadRequest, mes)
	}
	// → 処理中
	proccesing[requestIP] = true
	mu.Unlock()

	defer func() {
		mu.Lock()
		// → 処理終了
		proccesing[requestIP] = false
		mu.Unlock()
	}()

	hostName := c.Param("host")
	err := addHost(c, hostName)
	// 応答
	return err
}
