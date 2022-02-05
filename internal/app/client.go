package app

type Client struct {
	tableName struct{} `json:"-" pg:"auth_clients"`

	ID     string `json:"id" pg:"id"`
	Secret string `json:"secret" pg:"secret"`
	Domain string `json:"domain" pg:"domain"`
	UserID string `json:"user_id" pg:"user_id"`
}

// GetID client id
func (c *Client) GetID() string {
	return c.ID
}

// GetSecret client secret
func (c *Client) GetSecret() string {
	return c.Secret
}

// GetDomain client domain
func (c *Client) GetDomain() string {
	return c.Domain
}

// GetUserID user id
func (c *Client) GetUserID() string {
	return c.UserID
}
