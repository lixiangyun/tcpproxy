package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
)

func TlsClientConfig(cfg *TlsConfig, addr string) *tls.Config {
	var pool *x509.CertPool

	if cfg.CA != "" {
		buf, err := ioutil.ReadFile(cfg.CA)
		if err != nil {
			log.Fatal(err.Error())
			return nil
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(buf)
	}

	cert, err := tls.LoadX509KeyPair(cfg.Cert, cfg.Key)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	var bSkipVerify bool

	// 如果没有配置服务端根证书，则忽略校验服务端证书有效性。
	if pool == nil {
		bSkipVerify = true
	}

	return &tls.Config{
		ServerName:         addr,
		InsecureSkipVerify: bSkipVerify,
		RootCAs:            pool,
		Certificates:       []tls.Certificate{cert},
	}
}

func TlsServerConfig(cfg *TlsConfig) *tls.Config {
	var pool *x509.CertPool

	if cfg.CA != "" {
		//这里读取的是根证书
		buf, err := ioutil.ReadFile(cfg.CA)
		if err != nil {
			log.Fatal(err.Error())
			return nil
		}
		pool = x509.NewCertPool()
		pool.AppendCertsFromPEM(buf)
	}

	//加载服务端证书
	crt, err := tls.LoadX509KeyPair(cfg.Cert, cfg.Key)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	var authtype tls.ClientAuthType

	if pool != nil {
		authtype = tls.RequireAndVerifyClientCert
	} else {
		authtype = tls.RequireAnyClientCert
	}

	return &tls.Config{
		Certificates: []tls.Certificate{crt},
		ClientAuth:   authtype,
		ClientCAs:    pool,
	}
}
