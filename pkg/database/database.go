package database

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/imdario/mergo"
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

var default_options = DBOptions{
	MaxOpenConns: 50,
	MaxIdleConns: 10,
}

func GenerateDSN(options DBOptions) (string, error) {
	var dsn string

	if options.Driver == "mysql" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", options.User, options.Password,
			options.Host, options.Port, options.Database)

		if options.Options != "" {
			dsn += "?" + options.Options
		}
	} else {
		return "", fmt.Errorf("unsupported driver: %s", options.Driver)
	}

	return dsn, nil
}

func HandleCACert(driver string, CACert string) error {
	rootCertPool := x509.NewCertPool()
	pem, err := base64.StdEncoding.DecodeString(CACert)
	if err != nil {
		return err
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return fmt.Errorf("failed to append PEM")
	}

	// @TODO: support other drivers
	if driver == "mysql" {
		gomysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs: rootCertPool,
		})
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

	err := mergo.Merge(&options, default_options)
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

	var conn gorm.Dialector
	if options.Driver == "mysql" {
		conn = mysql.Open(dsn)
	}

	DB, err = gorm.Open(conn, &gorm.Config{Logger: NewLogger(options.Logger)})
	if err != nil {
		return err
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
