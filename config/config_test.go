package config

import (
	"reflect"

	"testing"
	"time"
)

type validatePair struct {
	expectedMessages  map[string][]string
	expectedHasErrors bool
	config            Config
}

var validateTescases = []validatePair{
	{
		map[string][]string{
			"obfuscate": []string{},
			"keep":      []string{},
			"ignore":    []string{},
			"truncate":  []string{},
			"overall":   []string{},
		},
		false,
		Config{
			Tables: &TableConfig{
				Keep:      []string{},
				Ignore:    []string{},
				Truncate:  []string{},
				Obfuscate: map[string]interface{}{},
			},
		},
	},
	{
		map[string][]string{
			"obfuscate": []string{},
			"keep":      []string{},
			"ignore":    []string{"d"},
			"truncate":  []string{},
			"overall":   []string{"f", "a"},
		},
		true,
		Config{
			Tables: &TableConfig{
				Keep:      []string{"a", "b", "c"},
				Ignore:    []string{"d", "e", "f", "d"},
				Truncate:  []string{"a"},
				Obfuscate: map[string]interface{}{"f": "b"},
			},
		},
	},
}

func TestValidateConfig(t *testing.T) {
	for _, testcase := range validateTescases {
		actualMessages, hasErrors := testcase.config.ValidateConfig()
		if !reflect.DeepEqual(actualMessages, testcase.expectedMessages) {
			t.Error(
				"Got different messages list:",
				"expected", testcase.expectedMessages,
				"got", actualMessages,
			)
		}
		if hasErrors != testcase.expectedHasErrors {
			t.Error(
				"Got different error status:",
				"expected", testcase.expectedHasErrors,
				"got", hasErrors,
			)
		}
	}
}

type getDumpFileNamePair struct {
	expectedFileName string
	config           Config
}

var getDumpFileNameTestcases = []getDumpFileNamePair{
	{
		"black_mamba-2022-06-01T010203.sql",
		Config{
			Output:   &OutputConfig{FileNameFormat: "%s-2006-01-02T150405"},
			Database: &DatabaseConfig{DatabaseName: "black_mamba"},
		},
	},
}

func TestGetDumpFileName(t *testing.T) {
	for _, testcase := range getDumpFileNameTestcases {

		testcase.config.clock = func() time.Time { return time.Date(2022, 06, 01, 01, 02, 03, 0, time.UTC) }
		fileName := testcase.config.GetDumpFileName()
		if testcase.expectedFileName != fileName {
			t.Error("Expected name is", testcase.expectedFileName, ", got", fileName)
		}
	}
}

type getMysqlConfigDSNPair struct {
	expectedDSN string
	config      DatabaseConfig
}

var getMysqlConfigDSNTestcases = []getMysqlConfigDSNPair{
	{
		"dbuser:dbpass@tcp(127.0.0.1:3306)/black_mamba",
		DatabaseConfig{Net: "tcp", DatabaseName: "black_mamba", User: "dbuser", Password: "dbpass", Hostname: "127.0.0.1", Port: "3306"},
	},
	{
		"unix(/tmp/mysql.sock)/black_mamba",
		DatabaseConfig{Net: "unix", DatabaseName: "black_mamba", Socket: "/tmp/mysql.sock"},
	},
	{
		"dbuser:dbpass@unix(/tmp/mysql.sock)/black_mamba",
		DatabaseConfig{Net: "unix", DatabaseName: "black_mamba", Socket: "/tmp/mysql.sock", User: "dbuser", Password: "dbpass"},
	},
}

func TestGetMysqlConfigDSN(t *testing.T) {
	for _, testcase := range getMysqlConfigDSNTestcases {
		mysqlDSN := testcase.config.GetMysqlConfigDSN()
		if testcase.expectedDSN != mysqlDSN {
			t.Error("Expected ", testcase.expectedDSN, " got ", mysqlDSN)
		}
	}
}
