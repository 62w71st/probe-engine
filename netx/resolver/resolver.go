package resolver

import (
	"context"
)

// Resolver is a DNS resolver. The *net.Resolver used by Go implements
// this interface, but other implementations are possible.
type Resolver interface {
	// LookupHost resolves a hostname to a list of IP addresses.
	LookupHost(ctx context.Context, hostname string) (addrs []string, err error)
}
