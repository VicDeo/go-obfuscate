package config

import (
	"fmt"
	"net"
	"path"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"github.com/vicdeo/go-obfuscate/faker"
)

type (
	// DatabaseConfig -- Database connection config
	DatabaseConfig struct {
		Net          string `yaml:"net,omitempty"`
		Socket       string `yaml:"socket,omitempty"`
		Hostname     string `yaml:"hostname,omitempty"`
		Port         string `yaml:"port,omitempty"`
		DatabaseName string `yaml:"databaseName,omitempty"`
		User         string `yaml:"user,omitempty"`
		Password     string `yaml:"password,omitempty"`
	}

	// OutputConfig -- dump-specific options
	OutputConfig struct {
		FileNameFormat string `yaml:"fileNameFormat"`
		Directory      string `yaml:"directory"`
	}

	// Config - global config
	Config struct {
		Database *DatabaseConfig        `yaml:"database"`
		Output   *OutputConfig          `yaml:"output"`
		Tables   map[string]interface{} `yaml:"tables"`
		clock    func() time.Time
	}
)

const (
	ignoreMarker   = "ignore"
	truncateMarker = "truncate"
)

// Create a new Config instance.
var (
	conf         *Config
	dumpFileName string
)

// GetConf - Read the config file and marshal into the conf Config struct.
func GetConf(configDir, configFileName string) (*Config, error) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	viper.AddConfigPath(configDir)
	viper.SetConfigName(configFileName)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	conf = &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

// IsIgnoredTable
func IsIgnoredTable(tableName string) bool {
	if val, ok := conf.Tables[tableName]; ok {
		return val == ignoreMarker
	}
	return false
}

// ShouldDumpData
func ShouldDumpData(tableName string) bool {
	if val, ok := conf.Tables[tableName]; ok {
		return val != truncateMarker && !IsIgnoredTable(tableName)
	}
	return true
}

// GetColumnFaker - get a proper data generator
func GetColumnFaker(tableName, columnName string) faker.FakeGenerator {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	if table, ok := conf.Tables[tableName]; ok {
		tableMap := table.(map[string]interface{})
		if column, ok := tableMap[columnName]; ok {
			columnMap := column.(map[string]interface{})
			return faker.New(columnMap)
		}
	}
	return nil
}

func (config *Config) GetDumpFullPath() string {
	return path.Join(config.Output.Directory, config.GetDumpFileName())
}

func (config *Config) GetDumpFileName() string {
	if dumpFileName == "" {
		// Uses time.Time.Format (https://golang.org/pkg/time/#Time.Format). format appended with '.sql'.
		dumpFileName = config.now().Format(config.Output.FileNameFormat)
		dumpFileName = fmt.Sprintf(dumpFileName, config.Database.DatabaseName) + ".sql"
	}
	return dumpFileName
}

func (config *Config) now() time.Time {
	if config.clock == nil {
		return time.Now()
	}
	return config.clock()
}

func (config *DatabaseConfig) GetMysqlConfigDSN() string {
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.DBName = config.DatabaseName
	mysqlConfig.Net = config.Net
	mysqlConfig.User = config.User
	mysqlConfig.Passwd = config.Password

	switch config.Net {
	case "tcp":
		mysqlConfig.Addr = net.JoinHostPort(config.Hostname, config.Port)
	case "unix":
		mysqlConfig.Addr = config.Socket
	}
	return mysqlConfig.FormatDSN()
}
