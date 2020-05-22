package api

// Option it's a interface for options func.
type Option func(*Client)

// OptDoer it's option for set http client.
func OptDoer(d doer) Option {
	return func(client *Client) {
		client.http = d
	}
}
