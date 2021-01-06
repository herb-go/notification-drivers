package wechatworkdelivery

import (
	"errors"
	"fmt"
)

func NewInvalidMsgType(msgtype string) error {
	return fmt.Errorf("wechatworkdelivery: invalid msgtype [%s]", msgtype)
}

var ErrNewsFormatError = errors.New("wechatworkdelivery: news format error")
