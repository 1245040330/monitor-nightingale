package apiTest

import (
	"github.com/didi/nightingale/src/modules/collector/stra"
	"time"
)

func Detect() {
	detect()
	go loopDetect()
}

func loopDetect() {
	for {
		time.Sleep(time.Second * 10)
		detect()
	}
}

func detect() {
	ps := stra.GetApiCollects()
	DelNoProcCollect(ps)
	AddNewProcCollect(ps)
}