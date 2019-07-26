package mapi

// TextValue is a key value pair in a struct
type Provider struct {
	ID   string `json:"Value"`
	Name string `json:"Text"`
}

func (c *Client) GetProviders() ([]Provider, error) {
	var providers []Provider
	err := c.get("nextrip/providers?format=json", &providers)
	return providers, err
}
