package discovery

type Device struct {
	IpAddress int64  `json:"ipaddress" bson:"ipaddress,omitempty"`
	Name      string `json:"name" bson:"name,omitempty"`
	State     bool   `json:"state" bson:"state,omitempty"`
	Status    string `json:"status" bson:"status,omitempty"`
}
