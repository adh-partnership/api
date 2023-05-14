package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"net/url"
	"os"
	"strings"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database"
	"sigs.k8s.io/yaml"
)

var Templates = map[string]string{
	"visiting_rejected": "visiting_rejected",
	"visiting_added":    "visiting_added",
	"visiting_removed":  "visiting_removed",
	"inactive_warning":  "inactive_warning",
	"inactive":          "inactive",
}

type Template struct {
	Subject string   `json:"subject"`
	CC      []string `json:"cc"`
	BCC     []string `json:"bcc"`
	Body    string   `json:"body"`
}

func fetchRole(role string) []string {
	var ret []string
	users, err := database.FindUsersWithRole(role)
	if err != nil {
		return ret
	}
	for _, user := range users {
		ret = append(ret, fmt.Sprintf("%s %s, %s", user.FirstName, user.LastName, strings.ToUpper(role)))
	}
	return ret
}

func GetTemplate(name string) (*Template, error) {
	if _, err := os.Stat(config.Cfg.Email.TemplateDir + "/" + name + ".tmpl"); err != nil {
		return nil, err
	}

	f, err := os.Open(config.Cfg.Email.TemplateDir + "/" + name + ".tmpl")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
	if err != nil {
		return nil, err
	}

	templ := &Template{}
	err = yaml.Unmarshal(buf.Bytes(), &templ)
	if err != nil {
		return nil, err
	}

	return templ, nil
}

func BuildBody(name string, data map[string]interface{}) (*bytes.Buffer, string, string, string, error) {
	templ, err := GetTemplate(name)
	if err != nil {
		return nil, "", "", "", err
	}

	t, err := template.New(name).Funcs(template.FuncMap{
		"urlEscape": url.QueryEscape,
		"fetchRole": fetchRole,
		"findRole":  fetchRole,
	}).Parse(templ.Body)
	if err != nil {
		return nil, "", "", "", err
	}

	out := new(bytes.Buffer)
	err = t.Execute(out, data)
	if err != nil {
		return nil, "", "", "", err
	}

	return out, templ.Subject, strings.Join(templ.CC, ", "), strings.Join(templ.BCC, ", "), nil
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

func Send(to, from, subject string, template string, data map[string]interface{}) error {
	body, subj, cc, bcc, err := BuildBody(template, data)
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
		tolist = append(tolist, strings.Split(bcc, ", ")...)
	}

	auth := smtp.PlainAuth("", config.Cfg.Email.User, config.Cfg.Email.Password, config.Cfg.Email.Host)
	return smtp.SendMail(config.Cfg.Email.Host+":"+config.Cfg.Email.Port, auth, from, tolist, msg)
}
