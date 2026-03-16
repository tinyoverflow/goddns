package provider

type Provider interface {
	SetIPAddress(ip string) error
}
