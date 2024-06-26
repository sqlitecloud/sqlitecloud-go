//
//                    ////              SQLite Cloud
//        ////////////  ///
//      ///             ///  ///        Product     : SQLite Cloud GO SDK
//     ///             ///  ///         Version     : 1.1.2
//     //             ///   ///  ///    Date        : 2021/10/13
//    ///             ///   ///  ///    Author      : Andreas Pfeil
//   ///             ///   ///  ///
//   ///     //////////   ///  ///      Description : Go Methods related to the
//   ////                ///  ///                     SQCloud class for managing
//     ////     //////////   ///                      the connection and executing
//        ////            ////                        queries.
//          ////     /////
//             ///                      Copyright   : 2021 by SQLite Cloud Inc.
//
// -----------------------------------------------------------------------TAB=2

package sqlitecloud

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/xo/dburl"
)

type SQCloudConfig struct {
	Host                  string
	Port                  int
	Username              string
	Password              string
	Database              string
	PasswordHashed        bool          // Password is hashed
	Timeout               time.Duration // Optional query timeout passed directly to TLS socket
	CompressMode          string        // eg: LZ4
	Compression           bool          // Enable compression
	Zerotext              bool          // Tell the server to zero-terminate strings
	Memory                bool          // Database will be created in memory
	Create                bool          // Create the database if it doesn't exist?
	Secure                bool          // Connect using plain TCP port, without TLS encryption, NOT RECOMMENDED (insecure)
	NonLinearizable       bool          // Request for immediate responses from the server node without waiting for linerizability guarantees
	TlsInsecureSkipVerify bool          // Accept invalid TLS certificates (no_verify_certificate)
	Pem                   string
	ApiKey                string
	NoBlob                bool // flag to tell the server to not send BLOB columns
	MaxData               int  // value to tell the server to not send columns with more than max_data bytes
	MaxRows               int  // value to control rowset chunks based on the number of rows
	MaxRowset             int  // value to control the maximum allowed size for a rowset
}

type SQCloud struct {
	SQCloudConfig

	sock *net.Conn

	psub *SQCloud
	// psubc    chan string
	Callback func(*SQCloud, string)

	cert *tls.Config

	uuid   string // 36 runes -> remove maybe????
	secret string // 36 runes -> remove maybe????

	ErrorCode    int
	ExtErrorCode int
	ErrorOffset  int
	ErrorMessage string
}

const SQLiteDefaultPort = 8860

const (
	CompressModeNo  = "NO"
	CompressModeLZ4 = "LZ4"
)

const SQLiteCloudCA = "SQLiteCloudCA"

func New(config SQCloudConfig) *SQCloud {
	return &SQCloud{SQCloudConfig: config}
}

// init registers the sqlitecloud scheme in the connection string parser.
func init() {
	dburl.Register(dburl.Scheme{
		Driver:    "sc", // sqlitecloud
		Generator: dburl.GenFromURL("sqlitecloud://user:pass@host.com:8860/dbname?timeout=10&compress=NO&tls=INTERN"),
		Transport: dburl.TransportTCP,
		Opaque:    false,
		Aliases:   []string{"sqlitecloud"},
		Override:  "",
	})
}

// Helper functions

// ParseConnectionString parses the given connection string and returns it's components.
// An empty string ("") or -1 is returned as the corresponding return value, if a component of the connection string was not present.
// No plausibility checks are done (see: CheckConnectionParameter).
// If the connection string could not be parsed, an error is returned.
func ParseConnectionString(ConnectionString string) (config *SQCloudConfig, err error) {
	u, err := dburl.Parse(ConnectionString)
	if err == nil {
		config = &SQCloudConfig{}

		config.Host = u.Hostname()
		config.Port = SQLiteDefaultPort
		config.Username = u.User.Username()
		config.Password, _ = u.User.Password()
		config.Database = strings.TrimPrefix(u.Path, "/")
		config.Timeout = 0
		config.Compression = true // enabled by default
		config.CompressMode = CompressModeLZ4
		config.Zerotext = false
		config.Memory = false
		config.Create = false
		config.Secure = true
		config.NonLinearizable = false
		config.TlsInsecureSkipVerify = false
		config.Pem = ""
		config.NoBlob = false
		config.MaxData = 0
		config.MaxRows = 0
		config.MaxRowset = 0
		config.ApiKey = ""

		sPort := strings.TrimSpace(u.Port())
		if len(sPort) > 0 {
			if config.Port, err = strconv.Atoi(sPort); err != nil {
				return nil, err
			}
		}

		for key, values := range u.Query() {
			lastLiteral := strings.TrimSpace(values[len(values)-1])
			switch strings.ToLower(strings.TrimSpace(key)) {
			case "timeout":
				if timeout, err := strconv.Atoi(lastLiteral); err != nil {
					return nil, err
				} else {
					config.Timeout = time.Duration(timeout) * time.Second
				}

			case "compress":
				mode := strings.ToUpper(lastLiteral)
				if mode == CompressModeNo {
					config.CompressMode = CompressModeNo
					config.Compression = false
				} else if enable, err := parseBool(mode, config.Compression); err == nil && !enable {
					config.CompressMode = CompressModeNo
					config.Compression = false
				}
			case "compression":
				if enable, err := parseBool(lastLiteral, config.Compression); err == nil && !enable {
					config.Compression = false
					config.CompressMode = CompressModeNo
				}
			case "zerotext":
				if b, err := parseBool(lastLiteral, config.Zerotext); err == nil {
					config.Zerotext = b
				}
			case "memory":
				if b, err := parseBool(lastLiteral, config.Memory); err == nil {
					config.Memory = b
				}
			case "create":
				if b, err := parseBool(lastLiteral, config.Create); err == nil {
					config.Create = b
				}
			case "secure":
				if b, err := parseBool(lastLiteral, config.Secure); err == nil {
					config.Secure = b
				}
			case "insecure":
				if b, err := parseBool(lastLiteral, config.Secure); err == nil {
					config.Secure = !b
				}
			case "non_linearizable", "nonlinearizable":
				if b, err := parseBool(lastLiteral, config.NonLinearizable); err == nil {
					config.NonLinearizable = b
				}
			case "no_verify_certificate":
				if b, err := parseBool(lastLiteral, config.TlsInsecureSkipVerify); err == nil {
					config.TlsInsecureSkipVerify = b
				}
			case "tls":
				config.Secure, config.TlsInsecureSkipVerify, config.Pem = ParseTlsString(lastLiteral)
			case "apikey":
				config.ApiKey = lastLiteral
			case "noblob":
				if b, err := parseBool(lastLiteral, config.NoBlob); err == nil {
					config.NoBlob = b
				}
			case "maxdata":
				if v, err := strconv.Atoi(lastLiteral); err == nil {
					config.MaxData = v
				}
			case "maxrows":
				if v, err := strconv.Atoi(lastLiteral); err == nil {
					config.MaxRows = v
				}
			case "maxrowset":
				if v, err := strconv.Atoi(lastLiteral); err == nil {
					config.MaxRowset = v
				}
			}
		}

		// fmt.Printf( "NO ERROR: Host=%s, Port=%d, User=%s, Password=%s, Database=%s, Timeout=%d, Compress=%s\r\n", host, port, user, pass, database, timeout, compress )
		return config, nil
	}
	return nil, err
}

func ParseTlsString(tlsconf string) (secure bool, tlsInsecureSkipVerify bool, pem string) {
	switch strings.ToUpper(strings.TrimSpace(tlsconf)) {
	case "", "0", "N", "NO", "FALSE", "OFF", "DISABLE", "DISABLED":
		return false, false, ""
	case "1", "Y", "YES", "TRUE", "ON", "ENABLE", "ENABLED":
		return true, false, ""
	case "skip", "SKIP":
		return true, true, ""
	case strings.ToUpper(SQLiteCloudCA), "INTERN", "<USE INTERNAL PEM>":
		return true, false, SQLiteCloudCA
	default:
		return true, false, strings.TrimSpace(tlsconf)
	}
}

// CheckConnectionParameter checks the given connection arguments for validly.
// Host is either a resolve able hostname or an IP address.
// Port is an unsigned int number between 1 and 65535.
// Timeout must be 0 (=no timeout) or a positive number.
// Compress must be "NO" or "LZ4".
// Username, Password and Database are ignored.
// If a given value does not fulfill the above criteria's, an error is returned.
func (this *SQCloud) CheckConnectionParameter() error {
	if strings.TrimSpace(this.Host) == "" {
		return fmt.Errorf("Invalid hostname (%s)", this.Host)
	}

	if this.Timeout < 0 {
		return fmt.Errorf("Invalid Timeout (%s)", this.Timeout.String())
	}

	return nil
}

// Creation

// reset resets all Connection attributes.
func (this *SQCloud) reset() {
	_ = this.Close()
	this.uuid = ""
	this.secret = ""
	this.resetError()
}

// Connect creates a new connection and tries to connect to the server using the given connection string.
// The given connection string is parsed and checked for correct parameters.
// Nil and an error is returned if the connection string had invalid values or a connection to the server could not be established,
// otherwise, a pointer to the newly established connection is returned.
func Connect(ConnectionString string) (*SQCloud, error) {
	config, err := ParseConnectionString(ConnectionString)
	if err != nil {
		return nil, err
	}

	connection := &SQCloud{SQCloudConfig: *config}

	if err = connection.Connect(); err != nil {
		_ = connection.Close()
		return nil, err
	}

	return connection, err
}

// Connection Functions

// Connect connects to a SQLite Cloud server instance using the given arguments.
// If Connect is called on an already established connection, the old connection is closed first.
// All arguments are checked for valid values (see: CheckConnectionParameter).
// invalid argument values where given or the connection could not be established.
func (this *SQCloud) Connect() error {
	this.reset() // also closes an open connection

	if err := this.CheckConnectionParameter(); err != nil {
		return err
	}

	return this.reconnect()
}

// reconnect closes and then reopens a connection to the SQLite Cloud database server.
func (this *SQCloud) reconnect() error {
	if this.sock != nil {
		return nil
	}

	this.resetError()

	if this.Secure {
		cert, err := getTlsConfig(&this.SQCloudConfig)
		if err != nil {
			return err
		}
		this.cert = cert
	}

	var dialer = net.Dialer{}
	dialer.Timeout = this.Timeout

	switch {
	case this.cert != nil:
		if tls_c, err := tls.DialWithDialer(&dialer, "tcp", net.JoinHostPort(this.Host, strconv.Itoa(this.Port)), this.cert); err != nil {
			this.ErrorCode = -1
			this.ErrorMessage = err.Error()
			return err
		} else {
			c := net.Conn(tls_c)
			this.sock = &c
		}
	default:
		// todo: use the dialer...
		if c, err := net.DialTimeout("tcp", net.JoinHostPort(this.Host, strconv.Itoa(this.Port)), this.Timeout); err != nil {
			this.ErrorCode = -1
			this.ErrorMessage = err.Error()
			return err
		} else {
			this.sock = &c
		}
	}

	commands, args := connectionCommands(this.SQCloudConfig)

	if commands != "" {
		if len(args) > 0 {
			if err := this.ExecuteArray(commands, args); err != nil {
				return err
			}
		} else {
			if err := this.Execute(commands); err != nil {
				return err
			}
		}
	}

	return nil
}

// Close closes the connection to the SQLite Cloud Database server.
// The connection can later be reopened (see: reconnect)
func (this *SQCloud) Close() error {
	var err_sock, err_psub error

	err_psub = this.psubClose()

	if this.sock != nil {
		err_sock = (*this.sock).Close()
	}
	this.sock = nil

	this.resetError()

	if err_sock != nil {
		this.setError(-1, err_sock.Error())
		return err_sock
	}

	if err_psub != nil {
		this.setError(-1, err_psub.Error())
		return err_psub
	}
	return nil
}

func getTlsConfig(config *SQCloudConfig) (*tls.Config, error) {
	var pool *x509.CertPool = nil
	pem := []byte{}

	switch _, _, trimmed := ParseTlsString(config.Pem); trimmed {
	case "":
		break
	case SQLiteCloudCA:
		pem = []byte(sqliteCloudCAPEM)
	default:
		// check if it is a filepath
		_, err := os.Stat(trimmed)
		if os.IsNotExist(err) {
			// not a filepath, use the string as a pem string
			pem = []byte(trimmed)
		} else {
			// its a file, read its content into the pem string
			switch bytes, err := os.ReadFile(trimmed); {
			case err != nil:
				return nil, fmt.Errorf("could not open PEM file in '%s'", trimmed)
			default:
				pem = bytes
			}
		}
	}

	if len(pem) > 0 {
		pool = x509.NewCertPool()

		if !pool.AppendCertsFromPEM(pem) {
			return nil, fmt.Errorf("could not append certs from PEM")
		}
	}

	return &tls.Config{
		RootCAs:            pool,
		InsecureSkipVerify: config.TlsInsecureSkipVerify,
		MinVersion:         tls.VersionTLS12,
	}, nil
}

func connectionCommands(config SQCloudConfig) (string, []interface{}) {
	buffer := ""
	args := []interface{}{}

	// it must be executed before authentication command
	if config.NonLinearizable {
		buffer += nonlinearizableCommand(config.NonLinearizable)
	}

	if config.ApiKey != "" {
		c, a := authWithKeyCommand(config.ApiKey)
		buffer += c
		args = append(args, a...)
	}

	if config.Username != "" && config.Password != "" {
		c, a := authCommand(config.Username, config.Password, config.PasswordHashed)
		buffer += c
		args = append(args, a...)
	}

	if config.Database != "" {
		create := config.Create && !config.Memory
		c, a := useDatabaseCommand(config.Database, create)
		buffer += c
		args = append(args, a...)
	}

	buffer += compressCommand(config.CompressMode)

	if config.Zerotext {
		buffer += zerotextCommand(config.Zerotext)
	}

	if config.NoBlob {
		buffer += noblobCommand(config.NoBlob)
	}

	if config.MaxData > 0 {
		buffer += maxdataCommand(config.MaxData)
	}

	if config.MaxRows > 0 {
		buffer += maxrowsCommand(config.MaxRows)
	}

	if config.MaxRowset > 0 {
		buffer += maxrowsetCommand(config.MaxRowset)
	}

	return buffer, args
}

func noblobCommand(NoBlob bool) string {
	if NoBlob {
		return "SET CLIENT KEY NOBLOB TO 1;"
	} else {
		return "SET CLIENT KEY NOBLOB TO 0;"
	}
}

func maxdataCommand(v int) string {
	return fmt.Sprintf("SET CLIENT KEY MAXDATA TO %d;", v)
}

func maxrowsCommand(v int) string {
	return fmt.Sprintf("SET CLIENT KEY MAXROWS TO %d;", v)
}

func maxrowsetCommand(v int) string {
	return fmt.Sprintf("SET CLIENT KEY MAXROWSET TO %d;", v)
}

func compressCommand(CompressMode string) string {
	switch compression := strings.ToUpper(CompressMode); {
	case compression == CompressModeNo:
		return "SET CLIENT KEY COMPRESSION TO 0;"
	case compression == CompressModeLZ4:
		return "SET CLIENT KEY COMPRESSION TO 1;"
	default:
		return ""
	}
}

func nonlinearizableCommand(NonLinearizable bool) string {
	if NonLinearizable {
		return "SET CLIENT KEY NONLINEARIZABLE TO 1;"
	} else {
		return "SET CLIENT KEY NONLINEARIZABLE TO 0;"
	}
}

func zerotextCommand(Zerotext bool) string {
	if Zerotext {
		return "SET CLIENT KEY ZEROTEXT TO 1;"
	} else {
		return "SET CLIENT KEY ZEROTEXT TO 0;"
	}
}

// Compress enabled or disables data compression for this connection.
// If enabled, the data is compressed with the LZ4 compression algorithm, otherwise no compression is applied the data.
func (this *SQCloud) Compress(CompressMode string) error {
	switch c := compressCommand(CompressMode); {
	case this.sock == nil:
		return errors.New("Not connected")
	case c == "":
		return fmt.Errorf("Invalid method (%s)", CompressMode)
	default:
		return this.Execute(c)
	}
}

// IsConnected checks the connection to the SQLite Cloud database server by sending a PING command.
// true is returned, if the connection is established and actually working, false otherwise.
func (this *SQCloud) IsConnected() bool {
	switch {
	case this.sock == nil:
		return false
	case this.Ping() != nil:
		return false
	default:
		return true
	}
}

// Error Methods

func (this *SQCloud) setError(ErrorCode int, ErrorMessage string) {
	this.ErrorCode = ErrorCode
	this.ErrorMessage = ErrorMessage
}

// resetError resets the error code and message of the last run command.
func (this *SQCloud) resetError() { this.setError(0, "") }

// GetErrorCode returns the error code of the last unsuccessful command as an int value.
// 0 is returned if the last command run successful.
func (this *SQCloud) GetErrorCode() int { return this.ErrorCode }

// GetExtErrorCode returns the error code of the last unsuccessful command as an int value.
// 0 is returned if the last command run successful.
func (this *SQCloud) GetExtErrorCode() int { return this.ExtErrorCode }

// GetErrorOffset returns the error code of the last unsuccessful command as an int value.
// 0 is returned if the last command run successful.
func (this *SQCloud) GetErrorOffset() int { return this.ErrorOffset }

// IsError checks the successful execution of the last method call / command.
// true is returned if the last command resulted in an error, false otherwise.
func (this *SQCloud) IsError() bool { return this.GetErrorCode() != 0 }

// GetErrorMessage returns the error message of the last unsuccuessful command as an error.
// nil is returned if the last command run successful.
func (this *SQCloud) GetErrorMessage() error {
	switch this.IsError() {
	case true:
		return errors.New(this.ErrorMessage)
	default:
		return nil
	}
}

// GetError returned the error code and message of the last unsuccessful command.
// 0 and nil is returned if the last command run successful.
func (this *SQCloud) GetError() (int, int, int, error) {
	return this.GetErrorCode(), this.GetExtErrorCode(), this.GetErrorOffset(), this.GetErrorMessage()
}

// Data Access Functions

// Select executes a query on an open SQLite Cloud database connection.
// If an error occurs during the execution of the query, nil and an error describing the problem is returned.
// On successful execution, a pointer to the result is returned.
func (this *SQCloud) Select(SQL string) (*Result, error) {
	this.resetError()

	if _, err := this.sendString(SQL); err != nil {
		this.ErrorCode = 100003
		this.ErrorMessage = fmt.Sprintf("Internal Error: sendString (%s)", err.Error())
		return nil, errors.New(this.ErrorMessage)
	}

	switch result, err := this.readResult(); {
	case result == nil:
		return nil, errors.New("nil")

	case result.IsError():
		this.ErrorCode, this.ExtErrorCode, this.ErrorOffset, this.ErrorMessage, _ = result.GetError()
		result.Free()
		return nil, errors.New(this.ErrorMessage)

	case err != nil:
		this.ErrorCode, this.ExtErrorCode, this.ErrorOffset, this.ErrorMessage = 100000, NO_EXTCODE, NO_OFFCODE, err.Error()
		result.Free()
		return nil, err

	default:
		return result, nil
	}
}

func (this *SQCloud) SelectArray(SQL string, values []interface{}) (*Result, error) {
	this.resetError()

	if _, err := this.sendArray(SQL, values); err != nil {
		this.ErrorCode = 100003
		this.ErrorMessage = fmt.Sprintf("Internal Error: sendArray (%s)", err.Error())
		return nil, errors.New(this.ErrorMessage)
	}

	switch result, err := this.readResult(); {
	case result == nil:
		return nil, errors.New("nil")

	case result.IsError():
		this.ErrorCode, this.ExtErrorCode, this.ErrorOffset, this.ErrorMessage, _ = result.GetError()
		result.Free()
		return nil, errors.New(this.ErrorMessage)

	case err != nil:
		this.ErrorCode, this.ExtErrorCode, this.ErrorOffset, this.ErrorMessage = 100000, NO_EXTCODE, NO_OFFCODE, err.Error()
		result.Free()
		return nil, err

	default:
		return result, nil
	}
}

func (this *SQCloud) SendBlob(data []byte) error {
	this.resetError()

	if _, err := this.sendBytes(data); err != nil {
		this.ErrorCode = 100003
		this.ErrorMessage = fmt.Sprintf("Internal Error: sendBytes (%s)", err.Error())
		return errors.New(this.ErrorMessage)
	}

	switch result, err := this.readResult(); {
	case result == nil:
		return errors.New("nil")

	case result.IsError():
		this.ErrorCode, this.ExtErrorCode, this.ErrorOffset, this.ErrorMessage, _ = result.GetError()
		result.Free()
		return errors.New(this.ErrorMessage)

	case err != nil:
		this.ErrorCode, this.ExtErrorCode, this.ErrorOffset, this.ErrorMessage = 100000, NO_EXTCODE, NO_OFFCODE, err.Error()
		result.Free()
		return err

	default:
		return nil
	}
}

// Execute executes the given query.
// If the execution was not successful, an error describing the reason of the failure is returned.
func (this *SQCloud) Execute(SQL string) error {
	if result, err := this.Select(SQL); result != nil {
		defer result.Free()

		if !result.IsOK() && !result.IsArray() {
			return errors.New("ERROR: Unexpected Result (-1)")
		}
		return err
	} else {
		return err
	}
}

func (this *SQCloud) ExecuteArray(SQL string, values []interface{}) error {
	if result, err := this.SelectArray(SQL, values); result != nil {
		defer result.Free()

		if !result.IsOK() && !result.IsArray() {
			return errors.New("ERROR: Unexpected Result (-1)")
		}
		return err
	} else {
		return err
	}
}
