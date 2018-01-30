package web

import (
	"testing"
)

const (
	account  = "mapleque@163.com"
	password = "" // client password
	host     = "smtp.163.com:465"
	from     = "mapleque@163.com"
	to       = "mapleque@163.com"
)

func TestMail(t *testing.T) {
	if len(password) == 0 {
		return
	}
	if err := MailHtml(
		account, password, host, to,
		"测试文本",
		"<h1>测试正文</h1>",
	); err != nil {
		t.Error(err)
	}
}

func TestMailHtml(t *testing.T) {
	if len(password) == 0 {
		return
	}
	if err := MailHtml(
		account, password, host, to,
		"测试Html",
		"<h1>测试Html</h1>",
	); err != nil {
		t.Error(err)
	}
}

func TestMailAttachments(t *testing.T) {
	if len(password) == 0 {
		return
	}
	attachments := []*MailAttachment{
		&MailAttachment{
			"attach.txt",
			"application/octet-stream",
			[]byte("这里是附件文本\nattachment words here!"),
		},
	}

	if err := MailAttachments(
		account, password, host, from, to,
		"测试附件",
		"<h1>测试附件正文</h1>",
		"Content-Type: text/html; charset=UTF-8",
		attachments,
	); err != nil {
		t.Error(err)
	}
}
