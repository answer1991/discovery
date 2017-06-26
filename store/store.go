package store

import "github.com/answer1991/discovery/types"

type ServiceStore interface {
	Register(service *types.Service) (err error)
	Deregister(service *types.Service) (err error)
	List() (services []*types.Service, err error)
}
