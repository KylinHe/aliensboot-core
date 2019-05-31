/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/10/31
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 * Desc: compatible log framework
 *******************************************************************************/
package log

import (
	"fmt"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"time"
)

var logger = NewLogger("",  &logrus.TextFormatter{}, false)

var logRoot string

var DEBUG = false

//调试版本日志带颜色
func Init(debug bool, tag string, logDir string) {
	DEBUG = debug
	format := &logrus.TextFormatter{}
	format.ForceColors = DEBUG
	format.DisableColors = !DEBUG
	format.DisableTimestamp = DEBUG
	logger.Formatter = format
	logRoot = logDir
	configLocalFilesystemLogger(tag, logger)
	//被lfshook修改需要改回来
	format.DisableColors = !DEBUG
}

func NewLogger(name string, formatter logrus.Formatter, local bool) *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = formatter
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.Out = os.Stdout
	// Only log the warning severity or above.
	logger.Level = logrus.DebugLevel
	if local {
		configLocalFilesystemLogger(name, logger)
	}
	return logger
}

// config logrus log to amqp  rabbitMQ
//func ConfigAmqpLogger(server, username, password, exchange, exchangeType, virtualHost, routingKey string) {
//	hook := logrus_amqp.NewAMQPHookWithType(server, username, password, exchange, exchangeType, virtualHost, routingKey)
//	log.AddHook(hook)
//}

// config logrus log to elasticsearch
//func ConfigESLogger(esUrl string, esHOst string, index string) {
//	client, err := elastic.NewClient(elastic.SetURL(esUrl))
//	if err != nil {
//		log.Errorf("config es logger error. %+v", errors.WithStack(err))
//	}
//	esHook, err := elogrus.NewElasticHook(client, esHOst, log.DebugLevel, index)
//	if err != nil {
//		log.Errorf("config es logger error. %+v", errors.WithStack(err))
//	}
//	log.AddHook(esHook)
//}

//config logrus log to local file
func configLocalFilesystemLogger(name string, logger *logrus.Logger) {
	maxAge := 30 * 24 * time.Hour
	rotationTime := 24 * time.Hour
	logFileName := name + ".log"

	os.Mkdir(logRoot, os.ModePerm)
	baseLogPath := path.Join(logRoot, logFileName)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		fmt.Println(fmt.Errorf("config local file system logger error. %+v", errors.WithStack(err)))
	}

	errWriter, err1 := rotatelogs.New(
		baseLogPath+".err.%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err1 != nil {
		fmt.Println(fmt.Errorf("config local file system err logger error. %+v", errors.WithStack(err1)))
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: errWriter,
		logrus.FatalLevel: errWriter,
		logrus.PanicLevel: errWriter,
	}, logger.Formatter)
	logger.AddHook(lfHook)
}

//Debugf Printf Infof Warnf Warningf Errorf Panicf Fatalf

//做一层适配，方便后续切换到其他日志框架或者自己写

//-----------format
func WithField(key string, value interface{}) *logrus.Entry {
	return logger.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}

func getLocation() string {
	pc, _, lineno, ok := runtime.Caller(2)
	src := ""
	if ok {
		src = fmt.Sprintf("[%s:%d] ", runtime.FuncForPC(pc).Name(), lineno)
	}
	return src
}

//Debugf Printf Infof Warnf Warningf Errorf Panicf Fatalf

//做一层适配，方便后续切换到其他日志框架或者自己写
func Debug(arg ...interface{}) {
	logger.Debug(arg...)
}

func Print(arg ...interface{}) {
	logger.Print(arg...)
}

func Info(arg ...interface{}) {
	logger.Info(arg...)
}

func Warn(arg ...interface{}) {
	logger.Warn(arg...)
}

func Error(arg ...interface{}) {
	logger.Error(arg...)
}

func Panic(arg ...interface{}) {
	logger.Panic(arg...)
}

func Fatal(arg ...interface{}) {
	logger.Fatal(arg...)
}

//-----------format

func Debugf(format string, arg ...interface{}) {
	if DEBUG {
		format = getLocation() + format
	}
	logger.Debugf(format, arg...)
}

func Printf(format string, arg ...interface{}) {
	if DEBUG {
		format = getLocation() + format
	}
	logger.Printf(format, arg...)
}

func Infof(format string, arg ...interface{}) {
	if DEBUG {
		format = getLocation() + format
	}
	logger.Infof(format, arg...)
}

func Warnf(format string, arg ...interface{}) {
	if DEBUG {
		format = getLocation() + format
	}
	logger.Warnf(format, arg...)
}

func Errorf(format string, arg ...interface{}) {
	if DEBUG {
		format = getLocation() + format
	}
	logger.Errorf(format, arg...)
}

func Panicf(format string, arg ...interface{}) {
	if DEBUG {
		format = getLocation() + format
	}
	logger.Panicf(format, arg...)
}

func Fatalf(format string, arg ...interface{}) {
	if DEBUG {
		format = getLocation() + format
	}
	logger.Fatalf(format, arg...)
}
