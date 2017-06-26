package consul

import (
	"github.com/answer1991/daily-roll-logrus"
	"github.com/answer1991/discovery/types"
	"github.com/hashicorp/consul/api"
)

var logger = drl.GetLogger("store-consul")

const (
	serviceConsul = "consul"
)

func NewConsulStore(address string) (store *ServiceConsulStore, err error) {
	client, err := api.NewClient(&api.Config{
		Address: address,
		Scheme:  "http",
	})

	if nil != err {
		return nil, err
	}

	info, err := client.Agent().Self()

	if nil != err {
		return nil, err
	}

	return &ServiceConsulStore{
		NodeName:    info["Config"]["NodeName"].(string),
		NodeId:      info["Config"]["NodeID"].(string),
		NodeAddress: info["Member"]["Addr"].(string),
		Client:      client,
	}, nil
}

type ServiceConsulStore struct {
	NodeName    string
	NodeId      string
	NodeAddress string
	*api.Client
}

func (this *ServiceConsulStore) List() (services []*types.Service, err error) {
	consulServices, err := this.Client.Agent().Services()

	if nil != err {
		return nil, err
	}

	services = make([]*types.Service, 0)

	for _, consulService := range consulServices {
		if serviceConsul == consulService.ID {
			continue
		}

		services = append(services, &types.Service{
			Id:          consulService.ID,
			Node:        this.NodeName,
			NodeAddress: this.NodeAddress,
			Name:        consulService.Service,
			Address:     consulService.Address,
			Ports:       []int{consulService.Port},
			Tags:        consulService.Tags,
		})
	}

	return services, nil
}

func (this *ServiceConsulStore) Register(service *types.Service) (err error) {
	port := 0

	if nil != service.Ports && len(service.Ports) > 0 {
		port = service.Ports[0]
	}

	reg := &api.AgentServiceRegistration{
		ID:      service.Id,
		Name:    service.Name,
		Tags:    service.Tags,
		Port:    port,
		Address: service.Address,
	}

	logger.WithField("service", *service).Info("Register Service")

	return this.Client.Agent().ServiceRegister(reg)
}

func (this *ServiceConsulStore) Deregister(service *types.Service) (err error) {
	if serviceConsul == service.Id {
		return nil
	}

	logger.WithField("serviceId", service.Id).Info("Deregister Service")
	return this.Client.Agent().ServiceDeregister(service.Id)
}
