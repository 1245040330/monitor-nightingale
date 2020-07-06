package apiTest

import "github.com/didi/nightingale/src/model"

var (
	Apis              = make(map[int]*model.ApiCollect)
	ApisWithScheduler = make(map[int]*ApiScheduler)
)

func DelNoProcCollect(newCollect map[int]*model.ApiCollect) {
	for currKey, currProc := range Apis {
		newProc, ok := newCollect[currKey]
		if !ok || currProc.LastUpdated != newProc.LastUpdated {
			deleteApi(currKey)
		}
	}
}

func AddNewProcCollect(newCollect map[int]*model.ApiCollect) {
	for target, newProc := range newCollect {
		if _, ok := Apis[target]; ok && newProc.LastUpdated == Apis[target].LastUpdated {
			continue
		}

		Apis[target] = newProc
		sch := NewApiScheduler(newProc)
		ApisWithScheduler[target] = sch
		sch.Schedule()
	}
}

func deleteApi(key int) {
	v, ok := ApisWithScheduler[key]
	if ok {
		v.Stop()
		delete(ApisWithScheduler, key)
	}
	delete(Apis, key)
}