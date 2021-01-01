package util

import (
	"fmt"
	"log"
	"runtime"
)

// LogOutput はログを出力します。
func LogOutput(mes string) {
	called := GetCallerPC().String()
	log.Output(2, called+": "+mes)
}

// CalledPC は呼び出し元を示す識別コードです。
type CalledPC uintptr

// String はstringerです。
func (c CalledPC) String() string {
	pc := uintptr(c)
	fpc := runtime.FuncForPC(pc)
	n, l := fpc.FileLine(pc)
	return fmt.Sprintf("%s (%s:%d)", fpc.Name(), n, l)
}

// GetCallerPC はそれを呼び出した関数の呼び出し元の PC を取得します。
func GetCallerPC() CalledPC {
	pc, _, _, _ := runtime.Caller(2)
	return CalledPC(pc)
}
