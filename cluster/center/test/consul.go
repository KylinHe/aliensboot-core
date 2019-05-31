/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/6/1
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package test

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ServiceInfo struct {
	ServiceID string
	IP        string
	Port      int
	Load      int
	Timestamp int //load updated ts
}
type ServiceList []ServiceInfo

type KVData struct {
	Load      int `json:"load"`
	Timestamp int `json:"ts"`
}

var (
	servics_map     = make(map[string]ServiceList)
	service_locker  = new(sync.Mutex)
	consul_client   *api.Client
	my_service_id   string
	my_service_name string
	my_kv_key       string
)

func CheckErr(err error) {
	if err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("check status.")
	fmt.Fprint(w, "status ok!")
}

func StartService(addr string) {
	http.HandleFunc("/status", StatusHandler)
	fmt.Println("start listen...")
	err := http.ListenAndServe(addr, nil)
	CheckErr(err)
}

func main() {
	var status_monitor_addr, service_name, service_ip, consul_addr, found_service string
	var service_port int
	flag.StringVar(&consul_addr, "consul_addr", "localhost:8500", "host:port of the service stuats monitor interface")
	flag.StringVar(&status_monitor_addr, "monitor_addr", "127.0.0.1:54321", "host:port of the service stuats monitor interface")
	flag.StringVar(&service_name, "service_name", "worker", "name of the service")
	flag.StringVar(&service_ip, "ip", "127.0.0.1", "service serve ip")
	flag.StringVar(&found_service, "found_service", "worker", "found the target service")
	flag.IntVar(&service_port, "port", 4300, "service serve port")
	flag.Parse()

	my_service_name = service_name

	DoRegistService(consul_addr, status_monitor_addr, service_name, service_ip, service_port)

	go DoDiscover(consul_addr, found_service)

	go StartService(status_monitor_addr)

	go WaitToUnRegistService()

	go DoUpdateKeyValue(consul_addr, service_name, service_ip, service_port)

	select {}
}

func DoRegistService(consul_addr string, monitor_addr string, service_name string, ip string, port int) {
	my_service_id = service_name + "-" + ip
	var tags []string
	service := &api.AgentServiceRegistration{
		ID:      my_service_id,
		Name:    service_name,
		Port:    port,
		Address: ip,
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			HTTP:     "http://" + monitor_addr + "/status",
			Interval: "5s",
			Timeout:  "1s",
		},
	}

	client, err := api.NewClient(api.DefaultConfig())

	if err != nil {
		log.Fatal(err)
	}
	consul_client = client
	if err := consul_client.Agent().ServiceRegister(service); err != nil {
		log.Fatal(err)
	}
	log.Printf("Registered service %q in consul with tags %q", service_name, strings.Join(tags, ","))
}

func WaitToUnRegistService() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	if consul_client == nil {
		return
	}
	if err := consul_client.Agent().ServiceDeregister(my_service_id); err != nil {
		log.Fatal(err)
	}
}

func DoDiscover(consul_addr string, found_service string) {
	t := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-t.C:
			DiscoverServices(consul_addr, true, found_service)
		}
	}
}

func DiscoverServices(addr string, healthyOnly bool, service_name string) {
	consulConf := api.DefaultConfig()
	consulConf.Address = addr
	client, err := api.NewClient(consulConf)
	CheckErr(err)

	services, _, err := client.Catalog().Services(&api.QueryOptions{})
	CheckErr(err)

	fmt.Println("--do discover ---:", addr)

	var sers ServiceList
	for name := range services {
		servicesData, _, err := client.Health().Service(name, "", healthyOnly,
			&api.QueryOptions{})
		CheckErr(err)
		for _, entry := range servicesData {
			if service_name != entry.Service.Service {
				continue
			}
			for _, health := range entry.Checks {
				if health.ServiceName != service_name {
					continue
				}
				fmt.Println("  health nodeid:", health.Node, " service_name:", health.ServiceName, " service_id:", health.ServiceID, " status:", health.Status, " ip:", entry.Service.Address, " port:", entry.Service.Port)

				var node ServiceInfo
				node.IP = entry.Service.Address
				node.Port = entry.Service.Port
				node.ServiceID = health.ServiceID

				//get data from kv store
				s := GetKeyValue(service_name, node.IP, node.Port)
				if len(s) > 0 {
					var data KVData
					err = json.Unmarshal([]byte(s), &data)
					if err == nil {
						node.Load = data.Load
						node.Timestamp = data.Timestamp
					}
				}
				fmt.Println("service node updated ip:", node.IP, " port:", node.Port, " serviceid:", node.ServiceID, " load:", node.Load, " ts:", node.Timestamp)
				sers = append(sers, node)
			}
		}
	}

	service_locker.Lock()
	servics_map[service_name] = sers
	service_locker.Unlock()
}

func DoUpdateKeyValue(consul_addr string, service_name string, ip string, port int) {
	t := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-t.C:
			StoreKeyValue(consul_addr, service_name, ip, port)
		}
	}
}

func StoreKeyValue(consul_addr string, service_name string, ip string, port int) {

	my_kv_key = my_service_name + "/" + ip + ":" + strconv.Itoa(port)

	var data KVData
	data.Load = rand.Intn(100)
	data.Timestamp = int(time.Now().Unix())
	bys, _ := json.Marshal(&data)

	kv := &api.KVPair{
		Key:   my_kv_key,
		Flags: 0,
		Value: bys,
	}

	_, err := consul_client.KV().Put(kv, nil)
	CheckErr(err)
	fmt.Println(" store data key:", kv.Key, " value:", string(bys))
}

func GetKeyValue(service_name string, ip string, port int) string {
	key := service_name + "/" + ip + ":" + strconv.Itoa(port)

	kv, _, err := consul_client.KV().Get(key, nil)
	if kv == nil {
		return ""
	}
	CheckErr(err)

	return string(kv.Value)
}
