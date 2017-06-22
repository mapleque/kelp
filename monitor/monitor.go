package monitor

import (
	"encoding/json"
	"net/http"
)

type Observable interface {
	GetInfo() interface{}
}

var observeList map[string]Observable

func init() {
	observeList = make(map[string]Observable)
}

func Observe(name string, subject Observable) {
	observeList[name] = subject
}

func Run(host string) {
	log.Info("monitor starting ...")
	mux := http.NewServeMux()

	mux.HandleFunc("/info", getInfo)

	err := http.ListenAndServe(host, mux)
	if err != nil {
		panic("monitor start faild " + err.Error())
	}
}

func formatResponse(res interface{}) []byte {
	ret, err := json.Marshal(res)
	if err != nil {
		log.Error("can not format response", res)
	}
	return ret
}

func getInfo(w http.ResponseWriter, r *http.Request) {
	ret := make(map[string]interface{})
	for name, subject := range observeList {
		ret[name] = subject.GetInfo()
	}
	w.Write(formatResponse(ret))
}
