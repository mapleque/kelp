package monitor

import (
	"encoding/json"
	"net/http"

	"github.com/kelp/crontab"
	"github.com/kelp/queue"
)

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
	w.Write(formatResponse(map[string]interface{}{
		"queue":    queue.GetInfo(),
		"crontab":  crontab.GetInfo(),
		"producer": queue.GetProducerInfo(),
		"consumer": queue.GetConsumerInfo()}))
}
