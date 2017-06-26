package engine

import (
	"context"
	"github.com/answer1991/discovery/comparator"
	"github.com/answer1991/discovery/daemon/docker"
	"github.com/answer1991/discovery/store/consul"
	"github.com/answer1991/discovery/watcher"
)

func NewEngine(ttl int, addr string, serviceLabelKeys []string, serviceTagLabelKeys []string) (engine *Engine, err error) {
	daemon, err := docker.NewDockerServiceDaemon(serviceLabelKeys, serviceTagLabelKeys)

	if nil != err {
		return nil, err
	}

	store, err := consul.NewConsulStore(addr)

	if nil != err {
		return nil, err
	}

	return &Engine{
		ServiceComparator:   comparator.NewServiceComparator(ttl, daemon, store),
		ServiceEventWatcher: watcher.NewWatcher(daemon, store),
	}, nil
}

type Engine struct {
	*comparator.ServiceComparator
	*watcher.ServiceEventWatcher
}

func (this *Engine) Start(cxt context.Context) {
	this.ServiceComparator.Run(cxt)
	this.ServiceEventWatcher.Watch(cxt)
}
