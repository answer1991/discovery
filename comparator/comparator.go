package comparator

import (
	"context"
	"github.com/answer1991/daily-roll-logrus"
	"github.com/answer1991/discovery/daemon"
	"github.com/answer1991/discovery/store"
	"github.com/answer1991/discovery/types"
	"time"
)

var logger = drl.GetLogger("comparator")

func NewServiceComparator(ttl int, serviceDaemon daemon.ServiceDaemon, store store.ServiceStore) *ServiceComparator {
	return &ServiceComparator{
		ttl:           ttl,
		serviceDaemon: serviceDaemon,
		store:         store,
	}
}

type ServiceComparator struct {
	ttl           int
	serviceDaemon daemon.ServiceDaemon
	store         store.ServiceStore
}

func (this *ServiceComparator) compare(ctx context.Context) (removeServices []*types.Service, appendServices []*types.Service, err error) {
	expectServices, err := this.store.List()

	if nil != err {
		return nil, nil, err
	}

	//logger.WithField("Service in Store", expectServices).Info("Get Services in Store")

	actualServices, err := this.serviceDaemon.ListService(ctx)

	if nil != err {
		return nil, nil, err
	}

	removeServices = make([]*types.Service, 0)
	appendServices = make([]*types.Service, 0)

	for _, expectService := range expectServices {
		exist := false

		for _, actualService := range actualServices {
			if expectService.Id == actualService.Id {
				exist = true
			}
		}

		if !exist {
			removeServices = append(removeServices, expectService)
		}
	}

	for _, actualService := range actualServices {
		exist := false

		for _, expectService := range expectServices {
			if expectService.Id == actualService.Id {
				exist = true
			}
		}

		if !exist {
			appendServices = append(appendServices, actualService)
		}
	}

	return removeServices, appendServices, err
}

func (this *ServiceComparator) runOnce(ctx context.Context) {
	removes, appends, err := this.compare(ctx)

	if nil != err {

	}

	for _, remove := range removes {
		err := this.store.Deregister(remove)

		if nil != err {
			logger.WithField("err", err).WithField("service", *remove).Error("Deregister Failed")
		} else {
			logger.WithField("service", remove).Info("Deregister Success")
		}
	}

	for _, append := range appends {
		err := this.store.Register(append)

		if nil != err {
			logger.WithField("err", err).WithField("service", *append).Error("Deregister Failed")
		} else {
			logger.WithField("service", append).Info("Register Success")
		}
	}
}

func (this *ServiceComparator) Run(ctx context.Context) {
	go func() {
		this.runOnce(ctx)
		ticker := time.NewTicker(time.Second * time.Duration(this.ttl))

		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				this.runOnce(ctx)
			}
		}
	}()
}
