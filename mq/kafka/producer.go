/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/4/21
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package kafka

import (
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/Shopify/sarama"
	"time"
)

type Producer struct {
	proxy sarama.AsyncProducer
}

func (this *Producer) Init(address []string, timeout int) error {
	if timeout == 0 {
		timeout = 5
	}
	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//随机的分区类型
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	//config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Timeout = time.Duration(timeout) * time.Second
	p, err := sarama.NewAsyncProducer(address, config)
	//defer p.Close()
	if err != nil {
		return err
	}

	this.proxy = p
	//必须有这个匿名函数内容
	go func(p sarama.AsyncProducer) {
		defer p.AsyncClose()
		errors := p.Errors()
		//success := p.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					log.Error(err)
				}
				//case succ := <-success:
				//	if succ != nil {
				//		log.Debug(succ)
				//	}
			}
		}
	}(this.proxy)
	return nil
}

func (this *Producer) Close() error {
	if this.proxy != nil {
		return this.proxy.Close()
	}
	return nil
}

func (this *Producer) Broadcast(service string, data []byte) {
	msg := &sarama.ProducerMessage{
		Topic: service,
		Value: sarama.ByteEncoder(data),
	}
	this.proxy.Input() <- msg
}

func (this *Producer) SendMessage(service string, node string, data []byte) {
	msg := &sarama.ProducerMessage{
		Topic: service + node,
		Value: sarama.ByteEncoder(data),
	}
	this.proxy.Input() <- msg
}
