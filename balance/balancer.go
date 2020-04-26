package balance

type Balancer interface {
	BalanceInit()
	GetNode(s string) string
}
