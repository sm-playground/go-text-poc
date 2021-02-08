package common

type RequestStatus struct {
	Status  string
	Message string
}

type TokenProcessStatus struct {
	Id     int
	Token  string
	Status RequestStatus
}
