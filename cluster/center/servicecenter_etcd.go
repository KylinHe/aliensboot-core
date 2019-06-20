/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 *
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package center

//服务中心，处理服务的调度和查询
import (
	"context"
	"encoding/json"
	"github.com/KylinHe/aliensboot-core/cluster/center/service"
	"github.com/KylinHe/aliensboot-core/config"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"gopkg.in/mgo.v2/bson"
	"sync"
	"time"
)

type ETCDServiceCenter struct {
	//sync.RWMutex

	*service.Container //服务容器 key 服务名
	client             *clientv3.Client

	node string //当前节点
	configRoot  string //配置根节点
	serviceRoot string //服务根节点

	ttl         int64
	ttlCheck    *sync.Map //map[string]string
	//ticker      *time.Ticker
}

func (this *ETCDServiceCenter) ConnectCluster(config config.ClusterConfig) {
	if config.ID == "" || config.ID == "-" {
		config.ID = bson.NewObjectId().Hex()
	}
	if config.Timeout == 0 {
		config.Timeout = 5
	}
	if config.TTL == 0 {
		config.TTL = 30
	}

	etcdConfig := clientv3.Config{
		Endpoints:   config.Servers,
		DialTimeout: time.Second * time.Duration(config.Timeout),
	}
	client, err := clientv3.New(etcdConfig)
	if err != nil {
		log.Fatal(err)
	}
	this.client = client
	this.ttlCheck = &sync.Map{} //make(map[string]string)
	this.serviceRoot = NodeSplit + "root" + NodeSplit + config.Name + NodeSplit + ServiceNodeName
	this.configRoot = NodeSplit + "root" + NodeSplit + config.Name + NodeSplit + ConfigNodeName

	this.node = config.ID
	this.ttl = config.TTL
	//this.listeners = make(map[string]struct{})

	this.Container = service.NewContainer()

	//开启ttl
	//go this.openTTLCheck()
}

func (this *ETCDServiceCenter) GetNodeID() string {
	return this.node
}

func (this *ETCDServiceCenter) IsConnect() bool {
	return this.client != nil
}

func (this *ETCDServiceCenter) Close() {
	if this.client != nil {
		_ = this.client.Close()
		this.client = nil
	}
}

//释放服务
func (this *ETCDServiceCenter) ReleaseService(service service.IService) {
	servicePath := this.serviceRoot + NodeSplit + service.GetName() + NodeSplit + service.GetID()
	ticker, ok := this.ttlCheck.Load(servicePath)
	if ok {
		ticker.(*time.Ticker).Stop()
		this.ttlCheck.Delete(servicePath)
	}
	this.Container.RemoveService(service.GetName(), service.GetID())

	_, err := this.client.Delete(newTimeoutContext(), servicePath)
	if err != nil {
		log.Errorf("release service %v err : %v", servicePath, err)
	} else {
		log.Infof("release service %v success", servicePath)
	}
}

//func (this *ETCDServiceCenter) AddServiceListener(listener service.Listener) {
//	this.Container.AddServiceListener(listener)
//}

func (this *ETCDServiceCenter) PublicService(service service.IService, config config.ServiceConfig) bool {
	if !service.IsLocal() {
		log.Error("service info is invalid")
		return false
	}

	serviceRootPath := this.serviceRoot + NodeSplit + service.GetName()
	servicePath := serviceRootPath + NodeSplit + service.GetID()

	rsp, err := this.client.Get(newTimeoutContext(), serviceRootPath, clientv3.WithPrefix())
	if err != nil {
		log.Errorf("get service %v error: %v", serviceRootPath, err)
		return false
	}

	if config.Unique && len(rsp.Kvs) > 0 {
		log.Errorf("unique service %v already exist.", service.GetName())
		return false
	}

	for _, v := range rsp.Kvs {
		path := string(v.Key)
		if servicePath == path {
			log.Errorf("service %v - %v already exist.", service.GetName(), servicePath)
			return false
		}
	}

	_, ok := this.ttlCheck.Load(servicePath)
	//this.RLock()
	//ttlData := this.ttlCheck[servicePath]
	//this.RUnlock()
	if ok {
		log.Errorf("service %v already public : %v", servicePath)
		return false
	}

	//ttlCheck : 10s
	data, err := json.Marshal(service)

	//this.PublicConfig("testconfig", data)
	if err != nil {
		log.Errorf("marshal json service data error : %v", err)
		return false
	}

	//允许本地调用 需要注册在服务容器中
	if config.Local && !this.UpdateService(service, false) {
		return false
	}

	resp, _ := this.client.Grant(context.TODO(), this.ttl)
	serviceData := string(data)
	_, err1 := this.client.Put(newTimeoutContext(), servicePath, serviceData, clientv3.WithLease(resp.ID))
	if err1 != nil {
		log.Errorf("create service error : %v", err1)
		return false
	}
	this.openTTLCheck(servicePath, serviceData)
	log.Infof("public %v success", servicePath)
	return true
}

func newTimeoutContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx
}

func (this *ETCDServiceCenter) openTTLCheck(path string, data string) {
	ticker := time.NewTicker(time.Second * time.Duration(this.ttl/2))
	this.ttlCheck.Store(path, ticker)
	go func(){
		for {
			select {
			case <-ticker.C:
				//this.ttlCheck.Range(this.check)
				resp, _ := this.client.Grant(context.TODO(), this.ttl)
				//log.Debugf("ttl updata %v - %v", path, data)
				_, err := this.client.Put(newTimeoutContext(), path, data, clientv3.WithLease(resp.ID))
				if err != nil {
					log.Debugf("ttl update %v", err)
				}

			}
		}
	}()
}

//func (this *ETCDServiceCenter) check(path, data interface{}) bool {
//	resp, _ := this.client.Grant(context.TODO(), this.ttl)
//	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
//	//log.Debugf("ttl updata %v - %v", path, data)
//	_, err := this.client.Put(ctx, path.(string), data.(string), clientv3.WithLease(resp.ID))
//	if err != nil {
//		log.Debugf("ttl update %v", err)
//	}
//	return true
//}

func (this *ETCDServiceCenter) SubscribeServices(serviceNames ...string) {
	for _, serviceName := range serviceNames {
		this.SubscribeService(true, serviceName)
	}
}

func (this *ETCDServiceCenter) SetLbs(serviceName string, lbs string) {
	this.Container.SetLbs(serviceName, lbs)
}

func (this *ETCDServiceCenter) SubscribeService(healthyOnly bool, serviceName string) {
	//this.SubscribeConfig("lbs"+NODE_SPLIT+serviceName, func(data []byte) {
	//	this.Container.SetLbs(serviceName, string(data))
	//})

	serviceRootPath := this.serviceRoot + NodeSplit + serviceName + NodeSplit
	prefixLen := len(serviceRootPath)

	rsp, err := this.client.Get(newTimeoutContext(), serviceRootPath, clientv3.WithPrefix())
	//serviceIDs, _, ch, err := this.zkCon.ChildrenW(serviceRootPath)
	if err != nil {
		log.Errorf("subscribe service %v error: %v", serviceRootPath, err)
		return
	}
	for _, v := range rsp.Kvs {
		this.handleService(mvccpb.PUT, v, serviceName, prefixLen)
	}
	go this.openListener(serviceName, serviceRootPath)
}

func (this *ETCDServiceCenter) openListener(serviceName string, serviceRootPath string) {
	ch := this.client.Watch(context.TODO(), serviceRootPath, clientv3.WithPrefix())
	prefixLen := len(serviceRootPath)
	for {
		//只要消息管道没有关闭，就一直等待用户请求
		event, _ := <-ch
		for _, serviceEvent := range event.Events {
			this.handleService(serviceEvent.Type, serviceEvent.Kv, serviceName, prefixLen)
		}
	}
}

func (this *ETCDServiceCenter) handleService(eventType mvccpb.Event_EventType, v *mvccpb.KeyValue, serviceName string, prefixLen int) {
	servicePath := string(v.Key)
	data := v.Value
	serviceID := servicePath[prefixLen:]

	if eventType == clientv3.EventTypePut {
		centerService := &service.CenterService{}
		err1 := json.Unmarshal(data, centerService)
		if err1 != nil {
			log.Errorf("unmarshal service %v data error: %v", servicePath, err1)
			return
		}
		service := service.NewService2(centerService, serviceID, serviceName)
		this.Container.UpdateService(service, false)
	} else if eventType == clientv3.EventTypeDelete {
		this.Container.RemoveService(serviceName, serviceID)
	}

}

func (this *ETCDServiceCenter) SubscribeConfig(configName string, configHandler ConfigListener) {
	configPath := this.configRoot + NodeSplit + configName
	rsp, err := this.client.Get(newTimeoutContext(), configPath)
	if err != nil {
		log.Fatalf("subscribe config %v error: %v", configPath, err)
		return
	}
	for _, v := range rsp.Kvs {
		configHandler(v.Value)
	}

	go func() {
		ch := this.client.Watch(context.TODO(), configPath)
		for {
			//只要消息管道没有关闭，就一直等待用户请求
			event, _ := <-ch
			for _, serviceEvent := range event.Events {
				if serviceEvent.Type == clientv3.EventTypePut {
					configHandler(serviceEvent.Kv.Value)
				}
			}
		}
	}()

}

func (this *ETCDServiceCenter) PublicConfig(configType string, configContent []byte) bool {
	if configType == "" {
		log.Info("config type con not be empty")
		return false
	}
	configPath := this.configRoot + NodeSplit + configType

	_, err := this.client.Put(newTimeoutContext(), configPath, string(configContent))
	if err != nil {
		log.Errorf("public config %v  err : %v", configType, err)
		return false
	}
	log.Infof("public config %v success", configType)
	return true
}

//上传数据
func (this *ETCDServiceCenter) UploadData(path string, configContent []byte) bool {
	_, err := this.client.Put(newTimeoutContext(), path, string(configContent))
	if err != nil {
		log.Errorf("upload data %v failed err : %v", path, err)
		return false
	}
	log.Infof("upload data %v success", path)
	return true
}

func (this *ETCDServiceCenter) DownloadData(path string) []byte {
	rsp, err := this.client.Get(newTimeoutContext(), path)
	if err != nil {
		log.Errorf("download data %v error: %v", path, err)
		return nil
	}
	for _, v := range rsp.Kvs {
		return v.Value
	}
	return nil
}
