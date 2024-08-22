package model

type HTTPResponse struct {
	Master string `json:"master"`
	Slave  string `json:"slave"`
}
