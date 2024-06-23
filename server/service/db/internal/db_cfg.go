package internal

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	tagNameForDsn = "dsnTag"
)

// = export ConnectionCfg

type ConnectionCfg struct {
	Ip           string
	Port         string
	UserName     string
	Password     string
	DatabaseName string
	ExtraConncetionCfg
}

func (c *ConnectionCfg) withOptions(opts ...connectionCfgOption) {
	for _, opt := range opts {
		opt.apply(c)
	}
}

func (c *ConnectionCfg) getBaseDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.UserName, c.Password, c.Ip, c.Port, c.DatabaseName)
}

func (c *ConnectionCfg) getExtraDsn() string {
	return (*c).ExtraConncetionCfg.getExtraDsn()
}

func (c *ConnectionCfg) GetDSN() string {
	dsn := c.getBaseDsn()
	extraDsn := c.getExtraDsn()
	if len(extraDsn) > 0 {
		dsn = strings.Join([]string{dsn, extraDsn}, "?")
	}
	return dsn
}

func parseToSimpleDsn(field reflect.StructField, fieldValue reflect.Value) string {
	tagName := field.Tag.Get(tagNameForDsn)
	if tagName != "" && fieldValue.String() != "" {
		return fmt.Sprintf("%s=%s", tagName, fieldValue.String())
	}
	return ""
}

// = ExtraConncetionCfg
type ExtraConncetionCfg struct {
	Charset      string `dsnTag:"charset"`
	ParseTime    string `dsnTag:"parseTime"`
	Loc          string `dsnTag:"loc"`
	Timeout      string `dsnTag:"timeout"`
	ReadTimeout  string `dsnTag:"readTimeout"`
	WriteTimeout string `dsnTag:"writeTimeout"`
}

func (ec *ExtraConncetionCfg) getExtraDsn() string {
	tmpDsnList := []string{}
	cfgType := reflect.TypeOf(*ec)
	cfgValue := reflect.ValueOf(*ec)

	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)
		fieldValue := cfgValue.Field(i)
		// Check if the field is set to its zero value
		if reflect.Zero(fieldValue.Type()).Interface() != fieldValue.Interface() {
			if tmpDsn := parseToSimpleDsn(field, fieldValue); tmpDsn != "" {
				tmpDsnList = append(tmpDsnList, tmpDsn)
			}
		}
	}
	return strings.Join(tmpDsnList, "&")
}

// = connectionCfgOption

type connectionCfgOption interface {
	apply(*ConnectionCfg)
}

type connectionCfgOptionFunc func(*ConnectionCfg)

func (fn connectionCfgOptionFunc) apply(c *ConnectionCfg) {
	fn(c)
}

// = export functions

func NewConnectionCfg(userName, password, databaseName string, opts ...connectionCfgOption) *ConnectionCfg {
	cfg := &ConnectionCfg{
		Ip:           "127.0.0.1",
		Port:         "3306",
		UserName:     userName,
		Password:     password,
		DatabaseName: databaseName,
	}
	cfg.withOptions(opts...)
	return cfg
}

func WithIp(ip string) connectionCfgOption {
	return connectionCfgOptionFunc(func(c *ConnectionCfg) {
		c.Ip = ip
	})
}

func WithPort(port string) connectionCfgOption {
	return connectionCfgOptionFunc(func(c *ConnectionCfg) {
		c.Port = port
	})
}

func WithCharset(charset string) connectionCfgOption {
	return connectionCfgOptionFunc(func(c *ConnectionCfg) {
		c.ExtraConncetionCfg.Charset = charset
	})
}

func WithParseTime(parseTime string) connectionCfgOption {
	return connectionCfgOptionFunc(func(c *ConnectionCfg) {
		c.ExtraConncetionCfg.ParseTime = parseTime
	})
}

func WithLoc(loc string) connectionCfgOption {
	return connectionCfgOptionFunc(func(c *ConnectionCfg) {
		c.ExtraConncetionCfg.Loc = loc
	})
}

func WithTimeout(timeout string) connectionCfgOption {
	return connectionCfgOptionFunc(func(c *ConnectionCfg) {
		c.ExtraConncetionCfg.Timeout = timeout
	})
}

func WithReadTimeout(readTimeout string) connectionCfgOption {
	return connectionCfgOptionFunc(func(c *ConnectionCfg) {
		c.ExtraConncetionCfg.ReadTimeout = readTimeout
	})
}

func WithWriteTimeout(writeTimeout string) connectionCfgOption {
	return connectionCfgOptionFunc(func(c *ConnectionCfg) {
		c.ExtraConncetionCfg.WriteTimeout = writeTimeout
	})
}
