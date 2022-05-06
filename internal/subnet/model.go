package subnet

type Subnet struct {
	Id     string `json:"id" bson:"_id,omitempty"`
	Subnet string `json:"subnet" bson:"subnet,omitempty"`
	State  bool   `json:"state" bson:"state,omitempty"`
}
