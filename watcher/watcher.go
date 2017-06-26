package watcher

import (
	"context"
	"github.com/answer1991/discovery/daemon"
	"github.com/answer1991/discovery/store"
	"github.com/answer1991/discovery/types"
)

func NewWatcher(serviceQuery daemon.ServiceDaemon, store store.ServiceStore) (watcher *ServiceEventWatcher) {
	return &ServiceEventWatcher{
		serviceDaemon: serviceQuery,
		serviceStore:  store,
	}
}

type ServiceEventWatcher struct {
	serviceDaemon daemon.ServiceDaemon
	serviceStore  store.ServiceStore
}

func (this *ServiceEventWatcher) Watch(ctx context.Context) {
	go func() {
		mCh, eCh := this.serviceDaemon.GetEventCh(ctx)

		for {
			select {
			case event := <-mCh:
				switch event.EventName {
				case types.EventAdd:
					this.serviceStore.Register(event.Service)
					break
				case types.EventRemove:
					this.serviceStore.Deregister(event.Service)
					break
				default:
					break
				}
			case <-eCh:
				mCh, eCh = this.serviceDaemon.GetEventCh(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()
}
