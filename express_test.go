package task

import (
	"testing"
	"fmt"
)

func TestParseExpress(t *testing.T){
	fmt.Println(parseExpress("1-5", ExpressType_WeekDay))
}