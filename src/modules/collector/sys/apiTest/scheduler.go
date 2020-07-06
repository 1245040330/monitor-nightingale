package apiTest

import (
	"fmt"
	"github.com/didi/nightingale/src/dataobj"
	"github.com/didi/nightingale/src/model"
	"github.com/didi/nightingale/src/modules/collector/sys/funcs"
	"github.com/toolkits/pkg/net/httplib"

	"time"
)

type ApiScheduler struct {
	Ticker *time.Ticker
	Api   *model.ApiCollect
	Quit   chan struct{}
}

func NewApiScheduler(p *model.ApiCollect) *ApiScheduler {
	scheduler := ApiScheduler{Api: p}
	scheduler.Ticker = time.NewTicker(time.Duration(p.Step) * time.Second)
	scheduler.Quit = make(chan struct{})
	return &scheduler
}

func (p *ApiScheduler) Schedule() {
	go func() {
		for {
			select {
			case <-p.Ticker.C:
				ApiCollect(p.Api)
			case <-p.Quit:
				p.Ticker.Stop()
				return
			}
		}
	}()
}

func (p *ApiScheduler) Stop() {
	close(p.Quit)
}
type Res struct {
	Code int64 `json:"code"`
	Data interface{} `json:"data"`
	Status string `json:"status"`
}
func ApiCollect(p *model.ApiCollect) {
	fmt.Println(p)
	var isOk=true
	//检测监控是否正常
	var res Res
	err := httplib.Get(p.Api).ToJSON(&res)
	if err != nil {
		err = fmt.Errorf("get collects from remote:%s failed, error:%v", p.Api, err)
		isOk=false
	}
	if(res.Code!=200 && res.Code!=20000 && res.Code!=2000){
		isOk=false
	}
	now := time.Now().Unix()
	item :=dataobj.MetricValue{}
	item.Metric="api.get"
	item.Endpoint = p.Ip
	item.Timestamp = now
	item.Step = int64(p.Step)
	if(isOk){
		item.Value=float64(200)
		item.ValueUntyped=float64(200)
	}else {
		item.Value=float64(0)
		item.ValueUntyped=float64(0)
	}
	item.CounterType="GAUGE"
	item.Tags="api="+p.Api
	items := []*dataobj.MetricValue{&item}

	funcs.Push(items)
}


