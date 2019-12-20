package providerinst

import (
	"fmt"

	"github.com/apparentlymart/go-versions/versions"

	"github.com/hashicorp/terraform/addrs"
)

// AvailableVersions returns all of the versions available for the provider
// with the given address, or an error if that result cannot be determined.
//
// If the request fails, the returned error might be an value of
// ErrHostNoProviders, ErrHostUnreachable, ErrUnauthenticated,
// ErrProviderNotKnown, or ErrQueryFailed. Callers must be defensive and
// expect errors of other types too, to allow for future expansion.
func (i *Installer) AvailableVersions(provider addrs.Provider) (versions.List, error) {
	client, err := i.registryClient(provider.Hostname)
	if err != nil {
		return nil, err
	}

	versionStrs, err := client.ProviderVersions(provider)
	if err != nil {
		return nil, err
	}

	if len(versionStrs) == 0 {
		return nil, nil
	}

	ret := make(versions.List, len(versionStrs))
	for i, str := range versionStrs {
		v, err := versions.ParseVersion(str)
		if err != nil {
			return nil, ErrQueryFailed{
				Provider: provider,
				Wrapped:  fmt.Errorf("registry response includes invalid version string %q: %s", str, err),
			}
		}
		ret[i] = v
	}
	ret.Sort() // lowest precedence first, preserving order when equal precedence
	return ret, nil
}
