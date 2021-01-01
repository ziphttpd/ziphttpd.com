package firebase

import (
	"context"
	"fmt"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/ziphttpd/ziphttpd.com/util"
)

// VerifyUser はIDトークンを検証した結果を返します。
func VerifyUser(w http.ResponseWriter, reqToken string) (*auth.Token, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app, err := firebase.NewApp(ctx, nil, GetOption())
	if err != nil {
		util.LogOutput(err.Error())
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return nil, err
	}
	auth, err := app.Auth(ctx)
	if err != nil {
		util.LogOutput(err.Error())
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return nil, err
	}
	token, err := auth.VerifyIDToken(ctx, reqToken)
	if err != nil {
		util.LogOutput(err.Error())
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return nil, err
	}
	return token, nil
}

// DeleteUser はユーザ情報を抹消します。
func DeleteUser(w http.ResponseWriter, token *auth.Token) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app, err := firebase.NewApp(ctx, nil, GetOption())
	if err != nil {
		util.LogOutput(err.Error())
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return err
	}
	auth, err := app.Auth(ctx)
	if err != nil {
		util.LogOutput(err.Error())
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return err
	}
	err = auth.DeleteUser(ctx, token.UID)
	if err != nil {
		util.LogOutput(err.Error())
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return err
	}
	return nil
}
