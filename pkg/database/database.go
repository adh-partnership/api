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

package database

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"dario.cat/mergo"
	gomysql "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type DBOptions struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Options  string

	MaxOpenConns int
	MaxIdleConns int

	CACert string
	Logger *logrus.Logger
}

var defaultOptions = DBOptions{
	MaxOpenConns: 50,
	MaxIdleConns: 10,
}

func GenerateDSN(options DBOptions) (string, error) {
	var dsn string

	if options.Driver == "mysql" {
		tls := ""
		if options.CACert != "" {
			tls = "&tls=custom"
		}
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true%s", options.User, options.Password,
			options.Host, options.Port, options.Database, tls)
		log.Debugf("dsn=%s", dsn)
		if options.Options != "" {
			dsn += "?" + options.Options
		}
	} else {
		return "", fmt.Errorf("unsupported driver: %s", options.Driver)
	}

	return dsn, nil
}

func HandleCACert(driver string, cacert string) error {
	rootCertPool := x509.NewCertPool()
	pem, err := base64.StdEncoding.DecodeString(cacert)
	if err != nil {
		return err
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return fmt.Errorf("failed to append PEM")
	}

	// @TODO: support other drivers
	if driver == "mysql" {
		err := gomysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs: rootCertPool,
		})
		if err != nil {
			return errors.New("error registering tls config: " + err.Error())
		}
		log.Debugf("registered tls config: %+v", rootCertPool)
	}

	return nil
}

func isValidDriver(driver string) bool {
	return driver == "mysql"
}

func Connect(options DBOptions) error {
	if !isValidDriver(options.Driver) {
		return errors.New("invalid driver: " + options.Driver)
	}

	err := mergo.Merge(&options, defaultOptions)
	if err != nil {
		return errors.New("failed to apply defaults: " + err.Error())
	}

	if options.Logger == nil {
		options.Logger = logrus.New()
	}

	if options.CACert != "" {
		err := HandleCACert(options.Driver, options.CACert)
		if err != nil {
			return err
		}
	}

	dsn, err := GenerateDSN(options)
	if err != nil {
		return err
	}

	var conn *sql.DB
	if options.Driver == "mysql" {
		conn, err = sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		DB, err = gorm.Open(mysql.New(mysql.Config{Conn: conn}), &gorm.Config{Logger: NewLogger(options.Logger)})
		if err != nil {
			return err
		}
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(options.MaxOpenConns)
	sqlDB.SetMaxIdleConns(options.MaxIdleConns)
	sqlDB.SetConnMaxIdleTime(time.Minute * 5)

	return nil
}
