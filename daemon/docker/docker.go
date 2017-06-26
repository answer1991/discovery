package docker

import (
	"context"
	"fmt"
	"github.com/answer1991/daily-roll-logrus"
	"github.com/answer1991/discovery/types"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
	"strings"
)

var logger = drl.GetLogger("docker-daemon")

func NewDockerServiceDaemon(serviceLabelKeys []string, serviceTagLabelKeys []string) (daemon *DockerServiceDaemon, err error) {
	dockerClient, err := client.NewEnvClient()
	if nil != err {
		return nil, err
	}

	return &DockerServiceDaemon{
		Client:              dockerClient,
		ServiceLabelKeys:    serviceLabelKeys,
		ServiceTagLabelKeys: serviceTagLabelKeys,
	}, nil
}

type DockerServiceDaemon struct {
	*client.Client
	ServiceLabelKeys    []string
	ServiceTagLabelKeys []string
}

func (this *DockerServiceDaemon) convertToService(ctx context.Context, containerId string) (service *types.Service, err error) {
	info, err := this.Client.ContainerInspect(ctx, containerId)

	if nil != err {
		return nil, err
	}

	serviceName := getServiceName(info, this.ServiceLabelKeys)

	if "" == serviceName {
		logger.
			WithField("containerName", info.Name).
			WithField("containerId", info.ID).
			WithField("containerLabels", info.Config.Labels).
			Info("Container Can not Found Service Name, Will Ignore")

		return nil, nil
	}

	return &types.Service{
		Id:      convertToServiceId(info.ID),
		Name:    serviceName,
		Address: info.NetworkSettings.IPAddress,
		Ports:   getServicePorts(info),
		Tags:    getServiceTags(info, this.ServiceTagLabelKeys),
	}, nil
}

func convertToServiceId(containerId string) (serviceId string) {
	return fmt.Sprintf("%s-%s-%s-%s-%s", containerId[0:8], containerId[8:12], containerId[12:16], containerId[16:20], containerId[20:32])
}

func getServiceName(container dockerTypes.ContainerJSON, serviceLabelKeys []string) (serviceName string) {
	for _, serviceLabelKey := range serviceLabelKeys {
		if value, ok := container.Config.Labels[serviceLabelKey]; ok {
			return value
		}
	}
	return ""
}

func getServiceTags(container dockerTypes.ContainerJSON, serviceTagLabelKeys []string) (tags []string) {
	tags = make([]string, 0)
	for _, tagKey := range serviceTagLabelKeys {
		if value, ok := container.Config.Labels[tagKey]; ok {
			tagArr := strings.Split(value, ",")
			tags = append(tags, tagArr...)
		}
	}
	return tags
}

func getServicePorts(container dockerTypes.ContainerJSON) (ports []int) {
	return []int{}
}

func (this *DockerServiceDaemon) processMessage(cxt context.Context, message *events.Message, eventCh chan *types.ServiceEvent) (err error) {
	logger.WithField("message", message).Debug("Get Message")
	if message.Type == events.ContainerEventType {
		logger.WithField("message", message).Debug("Get Container Message")

		if "start" == message.Action {
			logger.WithField("message", message).Debug("Get Container Start Message")
			containerId := message.Actor.ID
			service, err := this.convertToService(cxt, containerId)

			if nil != err {
				logger.WithField("err", err).Info("Get Container Info Failed")
			}

			if nil == service {
				logger.Info("Ignore Service which can not convert to a Service")
				return nil
			}

			eventCh <- &types.ServiceEvent{
				EventName: types.EventAdd,
				Service:   service,
			}

		} else if "die" == message.Action {
			logger.WithField("message", message).Debug("Get Container Start Message")
			containerId := message.Actor.ID

			eventCh <- &types.ServiceEvent{
				EventName: types.EventAdd,
				Service: &types.Service{
					Id: convertToServiceId(containerId),
				},
			}
		}
	}
	return nil
}

func (this *DockerServiceDaemon) GetEventCh(ctx context.Context) (eventCh chan *types.ServiceEvent, errCh chan error) {
	mch, ech := this.Client.Events(ctx, dockerTypes.EventsOptions{})

	eventCh = make(chan *types.ServiceEvent)
	errCh = make(chan error)

	go func() {
		defer close(errCh)
		defer close(eventCh)

		for {
			select {
			case message := <-mch:
				this.processMessage(ctx, &message, eventCh)
			case err := <-ech:
				logger.WithField("err", err).Info("Event ErrorCh Get Error")
				mch, ech = this.Client.Events(ctx, dockerTypes.EventsOptions{})
			case <-ctx.Done():
				return
			}
		}
	}()

	return eventCh, errCh
}

func (this *DockerServiceDaemon) ListService(ctx context.Context) (services []*types.Service, err error) {
	containers, err := this.Client.ContainerList(ctx, dockerTypes.ContainerListOptions{})

	if nil != err {
		return nil, err
	}

	services = make([]*types.Service, 0)
	for _, container := range containers {
		if service, err := this.convertToService(ctx, container.ID); nil == err {
			if nil != service {
				services = append(services, service)
			}
		}
	}

	return services, nil
}
