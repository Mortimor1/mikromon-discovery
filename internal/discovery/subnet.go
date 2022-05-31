package discovery

type Subnet struct {
	Subnet string `json:"subnet" bson:"subnet,omitempty"`
	State  bool   `json:"state" bson:"state,omitempty"`
}
