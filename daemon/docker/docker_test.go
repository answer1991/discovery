package docker

import (
	"context"
	"fmt"
	"testing"
)

//func TestDockerQuery_GetEventCh(t *testing.T) {
//	ctx, _ := context.WithCancel(context.Background())
//	query, err := NewDockerQuery(ctx)
//
//	if nil != err {
//		t.Fail()
//	}
//
//	mCh, eCh := query.GetEventCh()
//
//	for {
//		select {
//		case event := <-mCh:
//			fmt.Println(event)
//			fmt.Println(*event.Service)
//		case <-eCh:
//			return
//		}
//	}
//}

func TestDockerQuery_ListService(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())
	daemon, err := NewDockerServiceDaemon([]string{"io.answer1991.service.name"}, []string{"io.answer1991.service.tags"})

	if nil != err {
		t.Fail()
	}

	services, err := daemon.ListService(ctx)

	if nil != err {
		t.Fail()
	}

	for _, service := range services {
		fmt.Println(*service)
	}
}
