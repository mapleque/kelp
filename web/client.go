package web

import (
	"bytes"
	"crypto/tls"
	//	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
)

func PostJson(url string, data []byte) ([]byte, error) {
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

func Mail(account, password, host, from, to, subject, content, contentType string) error {
	return mail(account, password, host, from, to, subject, content, contentType, nil)
}

func MailHtml(from, password, host, to, subject, html string) error {
	contentType := "Content-Type: text/html; charset=UTF-8"
	return Mail(from, password, host, from, to, subject, html, contentType)
}

type MailAttachment struct {
	Filename string
	Mimetype string
	Content  []byte
}

func MailAttachments(
	account, password, host, from, to, subject, content, contentType string, attachements []*MailAttachment) error {
	return mail(account, password, host, from, to, subject, content, contentType, attachements)
}

func mail(account, password, host, from, to, subject, content, contentType string, attachments []*MailAttachment) error {
	msg := "To: " + to + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n"
	boundary := RandMd5()
	boundaryContent := RandMd5()
	if len(attachments) > 0 {
		msg += "Content-Type: multipart/mixed; boundary=" + boundary +
			"\r\n\r\n--" + boundary + "\r\n" +
			"Content-Type: multipart/alternative; boundary=" + boundaryContent +
			"\r\n\r\n--" + boundaryContent + "\r\n"
	}
	msg += contentType + "\r\n\r\n" + content + "\r\n"
	if len(attachments) > 0 {
		msg += "\r\n--" + boundaryContent + "--"
		for _, attach := range attachments {
			msg += "\r\n\r\n--" + boundary + "\r\n" +
				"Content-Type: " + attach.Mimetype + "\r\n" +
				"Content-Disposition: attachment; filename=\"" + attach.Filename + "\"\r\n\r\n"
			//buf := make([]byte, base64.StdEncoding.EncodedLen(len(attach.Content)))
			//base64.StdEncoding.Encode(buf, attach.Content)
			//for i, l := 0, len(buf); i < l; i++ {
			//	msg += string(buf[i])
			//	if (i+1)%76 == 0 {
			//		msg += "\r\n"
			//	}
			//}
			msg += string(attach.Content) + "\r\n"
			msg += "\r\n--" + boundary
		}
		msg += "--"
	}
	domain := strings.Split(host, ":")[0]
	auth := smtp.PlainAuth("", account, password, domain)
	conn, err := tls.Dial("tcp", host, nil)
	if err != nil {
		return err
	}
	client, err := smtp.NewClient(conn, domain)
	if err != nil {
		return err
	}
	if err := client.Auth(auth); err != nil {
		return err
	}
	if err := client.Mail(from); err != nil {
		return err
	}
	for _, addr := range strings.Split(to, ";") {
		err = client.Rcpt(addr)
		if err != nil {
			return err
		}
	}
	writer, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := writer.Write([]byte(msg)); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}
	client.Quit()
	return nil
}
