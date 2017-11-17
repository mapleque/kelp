package web

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
)

func PostJson(url string, data []byte) ([]byte, error) {
	log.Info("post", url, len(data))
	body := bytes.NewReader(data)
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return []byte(""), err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}
	return ret, nil
}

func PostForm(url string, values url.Values) ([]byte, error) {
	resp, err := http.PostForm(url, values)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}
	return body, nil
}

func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}
	return body, nil
}

func Mail(from, password, host, to, subject, body string) error {
	domain := strings.Split(host, ":")[0]
	auth := smtp.PlainAuth("", from, password, domain)
	contentType := "Content-Type: text/plain; charset=UTF-8"
	msg := "To: " + to + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		contentType + "\r\n\r\n" +
		body
	return smtp.SendMail(host, auth, from, strings.Split(to, ";"), []byte(msg))
}

func MailHtml(from, password, host, to, subject, html string) error {
	domain := strings.Split(host, ":")[0]
	auth := smtp.PlainAuth("", from, password, domain)
	contentType := "Content-Type: text/html; charset=UTF-8"
	msg := "To: " + to + "\r\n" +
		"From: " + from + ">\r\n" +
		"Subject: " + subject + "\r\n" +
		contentType + "\r\n\r\n" +
		html
	return smtp.SendMail(host, auth, from, strings.Split(to, ";"), []byte(msg))
}
