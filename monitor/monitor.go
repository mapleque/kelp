package monitor

import (
	"encoding/json"
	"net/http"

	"github.com/kelp/crontab"
	"github.com/kelp/log"
	"github.com/kelp/queue"
)

func Run(host string) {
	log.Info("monitor starting ...")
	mux := http.NewServeMux()

	mux.HandleFunc("/queue", getQueueInfo)
	mux.HandleFunc("/producer", getProducerInfo)
	mux.HandleFunc("/consumer", getConsumerInfo)
	mux.HandleFunc("/crontab", getCrontabInfo)

	err := http.ListenAndServe(host, mux)
	if err != nil {
		log.Fatal("monitor start faild", err.Error())
	}
}

func formatResponse(res interface{}) []byte {
	ret, err := json.Marshal(res)
	if err != nil {
		log.Error("can not format response", res)
	}
	return ret
}

func getQueueInfo(w http.ResponseWriter, r *http.Request) {
	w.Write(formatResponse(queue.GetInfo()))
}
func getProducerInfo(w http.ResponseWriter, r *http.Request) {
	w.Write(formatResponse(queue.GetProducerInfo()))
}
func getConsumerInfo(w http.ResponseWriter, r *http.Request) {
	w.Write(formatResponse(queue.GetConsumerInfo()))
}
func getCrontabInfo(w http.ResponseWriter, r *http.Request) {
	w.Write(formatResponse(crontab.GetInfo()))
}
