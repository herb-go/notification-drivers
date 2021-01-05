package emaildelivery

import (
	"mime"
)

func MustEscape(unescaped string) string {
	return mime.BEncoding.Encode("utf-8", unescaped)
}
