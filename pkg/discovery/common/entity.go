package common

// Entity holds the information of host entity.
type Entity struct {
	IPAddress string   `json:"ip_address"`
	Namespace string   `json:"namespace"`
	Services  []string `json:"service"`
	Prefix    string   `json:"prefix"`
	UID       string   `json:"uid"`
}
