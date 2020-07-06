package stra

import (
	"fmt"
	"github.com/didi/nightingale/src/model"
	"github.com/didi/nightingale/src/toolkits/str"
	"github.com/toolkits/pkg/net/httplib"
	"reflect"
	"time"
)


func NewApiCollect( step int, name,api string, modTime time.Time) *model.ApiCollect {
	return &model.ApiCollect{
		CollectType:   "api",
		Name:        name,
		Step:          step,
		Api:          api,
		LastUpdated:   modTime,
	}
}
type Res struct {
	Dat []model.ApiCollect `json:"dat"`
	Err string        `json:"err"`
}

func copyPoint(m *model.ApiCollect) *model.ApiCollect{
	vt := reflect.TypeOf(m).Elem()
	fmt.Println(vt)
	newoby := reflect.New(vt)
	newoby.Elem().Set(reflect.ValueOf(m).Elem())
	return newoby.Interface().(*model.ApiCollect)
}
func GetApiCollects() map[int]*model.ApiCollect {
	apis := make(map[int]*model.ApiCollect)
	var res Res
	var url ="http://fengjie.info/api/portal/collect/apiAllList"
	if StraConfig.Enable {
		err := httplib.Get(url).ToJSON(&res)
		if err != nil {
			err = fmt.Errorf("get collects from remote:%s failed, error:%v", url, err)
		}
		for k, v := range res.Dat {
			vP:=copyPoint(&v)
			apis[k] = vP
		}
		fmt.Println(res.Dat)
		for _, p := range res.Dat {
			tagsMap := str.DictedTagstring("api")
			tagsMap["api"] = p.Api

			p.Comment = str.SortedTags(tagsMap)
		}
	}
	return apis
}

