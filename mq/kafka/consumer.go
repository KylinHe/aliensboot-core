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
	"github.com/bsm/sarama-cluster" //support automatic consumer-group rebalancing and offset tracking
	"time"
)

type Consumer struct {
	proxy *cluster.Consumer
}

//address = strings.Split("localhost:9092", ",")
func (this *Consumer) Init(address []string, service string, node string, handle func(data []byte) error) error {
	config := cluster.NewConfig()
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetNewest //初始从最新的offset开始
	c, err := cluster.NewConsumer(address, service, []string{service, service + node}, config)
	if err != nil {
		return err
	}

	this.proxy = c
	go func(c *cluster.Consumer) {
		errors := c.Errors()
		noti := c.Notifications()
		for {
			select {
			case err := <-errors:
				if err != nil {
					log.Error(err)
				}
			case notify := <-noti:
				if notify != nil {
					log.Debug(notify)
				}
			}
		}
	}(this.proxy)

	go func() {
		defer c.Close()
		for msg := range c.Messages() {
			handle(msg.Value)
			//fmt.Fprintf(os.Stdout, "%s %s/%d/%d\t%s\n",service, msg.Topic, msg.Partition, msg.Offset, msg.Value)
			c.MarkOffset(msg, "") //MarkOffset 并不是实时写入kafka，有可能在程序crash时丢掉未提交的offset
		}
		//log.Error("consumer close!")
	}()

	return nil
}

func (this *Consumer) Close() error {
	if this.proxy != nil {
		return this.proxy.Close()
	}
	return nil
}
