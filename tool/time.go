package tool

import "time"

func NowTime() int {
	return int(time.Now().Unix())
}
