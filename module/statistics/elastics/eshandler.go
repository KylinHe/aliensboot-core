/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/8/3
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package elastics

import (
	"fmt"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/KylinHe/aliensboot-core/module/statistics/conf"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v5"
	"gopkg.in/sohlich/elogrus.v2"
	"time"
)

var tag = ""

var format = &logrus.JSONFormatter{DisableTimestamp: true}

func init() {
	//flag.StringVar(&tag, "etag", "", "testcase tag")  //测试用例标识
	//flag.Parse()
	//log.Debug("statistics tag is %v", tag)
}

func NewESHandler(prefix string) *esHandler {
	return &esHandler{
		prefix:    prefix,
		esLogs:    make(map[string]*logrus.Logger),
		dayPrefix: time.Now().Format("2006-01-02"),
	}
}

type esHandler struct {
	esClient  *elastic.Client
	prefix    string
	esLogs    map[string]*logrus.Logger
	dayPrefix string
}

func (this *esHandler) UpdateDayPrefix() {
	this.dayPrefix = time.Now().Add(time.Hour).Format("2006-01-02")
}

//func handleESLog(args []interface{}, ) {
//	index := args[0].(constant.EsLogIndex)
//	logger := getESLogger(index)
//	msg := args[1]
//	fields := args[2].(logrus.Fields)
//	logger.WithFields(fields).Debug(msg)
//}

func (this *esHandler) HandleDayESLog(index string, msg string, fields logrus.Fields) {
	index = index + "_" + this.dayPrefix
	logger := this.getESLogger(index)
	if tag != "" {
		fields["tag"] = tag
	}
	logger.WithFields(fields).Debug(msg)
}

func (this *esHandler) getESLogger(index string) *logrus.Logger {
	logger := this.esLogs[index]
	if logger == nil {
		logger = log.NewLogger(conf.LocalPrefix, format, conf.Local)
		err := this.attachES(logger, index)
		if err != nil {

		}
		this.esLogs[index] = logger
	}
	return logger
}

//关联elasticsearch
func (this *esHandler) attachES(logger *logrus.Logger, index string) error {
	if conf.Config.ES.Url == "" {
		return errors.New("invalid es config")
	}
	if this.esClient == nil {
		client, err := elastic.NewClient(elastic.SetURL(conf.Config.ES.Url))
		if client == nil || err != nil {
			return errors.New(fmt.Sprintf("%v config es logger error. %v", client, errors.WithStack(err)))
		}
		this.esClient = client
	}

	index = this.prefix + "_" + index
	esHook, err := elogrus.NewElasticHook(this.esClient, conf.Config.ES.Host, logrus.DebugLevel, index)
	//esHook, err := elogrus.NewElasticHook(client, esHOst, log.DebugLevel, index)
	if err != nil {
		return errors.New(fmt.Sprintf("config es logger error. %+v", errors.WithStack(err)))
	}
	logger.AddHook(esHook)
	return nil
}
