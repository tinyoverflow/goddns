package retriever

type Retriever interface {
	GetIPAddress() (string, error)
}
