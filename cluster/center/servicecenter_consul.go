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
	"fmt"
	"github.com/KylinHe/aliensboot-core/cluster/center/service"
	"github.com/KylinHe/aliensboot-core/common/util"
	"github.com/KylinHe/aliensboot-core/config"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/hashicorp/consul/api"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ConsulServiceCenter struct {
	sync.RWMutex

	*service.Container //服务容器 key 服务名
	client             *api.Client

	node      string
	listeners map[string]struct{}
}

func (this *ConsulServiceCenter) ConnectCluster(config config.ClusterConfig) {
	//if config.ID == "" {
	config.ID = bson.NewObjectId().Hex()
	//}
	this.node = config.ID

	this.listeners = make(map[string]struct{})

	consulConfig := api.DefaultConfig()
	consulConfig.Address = config.Servers[0]

	client, err := api.NewClient(consulConfig)

	if err != nil {
		log.Fatal(err)
	}
	this.client = client
	this.Container = service.NewContainer()

	go this.openListener()
}

func (this *ConsulServiceCenter) GetNodeID() string {
	return this.node
}

func (this *ConsulServiceCenter) IsConnect() bool {
	return this.client != nil
}

func (this *ConsulServiceCenter) Close() {
	if this.client != nil {
		this.client = nil
	}
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("check status.")
	fmt.Fprint(w, "status ok!")
}

func (this *ConsulServiceCenter) ReleaseService(service service.IService) {

}

func (this *ConsulServiceCenter) PublicService(service service.IService, unique bool) bool {
	var tags []string
	serviceID := service.GetName() + "-" + service.GetAddress() + "-" + util.IntToString(service.GetPort())
	checkAddress := ":" + util.IntToString(service.GetPort()+10)

	consulService := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    service.GetName(),
		Address: service.GetAddress(),
		Port:    service.GetPort(),
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			HTTP:     "http://" + checkAddress + "/status",
			Interval: "5s",
			Timeout:  "1s",
		},
	}
	if err := this.client.Agent().ServiceRegister(consulService); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/status", StatusHandler)
	fmt.Println("start listen...")
	go http.ListenAndServe(checkAddress, nil)

	log.Printf("Registered service %q in consul with tags %q", service.GetName(), strings.Join(tags, ","))
	return true

}

func (this *ConsulServiceCenter) openListener() {
	t := time.NewTicker(time.Second * 5) //5秒轮询刷新
	for {
		select {
		case <-t.C:
			this.handleListener()
		}
	}
}

func (this *ConsulServiceCenter) handleListener() {
	this.RLock()
	defer this.RUnlock()
	for serviceName, _ := range this.listeners {
		this.DiscoverService(true, serviceName)
	}
}

func (this *ConsulServiceCenter) SubscribeServices(serviceNames ...string) {
	this.Lock()
	defer this.Unlock()
	for _, serviceName := range serviceNames {
		this.listeners[serviceName] = struct{}{}
		this.DiscoverService(true, serviceName)
	}
}

func (this *ConsulServiceCenter) DiscoverService(healthyOnly bool, serviceName string) {
	servicesData, _, err := this.client.Health().Service(serviceName, "", healthyOnly, &api.QueryOptions{})
	if err != nil {
		return
	}

	services := []service.IService{}
	for _, entry := range servicesData {
		serviceEntry := entry.Service
		if serviceName != serviceEntry.Service {
			continue
		}
		for _, health := range entry.Checks {
			fmt.Println("  health nodeid:", health.Node, " serviceName:", health.ServiceName, " service_id:", health.ServiceID, " status:", health.Status)
			fmt.Println("  service id:", serviceEntry.ID, " serviceName:", serviceEntry.Service, " ip:", serviceEntry.Address, " port:", serviceEntry.Port)
			//if health.ServiceName != serviceName {
			//	continue
			//}

			iService, _ := service.NewService1(serviceEntry.ID, serviceEntry.Service, serviceEntry.Address, serviceEntry.Port, service.GRPC)
			services = append(services, iService)
			//node := newService1(health.ServiceID, serviceName, entry.Service.Address, entry.Service.Port, GRPC)
			////get data from kv store
			//s := GetKeyValue(serviceName, node.IP, node.Port)
			//if len(s) > 0 {
			//	var data KVData
			//	err = json.Unmarshal([]byte(s), &data)
			//	if err == nil {
			//		node.Load = data.Load
			//		node.Timestamp = data.Timestamp
			//	}
			//}
			//fmt.Println("service node updated ip:", node.IP, " port:", node.Port, " serviceid:", node.ServiceID, " load:", node.Load, " ts:", node.Timestamp)
			//sers = append(sers, node)
		}
	}
	this.Container.UpdateServices(serviceName, services)

	//service_locker.Lock()
	//servics_map[serviceName] = sers
	//service_locker.Unlock()
}

//func (this *ConsulServiceCenter) AddServiceListener(listener service.Listener) {
//	this.Container.AddServiceListener(listener)
//}
