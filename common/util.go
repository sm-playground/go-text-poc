package common

import c "github.com/sm-playground/go-text-poc/config"

type RequestStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type TokenProcessStatus struct {
	Id     int           `json:"id"`
	Token  string        `json:"token"`
	Status RequestStatus `json:"result"`
}

// GetServiceOwnerId - returns the value of service owner Id from configuration file
func GetServiceOwnerId() string {
	config := c.GetInstance().Get()

	return config.ServiceOwnerSourceId
}
