package database

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-sql-driver/mysql"
	_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewFromEnv sets up the database connections using the configuration in the
// process's environment variables. This should be called just once per server
// instance.
func NewFromEnv(ctx context.Context, cfg *Config) (*DB, error) {
	db, err := gorm.Open(_mysql.Open(dbToMysqlDSN(cfg)), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %v", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDb.SetMaxOpenConns(cfg.PoolMaxConnections)
	sqlDb.SetMaxIdleConns(cfg.PoolMaxIdleConnections)

	return &DB{db: db}, nil
}

// dbToMysqlDSN builds a connection string suitable for the mysql driver, using
// the values of vars.
func dbToMysqlDSN(cfg *Config) string {
	mySqlConfig := mysql.NewConfig()
	mySqlConfig.Addr = cfg.Address
	mySqlConfig.Passwd = cfg.Password
	mySqlConfig.Net = cfg.Protocol
	mySqlConfig.User = cfg.User
	mySqlConfig.Timeout = time.Duration(cfg.ConnectionTimeout) * time.Second
	mySqlConfig.DBName = cfg.Name
	mySqlConfig.ParseTime = true

	if cfg.EnableSSL {
		tlsConfigName := "custom"
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(cfg.SSLRootCertPath)
		if err != nil {
			panic(err)
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			panic("Failed to append PEM.")
		}
		clientCert := make([]tls.Certificate, 0, 1)
		certs, err := tls.LoadX509KeyPair(cfg.SSLCertPath, cfg.SSLKeyPath)
		if err != nil {
			panic(err)
		}
		clientCert = append(clientCert, certs)
		err = mysql.RegisterTLSConfig(tlsConfigName, &tls.Config{
			RootCAs:      rootCertPool,
			Certificates: clientCert,
		})
		if err != nil {
			panic(err)
		}
		mySqlConfig.TLSConfig = tlsConfigName
	}
	return mySqlConfig.FormatDSN()
}
