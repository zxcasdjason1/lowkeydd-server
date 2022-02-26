package consul

import (
	"fmt"
	"log"
	"lowkeydd-server/share"
	"os"
	"strconv"
	"strings"
	"sync"

	consul_api "github.com/hashicorp/consul/api"
)

type Service = consul_api.AgentServiceRegistration

type Setting struct {
	Server struct {
		IP   string `json:"IP"`
		Port string `json:"Port"`
	} `json:"server"`
	Services []Service `json:"services"`
}

type Driver struct {
	Services []Service
	client   *consul_api.Client
}

var (
	lock    = &sync.Mutex{}
	driver  *Driver
	setting Setting
)

func GetInstance() *Driver {
	if driver == nil {
		// 只允許一個goroutine訪問
		lock.Lock()
		defer lock.Unlock()
		if driver == nil {
			NewDriver()
		}
	}
	return driver
}

func NewDriver() {

	serviceIP := os.Getenv("SERVICE_IP")
	servicePort := os.Getenv("SERVICE_PORT")
	serviceType := strings.ToUpper(os.Getenv("SERVICE_TYPE"))

	if serviceType == "SERVER" {
		share.JSONFileLoader("setting/consul_server.json", &setting)
	} else {
		share.JSONFileLoader("setting/consul_client.json", &setting)
	}
	// 如果有傳入 IP 跟 PORT 就會重設，沒有就使用原先json檔裡的設定。
	if serviceIP != "" {
		setting.Server.IP = serviceIP
		setting.Services[0].Address = serviceIP
		log.Printf("[Consul] SERVICE_IP :> %s \n", serviceIP)
	}
	if servicePort != "" {
		// 設置該service的port，這是提供consul去檢查的資訊
		if port, err := strconv.Atoi(servicePort); err != nil && port > 0 {
			setting.Services[0].Port = port
		}
		// 設置該service的address
		addr := fmt.Sprintf("%s:%s", setting.Server.IP, servicePort)
		setting.Services[0].Address = addr
		setting.Services[0].Checks[0].HTTP = fmt.Sprintf("http://%s/health", addr)
		if serviceType == "SERVER" {
			setting.Services[0].Checks[1].HTTP = fmt.Sprintf("http://%s/crawler/update/", addr)
		}
		log.Printf("[Consul] SERVICE_PORT :> %s \n", servicePort)
	}

	// setup driver
	driver = &Driver{
		Services: setting.Services,
	}

	// Get a new client
	conf := consul_api.DefaultConfig()
	conf.Address = setting.Server.IP + ":" + setting.Server.Port

	// 創建consul Client
	client, err := consul_api.NewClient(conf)
	if err != nil {
		log.Fatal(err)
	}
	driver.client = client

}

func (this *Driver) RegisterService() {

	for _, service := range this.Services {

		log.Printf("RegisterService ID:> %s", service.ID)
		// log.Printf("Name:> %s", service.Name)
		// log.Printf("Tags:> %v", service.Tags)
		// log.Printf("Address:> %s", service.Address)
		// log.Printf("Port:> %d", service.Port)
		// log.Printf("Check:> %v", service.Check)

		err := this.client.Agent().ServiceRegister(&service)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func (this *Driver) KillService() {

	for _, service := range this.Services {

		log.Printf("DeRegisterService ID:> %s", service.ID)

		err := this.client.Agent().ServiceDeregister(service.ID)

		if err != nil {
			log.Fatal(err)
		}
	}
}
