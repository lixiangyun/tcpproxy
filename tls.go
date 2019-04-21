package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
)

func TlsClientConfig(cfg *TlsConfig) *tls.Config {

	//这里读取的是根证书
	buf, err := ioutil.ReadFile(cfg.CA)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(buf)

	//加载客户端证书
	//这里加载的是服务端签发的
	cert, err := tls.LoadX509KeyPair(cfg.Cert, cfg.Key)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            pool,
		Certificates:       []tls.Certificate{cert},
	}
}

func TlsServerConfig(cfg *TlsConfig) *tls.Config {
	//这里读取的是根证书
	buf, err := ioutil.ReadFile(cfg.CA)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(buf)

	//加载服务端证书
	cert, err := tls.LoadX509KeyPair(cfg.Cert, cfg.Key)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    pool,
	}
}
