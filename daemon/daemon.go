package daemon

import (
	"context"
	"github.com/answer1991/discovery/types"
)

type ServiceDaemon interface {
	GetEventCh(ctx context.Context) (eventCh chan *types.ServiceEvent, errCh chan error)
	ListService(ctx context.Context) (services []*types.Service, err error)
}
