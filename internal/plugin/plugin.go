package plugin

type Provider interface {
	SetIPAddress(ip string) error
}

type Retriever interface {
	GetIPAddress() (string, error)
}
