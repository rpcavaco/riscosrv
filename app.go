package main 

import (
	//"os"
	"fmt"
	"strconv"
	//"strings"
	//"regexp"
	//"errors"
	//"io/ioutil"
	//"encoding/json"
	//"github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	//"github.com/smallfish/simpleyaml"
	//"github.com/fasthttp-contrib/websocket"
	//"github.com/jackc/pgx"
)

var ADDR_HS string = ":8010"
//var ADDR_WS string = ":8094"
var SHUTDDELAY_SECS  int = 5
var TIMEDSERVER_MINS int = 20
var SVC_NAME string = "LocService"
//var APPPATH string = "C:\\GOWorkspace\\src\\github.com\\rpcavaco\\riscosrv"
var LOGPATH string = "\\log.txt"
var DBCONNCFG = "\\dbconn_config.json"
var DBSQLSTATMTS = "\\sqlstatements.yaml"
var DBINITSMTGRP = "initprepared"
//var SRVROOT = "C:\\www\\docsplat"


/*
func validFileServerExtension(path string) bool {
	
	var ret bool = false
	
	if path == "/" {
		return true
	}
	
	lowerpath := strings.ToLower(path)

	htmpatt := regexp.MustCompile("\\.(htm[l]?|json)$")
	imgpatt := regexp.MustCompile("\\.(jp[e]?g|png|gif|tif[f]?)$")
	webpatt := regexp.MustCompile("\\.(svg|js|css|ttf)$")
	pltxtpatt := regexp.MustCompile("\\.(txt|md|mkd|csv)$")

	switch {
		case htmpatt.MatchString(lowerpath):
			ret = true
		case imgpatt.MatchString(lowerpath):
			ret = true
		case webpatt.MatchString(lowerpath):
			ret = true
		case pltxtpatt.MatchString(lowerpath):
			ret = true
	} 
	
	return ret
	
}


// HTTP Server

var fs *fasthttp.FS = &fasthttp.FS{
		Root:               SRVROOT,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: false,
		Compress:           false,
		AcceptByteRange:    false,
	}



var fsHandler func(hsctx *fasthttp.RequestCtx) = fs.NewRequestHandler()
*/
	
func (s* appServer) featsHandler(hsctx *fasthttp.RequestCtx) {
	
	var vcnt, chunks, chunk int64
	var ferr error
	var outj string	
	
	sreqid := hsctx.QueryArgs().Peek("reqid")
	if len(sreqid) < 1 {
		LogErrorf("featsHandler parse params error, no reqid")
		hsctx.Error("featsHandler parse params error, no reqid", fasthttp.StatusInternalServerError)
		return;
	}	

	svertxcnt := hsctx.QueryArgs().Peek("vertxcnt")
	schunks := hsctx.QueryArgs().Peek("chunks")
	lname := hsctx.QueryArgs().Peek("lname")
	
	vcnt, ferr = strconv.ParseInt(string(svertxcnt), 10, 64)
	if ferr == nil {
		chunks, ferr = strconv.ParseInt(string(schunks), 10, 64)
	}	
	if ferr == nil {
		schunk := hsctx.QueryArgs().Peek("chunk")
		if len(schunk) < 1 {
			chunk = 1
		} else {
			chunk, ferr = strconv.ParseInt(string(schunk), 10, 64)
		}
	}	
	if ferr != nil {
		LogErrorf("featsHandler parse params error %s", ferr.Error())
		hsctx.Error("featsHandler parse params error", fasthttp.StatusInternalServerError)
	} else {
		lname = hsctx.QueryArgs().Peek("lname")
		if len(lname) < 1 {
			LogErrorf("featsHandler parse params error, no layer name")
			hsctx.Error("featsHandler parse params error, no layer name", fasthttp.StatusInternalServerError)
		} else {
			qryname := "initprepared.getfeats"
			LogTwitf("feats: %s %s %d %d %d", string(sreqid), lname, chunks, vcnt, chunk)
			row := s.db_connpool.QueryRow(qryname, sreqid, lname, chunks, vcnt, chunk)
			err := row.Scan(&outj)
			if err != nil {
				LogErrorf("featsHandler dbquery return read error %s, stmt name: '%s'", err.Error(), qryname)
				hsctx.Error("dbquery return read error", fasthttp.StatusInternalServerError)
			} else {			
				fmt.Fprintf(hsctx, outj)				
				hsctx.SetContentType("application/json; charset=utf8")
			}
		}		
	}	
}

func (s* appServer) statsHandler(hsctx *fasthttp.RequestCtx) {

	var cenx, ceny, wid, hei, pixsz float64
	var ferr error
	var outj string	
	
	scenx := hsctx.QueryArgs().Peek("cenx")
	sceny := hsctx.QueryArgs().Peek("ceny")
	swid := hsctx.QueryArgs().Peek("wid")
	shei := hsctx.QueryArgs().Peek("hei")
	spixsz := hsctx.QueryArgs().Peek("pixsz")
	
	cenx, ferr = strconv.ParseFloat(string(scenx), 64)
	if ferr == nil {
		ceny, ferr = strconv.ParseFloat(string(sceny), 64)
	}
	if ferr == nil {
		wid, ferr = strconv.ParseFloat(string(swid), 64)
	}
	if ferr == nil {
		hei, ferr = strconv.ParseFloat(string(shei), 64)
	}
	if ferr == nil {
		pixsz, ferr = strconv.ParseFloat(string(spixsz), 64)
	}
	if ferr != nil {
		LogErrorf("statsHandler parse params error %s", ferr.Error())
		hsctx.Error("statsHandler parse params error", fasthttp.StatusInternalServerError)
	} else {

		mapname := hsctx.QueryArgs().Peek("map")
		if len(mapname) < 1 {
			LogErrorf("statsHandler parse params error, no map name")
			hsctx.Error("statsHandler parse params error, no map name", fasthttp.StatusInternalServerError)
		} else {
			vizlayers := hsctx.QueryArgs().Peek("vizlrs")
			filter_lname := hsctx.QueryArgs().Peek("flname")
			filter_fname := hsctx.QueryArgs().Peek("ffname")
			filter_value := hsctx.QueryArgs().Peek("fval")

			qryname := "initprepared.fullchunkcalc"
			LogTwitf("stats: %f %f %f %f %f %s %s %s %s %s", cenx, ceny, pixsz, wid, hei, mapname, vizlayers, filter_lname, filter_fname, filter_value)
			row := s.db_connpool.QueryRow(qryname, cenx, ceny, pixsz, wid, hei, mapname, vizlayers, filter_lname, filter_fname, filter_value)
			err := row.Scan(&outj)
			if err != nil {
				LogErrorf("statsHandler dbquery return read error %s, stmt name: '%s'", err.Error(), qryname)
				hsctx.Error("dbquery return read error", fasthttp.StatusInternalServerError)
			} else {			
				fmt.Fprintf(hsctx, outj)				
				hsctx.SetContentType("application/json; charset=utf8")
			}
		}		
	}
}

func (s* appServer) testRequestHandler(hsctx *fasthttp.RequestCtx) {
	
	fmt.Fprintf(hsctx, "Hello, world!\n\n")	
	hsctx.SetContentType("text/plain; charset=utf8")

}

func (s* appServer) hsmux(hsctx *fasthttp.RequestCtx) {
		LogTwitf("acesso HTTP: %s", hsctx.Path())
		switch string(hsctx.Path()) {
			case "/x":
				s.testRequestHandler(hsctx)
			case "/stats":
				s.statsHandler(hsctx)
			case "/feats":
				s.featsHandler(hsctx)
			default:
				/*if validFileServerExtension(string(hsctx.Path())) { 
					fsHandler(hsctx)
				} else { */
				hsctx.Error("not found", fasthttp.StatusNotFound)
				LogWarningf("HTTP not found: %s", string(hsctx.Path()))
				//}
		}
}

