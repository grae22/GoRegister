package utils

import (
	"fmt"
	"strconv"
)

func CentsToRandsStr(c int) string {
	r := c / 100
	c = c % 100

	if c == 0 {
		return strconv.Itoa(r)
	} else {
		return fmt.Sprintf("%d.%0*d", r, c, 2)
	}
}
