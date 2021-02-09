package common

type RequestStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type TokenProcessStatus struct {
	Id     int           `json:"id"`
	Token  string        `json:"token"`
	Status RequestStatus `json:"result"`
}
