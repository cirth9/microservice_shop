package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type ConsulRegistry struct {
	Host string
	Port int
}
type RegistryClient interface {
	Register(ip string, port int, tags []string, name string, id string) error
	DeRegister(serviceID string) error
}

func NewRegistry(IP string, Port int) RegistryClient {
	return &ConsulRegistry{
		Host: IP,
		Port: Port,
	}
}

func (rc *ConsulRegistry) Register(ip string, port int, tags []string, name string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d",
		rc.Host,
		rc.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", ip, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	registration.Address = ip
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}
	return nil
}

func (rc *ConsulRegistry) DeRegister(serviceID string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", rc.Host, rc.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	err = client.Agent().ServiceDeregister(serviceID)
	return err
}
