package types

type Service struct {
	Id          string
	Node        string
	NodeAddress string
	Name        string
	Address     string
	Ports       []int
	Tags        []string
}

const EventAdd = "Add"
const EventRemove = "Remove"

type ServiceEvent struct {
	*Service
	EventName string
}
