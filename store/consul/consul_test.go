package consul

import (
	"fmt"
	//"github.com/answer1991/registrator/types"
	"testing"
)

func TestServiceConsulStore_List(t *testing.T) {
	client, err := NewConsulStore("10.244.42.35:8500")

	if nil != err {
		t.Fail()
		return
	}

	services, err := client.List()

	if nil != err {
		t.Fail()
		return
	}

	fmt.Println(len(services))

	for _, service := range services {
		fmt.Println(service)
		fmt.Println(service.Id)
	}
}

//func TestServiceConsulStore_Register(t *testing.T) {
//	client, err := NewConsulStore("e010244042035.ztt", "10.244.42.35:8500")
//
//	if nil != err {
//		t.Fail()
//		return
//	}
//
//	id := "eb7d2cff0e24c779899aa702b784856fbc73115d84e44401d2e2658e49483918"
//	uuid := fmt.Sprintf("%s-%s-%s-%s-%s", id[0:8], id[8:12], id[12:16], id[16:20], id[20:32])
//
//	service := &types.Service{
//		Id:          uuid,
//		Node:        "e010244042035.ztt",
//		NodeAddress: "10.244.42.35",
//		Name:        "busybox",
//		Address:     "172.17.0.4",
//		//Ports:       []int{22},
//	}
//
//	err = client.Register(service)
//
//	if nil != err {
//		fmt.Println(err.Error())
//		t.Fail()
//		return
//	}
//}

//func TestServiceConsulStore_Deregister(t *testing.T) {
//	client, err := NewConsulStore("10.244.42.35:8500")
//
//	if nil != err {
//		t.Fail()
//		return
//	}
//
//	id := "eb7d2cff0e24c779899aa702b784856fbc73115d84e44401d2e2658e49483918"
//	uuid := fmt.Sprintf("%s-%s-%s-%s-%s", id[0:8], id[8:12], id[12:16], id[16:20], id[20:32])
//
//	service := &types.Service{
//		Id:          uuid,
//		Node:        "e010244042035.ztt",
//		NodeAddress: "10.244.42.35",
//		Name:        "busybox",
//		Address:     "172.17.0.4",
//		//Ports:       []int{22},
//	}
//
//	err = client.Deregister(service)
//
//	if nil != err {
//		fmt.Println(err.Error())
//		t.Fail()
//		return
//	}
//}
