package consul

import (
	"log"
	"lowkeydd-server/share"
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

	share.JSONFileLoader("setting/consul.json", &setting)

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
