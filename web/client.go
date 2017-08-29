package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Get(url string, data url.Values) string {
	resp, err := http.Get(url + "?" + data.Encode())
	if err != nil {
		log.Error("get error", err, url)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("read error", err, url)
		return ""
	}
	return string(body)
}

func GetJson(url string, data url.Values) map[string]interface{} {
	body := Get(url, data)
	var ret map[string]interface{}
	if err := json.Unmarshal([]byte(body), &ret); err != nil {
		log.Error("get json decode error", err, body, url)
		return nil
	}
	return ret
}

func Post(url string, data url.Values) string {
	resp, err := http.PostForm(url, data)
	if err != nil {
		log.Error("post error", err, url)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("read error", err, url)
		return ""
	}
	return string(body)
}

func PostJson(url string, data url.Values) map[string]interface{} {
	body := Post(url, data)
	var ret map[string]interface{}
	if err := json.Unmarshal([]byte(body), &ret); err != nil {
		log.Error("post json decode error", err, body, url)
		return nil
	}
	return ret
}
