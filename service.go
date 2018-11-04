package main 

import (
	"os"
	"log"
	"net"
	"time"
	"flag"
	"path/filepath"
	//"fmt"
	//"errors"
	//"io/ioutil"
	"github.com/kardianos/service"
	"gopkg.in/natefinch/lumberjack.v2" // logging
	"github.com/jackc/pgx"
	//"github.com/fasthttp-contrib/websocket"
)


type appServer struct {
	http_listener net.Listener
	appPath string

	//ws_listener net.Listener
	db_connpool *pgx.ConnPool
	//ws_upgrader websocket.Upgrader
}


func (s *appServer) Start(svc service.Service) error {
	// Start should not block. Do the actual work async.
	LogInfo("start servico")
	go s.run()
	return nil
}

func (s *appServer) run() error {
	
	var err error
	var outerr error
	
	outerr = nil

	s.db_connpool, err = DoPoolConnect(s.appPath + DBCONNCFG, s.appPath + DBSQLSTATMTS, DBINITSMTGRP)
	if err != nil {
		LogCriticalf("database: FATAL error, no connection: %s", err)
	}
		
	err, s.http_listener = AsyncListenAndServe(ADDR_HS, s.hsmux, time.Duration(SHUTDDELAY_SECS) * time.Second, "HTTP")
	if err != nil {
		LogCriticalf("httpserver: FATAL error in AsyncListenAndServe: %s", err)
		outerr = err
	}
	
	/* REMOVER para ter WebSockets
	s.prepareWebsockets()

	err, s.ws_listener = AsyncListenAndServe(ADDR_WS, s.wsmux, time.Duration(SHUTDDELAY_SECS) * time.Second, "websocket")
	if err != nil {
		LogCriticalf("wsockserver: FATAL error in AsyncListenAndServe: %s", err)
	}
	* */
	
	return outerr
}

func (s *appServer) stop() {
	s.db_connpool.Close()
	//s.ws_listener.Close()
	s.http_listener.Close()
}

func (s *appServer) Stop(svc service.Service) error {
	LogInfo("stop servico")
	s.stop()
	return nil
}

func timedserve(appPath string) {
	as := &appServer{appPath: appPath}
	err := as.run()
	if err != nil {
		LogInfo("Timed server closing in error")
		as.stop()
	} else {
		time.Sleep(time.Duration(TIMEDSERVER_MINS) * time.Minute)
		LogInfo("a fechar timed server")
		as.stop()
	}
}

func serviceMain(appPath string) {
	
	svcConfig := &service.Config{
		Name:   SVC_NAME,
		/*DisplayName: "XXXXX",
		Description: "yYYYYYYYY",*/
	}

	prg := &appServer{appPath: appPath}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		LogError(err.Error())
	}

	err = s.Run()
	if err != nil {
		LogCritical(err.Error())
	}
}


func main() {
	
	var noservice bool

	
	/* ATENCAO -- TEMPORARIO !!! */
	flag.BoolVar(&noservice, "noservice", false, "como processo autonomo")	
	flag.Parse()
	
	//noservice = true
    ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)
        
	log.SetOutput(&lumberjack.Logger{
		Filename:   exPath + LOGPATH,
		MaxSize:    1, // megabytes
		MaxBackups: 5,
		MaxAge:     28, //days
	})
	
	if noservice { 
		timedserve(exPath)
	} else {
		serviceMain(exPath)
	}
}
