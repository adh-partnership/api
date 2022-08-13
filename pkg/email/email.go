package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"net/url"
	"strings"

	"github.com/kzdv/api/pkg/config"
	"github.com/kzdv/api/pkg/database"
)

func BuildBody(name string, data map[string]interface{}) (*bytes.Buffer, error) {
	templ, err := database.FindEmailTemplate(name)
	if err != nil {
		return nil, err
	}

	t, err := template.New(name).Funcs(template.FuncMap{
		"urlEscape": url.QueryEscape,
	}).Parse(templ.Body)
	if err != nil {
		return nil, err
	}

	out := new(bytes.Buffer)
	err = t.Execute(out, data)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func BuildEmail(from, to, subject string, cc []string, body *bytes.Buffer) []byte {
	var msg string
	msg = "To: " + to + "\r\n"
	msg += "From: " + from + "\r\n"
	if len(cc) > 0 {
		msg += "Cc: " + strings.Join(cc, ", ") + "\n"
	}
	msg += "Subject: " + subject + "\r\n"
	msg += fmt.Sprintf(`MIME-Version: 1.0
Content-Type: text/html; charset="UTF-8"
Content-Transfer-Encoding: quoted-printable

%s`, body.String())
	return []byte(msg)
}

func Send(to, from, subject string, cc []string, bcc []string, template string, data map[string]interface{}) error {
	body, err := BuildBody(template, data)
	if err != nil {
		return err
	}

	if from == "" {
		from = config.Cfg.Email.From
	}

	msg := BuildEmail(from, to, subject, cc, body)

	var tolist []string
	tolist = append(tolist, to)
	if len(cc) > 0 {
		tolist = append(tolist, cc...)
	}
	if len(bcc) > 0 {
		tolist = append(tolist, bcc...)
	}

	auth := smtp.PlainAuth("", config.Cfg.Email.User, config.Cfg.Email.Password, config.Cfg.Email.Host)
	return smtp.SendMail(config.Cfg.Email.Host+":"+config.Cfg.Email.Port, auth, from, tolist, msg)
}
