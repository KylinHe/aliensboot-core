/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/4/8
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package mq

import (
	"fmt"
	"github.com/KylinHe/aliensboot-core/mq/kafka"
	"github.com/pkg/errors"
)

type Type string

const (
	KAFKA Type = "kafka"
)

//消息生产者
type IProducer interface {
	Init(address []string, timeout int) error
	SendMessage(service string, node string, data []byte) //异步发送数据
	Broadcast(service string, data []byte)                //广播数据
	Close() error
}

//消息消费者
type IConsumer interface {
	Init(address []string, service string, node string, handle func(data []byte) error) error
	Close() error
}

func NewProducer(config Config) (producer IProducer, err error) {
	if config.Type == KAFKA {
		producer = &kafka.Producer{}
	}
	if producer != nil {
		err = producer.Init(config.Address, config.Timeout)
	} else {
		err = errors.New(fmt.Sprintf("un expect mq producer type %v", config.Type))
	}
	return
}

func NewConsumer(config Config, service string, node string, handle func(data []byte) error) (consumer IConsumer, err error) {
	if config.Type == KAFKA {
		consumer = &kafka.Consumer{}
	}
	if consumer != nil {
		err = consumer.Init(config.Address, service, node, handle)
	} else {
		err = errors.New(fmt.Sprintf("un expect mq producer type %v", config.Type))
	}
	return
}
