package comparator

import (
	"context"
	"github.com/answer1991/discovery/daemon/docker"
	"github.com/answer1991/discovery/store/consul"
	"testing"
)

func TestServiceComparator_Run(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())

	store, err := consul.NewConsulStore("10.244.42.35:8500")

	if nil != err {
		t.Fail()
	}

	query, err := docker.NewDockerServiceDaemon([]string{"io.answer1991.service.name"}, []string{"io.answer1991.service.tags"})

	if nil != err {
		t.Fail()
	}

	comparator := NewServiceComparator(60, query, store)

	comparator.Run(ctx)

	<-make(chan string)
}
