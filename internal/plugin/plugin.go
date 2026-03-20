package plugin

// Provider updates a DNS record to a given IP address.
type Provider interface {
	SetIPAddress(ip string) error
}

// Retriever fetches the current public IP address.
type Retriever interface {
	GetIPAddress() (string, error)
}
