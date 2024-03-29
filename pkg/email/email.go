/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/Shopify/gomail"
	"sigs.k8s.io/yaml"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/logger"
)

var Templates = map[string]string{
	"visiting_rejected": "visiting_rejected",
	"visiting_added":    "visiting_added",
	"visiting_removed":  "visiting_removed",
	"inactive_warning":  "inactive_warning",
	"inactive":          "inactive",
}

var log = logger.Logger.WithField("component", "email")

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
		return nil, fmt.Errorf("error stating template: %s", err)
	}

	f, err := os.Open(config.Cfg.Email.TemplateDir + "/" + name + ".tmpl")
	if err != nil {
		return nil, fmt.Errorf("error opening template: %s", err)
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
		return nil, fmt.Errorf("error unmarshalling template: %s", err)
	}

	return templ, nil
}

func BuildBody(name string, data map[string]interface{}) (*bytes.Buffer, string, string, string, error) {
	templ, err := GetTemplate(name)
	if err != nil {
		return nil, "", "", "", fmt.Errorf("error getting template: %s", err)
	}

	t, err := template.New(name).Funcs(template.FuncMap{
		"urlEscape": url.QueryEscape,
		"fetchRole": fetchRole,
		"findRole":  fetchRole,
	}).Parse(templ.Body)
	if err != nil {
		return nil, "", "", "", fmt.Errorf("error parsing template: %s", err)
	}

	out := new(bytes.Buffer)
	err = t.Execute(out, data)
	if err != nil {
		return nil, "", "", "", fmt.Errorf("error executing template: %s", err)
	}

	return out, templ.Subject, strings.Join(templ.CC, ", "), strings.Join(templ.BCC, ", "), nil
}

func Send(to, from, subject string, template string, data map[string]interface{}) error {
	body, subj, cc, bcc, err := BuildBody(template, data)
	if err != nil {
		return fmt.Errorf("error building email body: %s", err)
	}

	log.Infof("template built, subj=%s, cc=%s, bcc=%s", subj, cc, bcc)

	if from == "" {
		from = config.Cfg.Email.From
	}

	if subject == "" {
		subject = subj
	}

	i, err := strconv.Atoi(config.Cfg.Email.Port)
	if err != nil {
		return err
	}
	d := gomail.NewDialer(config.Cfg.Email.Host, i, config.Cfg.Email.User, config.Cfg.Email.Password)
	d.StartTLSPolicy = gomail.OpportunisticStartTLS

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", strings.Split(to, ", ")...)
	if cc != "" {
		m.SetHeader("Cc", strings.Split(cc, ", ")...)
	}
	if bcc != "" {
		m.SetHeader("Bcc", strings.Split(bcc, ", ")...)
	}
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	if err := d.DialAndSend(m); err != nil {
		log.Errorf("Failed to send email: %s", err)
		return err
	}

	return nil
}
