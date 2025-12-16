package test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	sqlitecloud "github.com/sqlitecloud/sqlitecloud-go"
	"github.com/stretchr/testify/assert"
)

func TestParseConnectionString(t *testing.T) {
	connectionString := "sqlitecloud://myproject.sqlite.cloud:8860/mydatabase?timeout=11&compress=yes"

	expectedConfig := &sqlitecloud.SQCloudConfig{
		Host:                  "myproject.sqlite.cloud",
		Port:                  8860,
		ProjectID:             "myproject",
		Username:              "",
		Password:              "",
		Database:              "mydatabase",
		Timeout:               time.Duration(11) * time.Second,
		Compression:           true,
		CompressMode:          sqlitecloud.CompressModeLZ4,
		Secure:                true,
		TlsInsecureSkipVerify: false,
		Pem:                   "",
		ApiKey:                "",
		NoBlob:                false,
		MaxData:               0,
		MaxRows:               0,
		MaxRowset:             0,
	}

	config, err := sqlitecloud.ParseConnectionString(connectionString)
	if err != nil {
		t.Fatalf("Failed to parse connection string: %v", err)
	}

	if !reflect.DeepEqual(config, expectedConfig) {
		t.Fatalf("Parsed config does not match expected config.\nExpected: %+v\nGot: %+v", expectedConfig, config)
	}
}

func TestParseConnectionStringWithAPIKey(t *testing.T) {
	connectionString := "sqlitecloud://user:pass@host.com:8860/dbname?apikey=abc123&timeout=11&compress=true"
	expectedConfig := &sqlitecloud.SQCloudConfig{
		Host:                  "host.com",
		Port:                  8860,
		Username:              "user",
		Password:              "pass",
		Database:              "dbname",
		Timeout:               time.Duration(11) * time.Second,
		Compression:           true,
		CompressMode:          sqlitecloud.CompressModeLZ4,
		Secure:                true,
		TlsInsecureSkipVerify: false,
		Pem:                   "",
		ApiKey:                "abc123",
		NoBlob:                false,
		MaxData:               0,
		MaxRows:               0,
		MaxRowset:             0,
	}

	config, err := sqlitecloud.ParseConnectionString(connectionString)
	assert.NoError(t, err)
	assert.Truef(t, reflect.DeepEqual(expectedConfig, config), "Expected: %+v\nGot: %+v", expectedConfig, config)
}

func TestParseConnectionStringWithToken(t *testing.T) {
	connectionString := "sqlitecloud://host.com:8860/dbname?token=123|tok123&compress=true"
	expectedConfig := &sqlitecloud.SQCloudConfig{
		Host:                  "host.com",
		Port:                  8860,
		Username:              "",
		Password:              "",
		Database:              "dbname",
		Timeout:               0,
		Compression:           true,
		CompressMode:          sqlitecloud.CompressModeLZ4,
		Secure:                true,
		TlsInsecureSkipVerify: false,
		Pem:                   "",
		Token:                 "123|tok123",
		NoBlob:                false,
		MaxData:               0,
		MaxRows:               0,
		MaxRowset:             0,
	}

	config, err := sqlitecloud.ParseConnectionString(connectionString)
	assert.NoError(t, err)
	assert.Truef(t, reflect.DeepEqual(expectedConfig, config), "Expected: %+v\nGot: %+v", expectedConfig, config)
}

func TestParseConnectionStringWithCredentials(t *testing.T) {
	connectionString := "sqlitecloud://user:pass@host.com:8860"
	config, err := sqlitecloud.ParseConnectionString(connectionString)

	assert.NoError(t, err)
	assert.Equal(t, "user", config.Username)
	assert.Equal(t, "pass", config.Password)
	assert.Equal(t, "host.com", config.Host)
	assert.Equal(t, 8860, config.Port)
	assert.Empty(t, config.Database)
}

func TestParseConnectionStringWithoutCredentials(t *testing.T) {
	connectionString := "sqlitecloud://host.com"
	config, err := sqlitecloud.ParseConnectionString(connectionString)

	assert.NoError(t, err)
	assert.Empty(t, config.Username)
	assert.Empty(t, config.Password)
	assert.Equal(t, "host.com", config.Host)
}

func TestParseConnectionStringWithParameters(t *testing.T) {
	tests := []struct {
		param         string
		configParam   string
		value         string
		expectedValue any
	}{
		{
			param:         "timeout",
			configParam:   "Timeout",
			value:         "11",
			expectedValue: time.Duration(11) * time.Second,
		},
		{
			param:         "compression",
			configParam:   "Compression",
			value:         "false",
			expectedValue: false,
		},
		{
			param:         "zerotext",
			configParam:   "Zerotext",
			value:         "true",
			expectedValue: true,
		},
		{
			param:         "memory",
			configParam:   "Memory",
			value:         "true",
			expectedValue: true,
		},
		{
			param:         "create",
			configParam:   "Create",
			value:         "true",
			expectedValue: true,
		},
		{
			param:         "insecure",
			configParam:   "Secure",
			value:         "true",
			expectedValue: false,
		},
		{
			param:         "secure",
			configParam:   "Secure",
			value:         "false",
			expectedValue: false,
		},
		{
			param:         "non_linearizable",
			configParam:   "NonLinearizable",
			value:         "true",
			expectedValue: true,
		},
		{
			param:         "nonlinearizable",
			configParam:   "NonLinearizable",
			value:         "true",
			expectedValue: true,
		},
		{
			param:         "no_verify_certificate",
			configParam:   "TlsInsecureSkipVerify",
			value:         "true",
			expectedValue: true,
		},
		{
			param:         "noblob",
			configParam:   "NoBlob",
			value:         "true",
			expectedValue: true,
		},
		{
			param:         "maxdata",
			configParam:   "MaxData",
			value:         "10",
			expectedValue: 10,
		},
		{
			param:         "maxrows",
			configParam:   "MaxRows",
			value:         "11",
			expectedValue: 11,
		},
		{
			param:         "maxrowset",
			configParam:   "MaxRowset",
			value:         "12",
			expectedValue: 12,
		},
		{
			param:         "apikey",
			configParam:   "ApiKey",
			value:         "abc123",
			expectedValue: "abc123",
		},
		{
			param:         "token",
			configParam:   "Token",
			value:         "123|tok123",
			expectedValue: "123|tok123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.param, func(t *testing.T) {
			config, err := sqlitecloud.ParseConnectionString("sqlitecloud://myhost.sqlite.cloud/mydatabase?" + tt.param + "=" + tt.value)

			assert.NoError(t, err)

			actualValue := reflect.ValueOf(*config).FieldByName(tt.configParam)
			if !actualValue.IsValid() {
				t.Fatalf("Field %s not found in config", tt.configParam)
			} else {
				assert.Equal(t, tt.expectedValue, actualValue.Interface())
			}
		})
	}
}

func TestParseConnectionStringCompressionCombinations(t *testing.T) {
	tests := []struct {
		param         string
		configParam   string
		value         string
		expectedValue any
	}{
		{
			param:         "compression",
			configParam:   "Compression",
			value:         "false",
			expectedValue: false,
		},
		{
			param:         "compression",
			configParam:   "Compression",
			value:         "disabled",
			expectedValue: false,
		},
		{
			param:       "compression",
			configParam: "Compression",
			// value not supported
			value:         "no",
			expectedValue: true,
		},
		{
			param:         "compression",
			configParam:   "Compression",
			value:         "0",
			expectedValue: false,
		},
		{
			param:         "compression",
			configParam:   "Compression",
			value:         "true",
			expectedValue: true,
		},
		{
			param:         "compression",
			configParam:   "Compression",
			value:         "enabled",
			expectedValue: true,
		},
		{
			param:         "compression",
			configParam:   "Compression",
			value:         "yes",
			expectedValue: true,
		},
		{
			param:         "compression",
			configParam:   "Compression",
			value:         "1",
			expectedValue: true,
		},
		{
			param:         "compress",
			configParam:   "CompressMode",
			value:         "lz4",
			expectedValue: sqlitecloud.CompressModeLZ4,
		},
		{
			param:         "compress",
			configParam:   "Compression",
			value:         "lz4",
			expectedValue: true,
		},
		{
			param:         "compress",
			configParam:   "CompressMode",
			value:         "no",
			expectedValue: sqlitecloud.CompressModeNo,
		},
		{
			param:         "compress",
			configParam:   "Compression",
			value:         "no",
			expectedValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.param, func(t *testing.T) {
			config, err := sqlitecloud.ParseConnectionString("sqlitecloud://myhost.sqlite.cloud/mydatabase?" + tt.param + "=" + tt.value)

			assert.NoError(t, err)

			actualValue := reflect.ValueOf(*config).FieldByName(tt.configParam)
			if !actualValue.IsValid() {
				t.Fatalf("Field %s not found in config", tt.configParam)
			} else {
				assert.Equal(t, tt.expectedValue, actualValue.Interface())
			}
		})
	}
}

func TestParseConnectionStringWithTLSParameter(t *testing.T) {
	tests := []struct {
		tlsValue         string
		expectedSecure   bool
		expectedInsecure bool
		expectedPem      string
	}{
		{
			tlsValue:         "true",
			expectedSecure:   true,
			expectedInsecure: false,
			expectedPem:      "",
		},
		{
			tlsValue:         "false",
			expectedSecure:   false,
			expectedInsecure: false,
			expectedPem:      "",
		},
		{
			tlsValue:         "y",
			expectedSecure:   true,
			expectedInsecure: false,
			expectedPem:      "",
		},
		{
			tlsValue:         "n",
			expectedSecure:   false,
			expectedInsecure: false,
			expectedPem:      "",
		},
		{
			tlsValue:         "yes",
			expectedSecure:   true,
			expectedInsecure: false,
			expectedPem:      "",
		},
		{
			tlsValue:         "no",
			expectedSecure:   false,
			expectedInsecure: false,
			expectedPem:      "",
		},
		{
			tlsValue:         "true",
			expectedSecure:   true,
			expectedInsecure: false,
			expectedPem:      "",
		},
		{
			tlsValue:         "on",
			expectedSecure:   true,
			expectedInsecure: false,
			expectedPem:      "",
		},
		{
			tlsValue:         "off",
			expectedSecure:   false,
			expectedInsecure: false,
			expectedPem:      "",
		},
		{
			tlsValue:         "enable",
			expectedSecure:   true,
			expectedInsecure: false,
			expectedPem:      "",
		},
		{
			tlsValue:         "disable",
			expectedSecure:   false,
			expectedInsecure: false,
			expectedPem:      "",
		},
		{
			tlsValue:         "enabled",
			expectedSecure:   true,
			expectedInsecure: false,
			expectedPem:      "",
		},
		{
			tlsValue:         "disabled",
			expectedSecure:   false,
			expectedInsecure: false,
			expectedPem:      "",
		},
		{
			tlsValue:         "skip",
			expectedSecure:   true,
			expectedInsecure: true,
			expectedPem:      "",
		},
		{
			tlsValue:         "intern",
			expectedSecure:   true,
			expectedInsecure: false,
			expectedPem:      sqlitecloud.SQLiteCloudCA,
		},
		{
			tlsValue:         "<use internal pem>",
			expectedSecure:   true,
			expectedInsecure: false,
			expectedPem:      sqlitecloud.SQLiteCloudCA,
		},
		{
			tlsValue:         "anything-default",
			expectedSecure:   true,
			expectedInsecure: false,
			expectedPem:      "anything-default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.tlsValue, func(t *testing.T) {
			connectionString := fmt.Sprintf("sqlitecloud://myhost.sqlite.cloud:8860/mydatabase?tls=%s", tt.tlsValue)

			config, err := sqlitecloud.ParseConnectionString(connectionString)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedSecure, config.Secure)
			assert.Equal(t, tt.expectedInsecure, config.TlsInsecureSkipVerify)
			assert.Equal(t, tt.expectedPem, config.Pem)
		})
	}
}
