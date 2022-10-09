package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"net/url"
	"strings"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database"
)

func BuildBody(name string, data map[string]interface{}) (*bytes.Buffer, string, string, error) {
	templ, err := database.FindEmailTemplate(name)
	if err != nil {
		return nil, "", "", err
	}

	t, err := template.New(name).Funcs(template.FuncMap{
		"urlEscape": url.QueryEscape,
		"fetchRole": func(role string) []string {
			var ret []string
			users, err := database.FindUsersWithRole(role)
			if err != nil {
				return ret
			}
			for _, user := range users {
				ret = append(ret, fmt.Sprintf("%s %s, %s", user.FirstName, user.LastName, strings.ToUpper(role)))
			}
			return ret
		},
	}).Parse(templ.Body)
	if err != nil {
		return nil, "", "", err
	}

	out := new(bytes.Buffer)
	err = t.Execute(out, data)
	if err != nil {
		return nil, "", "", err
	}

	return out, templ.Subject, templ.CC, nil
}

func BuildEmail(from, to, subject string, cc string, body *bytes.Buffer) []byte {
	var msg string
	msg = "To: " + to + "\r\n"
	msg += "From: " + from + "\r\n"
	if len(cc) > 0 {
		msg += "Cc: " + cc + "\n"
	}
	msg += "Subject: " + subject + "\r\n"
	msg += fmt.Sprintf(`MIME-Version: 1.0
Content-Type: text/html; charset="UTF-8"
Content-Transfer-Encoding: quoted-printable

%s`, body.String())
	return []byte(msg)
}

func Send(to, from, subject string, bcc []string, template string, data map[string]interface{}) error {
	body, subj, cc, err := BuildBody(template, data)
	if err != nil {
		return err
	}

	if from == "" {
		from = config.Cfg.Email.From
	}

	if subject == "" {
		subject = subj
	}

	msg := BuildEmail(from, to, subject, cc, body)

	var tolist []string
	tolist = append(tolist, to)
	if len(cc) > 0 {
		tolist = append(tolist, strings.Split(cc, ", ")...)
	}
	if len(bcc) > 0 {
		tolist = append(tolist, bcc...)
	}

	auth := smtp.PlainAuth("", config.Cfg.Email.User, config.Cfg.Email.Password, config.Cfg.Email.Host)
	return smtp.SendMail(config.Cfg.Email.Host+":"+config.Cfg.Email.Port, auth, from, tolist, msg)
}
