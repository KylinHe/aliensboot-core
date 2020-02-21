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
	"github.com/KylinHe/aliensboot-core/exception"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/KylinHe/aliensboot-core/task"
	"github.com/coreos/etcd/clientv3"
	"gopkg.in/mgo.v2/bson"
	"sync"
	"time"
)

type ETCDServiceCenter struct {
	sync.RWMutex

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
		Username:config.Username,
		Password:config.Password,
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

func (this *ETCDServiceCenter) GetConfigRoot() string {
	return this.configRoot
}

func (this *ETCDServiceCenter) GetServiceRoot() string {
	return this.serviceRoot
}

func (this *ETCDServiceCenter) GetNodeID() string {
	return this.node
}

func (this *ETCDServiceCenter) IsConnect() bool {
	return this.client != nil
}

func (this *ETCDServiceCenter) Close() {
	this.Lock()
	defer this.Unlock()
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

	serviceRootPath := this.GetServiceRootPath(service.GetName())
	servicePath := serviceRootPath + NodeSplit + service.GetID()

	rsp, err := this.client.Get(newTimeoutContext(), serviceRootPath, clientv3.WithPrefix())
	if err != nil {
		log.Errorf("get service %v error: %v", serviceRootPath, err)
		return false
	}

	for _, v := range rsp.Kvs {
		path := string(v.Key)
		if config.Unique && serviceRootPath != path {
			log.Errorf("unique service %v - %v already exist.", service.GetName(), path)
			return false
		}
		if servicePath == path {
			log.Errorf("service %v - %v already exist.", service.GetName(), path)
			return false
		}
	}

	_, err = this.client.Put(newTimeoutContext(), serviceRootPath, config.Lbs)
	if err != nil {
		log.Errorf("public service config %v  err : %v", serviceRootPath, err)
		return false
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
	//if config.Local && !this.UpdateService(service, false) {
	//	return false
	//}

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
	task.SafeGo(func(){
		if err := recover(); err != nil {
			exception.PrintStackDetail(err)
		}
		for {
			select {
			case <-ticker.C:
				//this.ttlCheck.Range(this.check)
				this.RLock()
				if this.client != nil {
					resp, err := this.client.Grant(newTimeoutContext(), this.ttl)
					if err != nil {
						log.Debugf("ttl grant %v", err)
					} else {
						_, err = this.client.Put(newTimeoutContext(), path, data, clientv3.WithLease(resp.ID))
						if err != nil {
							log.Debugf("ttl update %v", err)
						}
					}
				}
				this.RUnlock()
			}
		}
	})
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

//func (this *ETCDServiceCenter) SetLbs(serviceName string, lbs string) {
//	this.Container.SetLbs(serviceName, lbs)
//}

func (this *ETCDServiceCenter) GetServiceRootPath(service string) string {
	return this.serviceRoot + NodeSplit + service
}

func (this *ETCDServiceCenter) SubscribeService(healthyOnly bool, serviceName string) {
	this.SubscribeData(this.GetServiceRootPath(serviceName), func(data []byte, init bool) {
		this.Container.SetLbs(serviceName, string(data))
	}, false)

	serviceRootPath := this.serviceRoot + NodeSplit + serviceName + NodeSplit
	err := this.AddDataPrefixListener(serviceRootPath, serviceName, this.handleService)
	if err != nil {
		log.Errorf("subscribe service %v error: %v", serviceRootPath, err)
	}
}


func (this *ETCDServiceCenter) AddDataPrefixListener(dataRootPath string, dataRootName string, handler DataPrefixListener) error {
	chidlrenData := this.GetChildrenData(dataRootPath)
	for dataName, dataValue := range chidlrenData {
		handler(PUT, dataValue, dataRootName, dataName, true)
	}
	ch := this.client.Watch(context.TODO(), dataRootPath, clientv3.WithPrefix())
	task.SafeGo(func() {
		if err := recover(); err != nil {
			exception.PrintStackDetail(err)
		}
		prefixLen := len(dataRootPath)
		for {
			//只要消息管道没有关闭，就一直等待用户请求
			event, _ := <-ch
			for _, dataEvent := range event.Events {
				dataPath := string(dataEvent.Kv.Key)
				dataName := dataPath[prefixLen:]
				if dataName != "" {
					handler(DataEventType(dataEvent.Type), dataEvent.Kv.Value, dataRootName, dataName, false)
				}
			}
		}
	})
	return nil
}

func (this *ETCDServiceCenter) handleService(eventType DataEventType, data []byte, dataRootName string, dataName string, init bool) {
	if eventType == PUT {
		centerService := &service.CenterService{}
		err1 := json.Unmarshal(data, centerService)
		if err1 != nil {
			log.Errorf("unmarshal service %v data error: %v", dataRootName, err1)
			return
		}
		service, _ := service.NewService2(centerService, dataName, dataRootName)
		this.Container.UpdateService(service, false)
	} else if eventType == DELETE {
		this.Container.RemoveService(dataRootName, dataName)
	}

}

func (this *ETCDServiceCenter) SubscribeConfig(configName string, configHandler DataListener, options ...Option) {
	configPath := this.configRoot + NodeSplit + configName
	ensure := !haveOption(OptionEmpty, options)
	this.SubscribeData(configPath, configHandler, ensure)
}

func (this *ETCDServiceCenter) GetConfigData(configName string) []byte {
	configPath := this.configRoot + NodeSplit + configName
	return this.DownloadData(configPath)
}


func haveOption(option Option, options []Option) bool {
	for _, op := range options {
		if op == option {
			return true
		}
	}
	return false
}

// 订阅前缀配置
func (this *ETCDServiceCenter) SubscribeConfigWithPrefix(configName string, listener DataPrefixListener) {
	configPath := this.configRoot + NodeSplit + configName + NodeSplit
	err := this.AddDataPrefixListener(configPath, configName, listener)
	if err != nil {
		log.Errorf("subscribe config with prefix %v error: %v", configName, err)
	}
}


func (this *ETCDServiceCenter) SubscribeData(path string, configHandler DataListener, ensure bool) {
	rsp, err := this.client.Get(newTimeoutContext(), path)
	if err != nil || rsp.Kvs == nil{
		if ensure {
			log.Fatalf("subscribe config %v error: %v", path, err)
		}
	} else {
		for _, v := range rsp.Kvs {
			configHandler(v.Value, true)
		}
	}
	ch := this.client.Watch(context.TODO(), path)
	task.SafeGo(func() {
		if err := recover(); err != nil {
			exception.PrintStackDetail(err)
		}
		for {
			//只要消息管道没有关闭，就一直等待用户请求
			event, _ := <-ch
			for _, serviceEvent := range event.Events {
				if serviceEvent.Type == clientv3.EventTypePut {
					if serviceEvent.Kv.Value == nil || len(serviceEvent.Kv.Value) == 0 {
						log.Errorf("invalid config %v", path)
					} else {
						configHandler(serviceEvent.Kv.Value, false)
					}
				}
			}
		}
	})
}

func (this *ETCDServiceCenter) PublicConfigData(configName string, data interface{}) bool {
	content, _ := json.Marshal(data)
	return this.PublicConfig(configName, content)
}

func (this *ETCDServiceCenter) DeleteConfig(configName string) bool {
	configPath := this.configRoot + NodeSplit + configName
	return this.DeleteData(configPath)
}

func (this *ETCDServiceCenter) PublicConfig(configName string, configContent []byte) bool {
	if configName == "" {
		log.Info("config type con not be empty")
		return false
	}
	configPath := this.configRoot + NodeSplit + configName
	return this.UploadData(configPath, configContent)
}

func (this *ETCDServiceCenter) GetChildrenData(path string) map[string][]byte {
	result := make(map[string][]byte)
	prefixLen := len(path)
	rsp, err := this.client.Get(newTimeoutContext(), path, clientv3.WithPrefix())
	if err != nil {
		log.Errorf("download data %v error: %v", path, err)
		return result
	}
	for _, v := range rsp.Kvs {
		dataPath := string(v.Key)
		dataName := dataPath[prefixLen:]
		if dataName != "" {
			result[dataName] = v.Value
		}
	}
	return result
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

func (this *ETCDServiceCenter) DeleteData(path string) bool {
	_, err := this.client.Delete(newTimeoutContext(), path)
	if err != nil {
		log.Errorf("delete data %v failed err : %v", path, err)
		return false
	}
	log.Infof("delete data %v success", path)
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
