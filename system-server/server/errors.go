package myServer

import (
	"fmt"
	"github.com/iWyh2/myTcpSystem-server/utils"
)

// ErrMsg 错误提示
func ErrMsg(err error) {
	fmt.Printf("[%s] system> %v\n", utils.TimeStr(), err.Error())
}
