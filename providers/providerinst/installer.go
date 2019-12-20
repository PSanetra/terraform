package providerinst

import (
	"fmt"

	svchost "github.com/hashicorp/terraform-svchost"
	disco "github.com/hashicorp/terraform-svchost/disco"
)

// Installer is a type that knows how to find and install providers from
// provider registries.
type Installer struct {
	services *disco.Disco
}

// NewInstaller creates and returns a new installer that will use the given
// service discovery object to perform discovery for provider registry
// services.
func NewInstaller(services *disco.Disco) *Installer {
	return &Installer{
		services: services,
	}
}

func (i *Installer) registryClient(hostname svchost.Hostname) (*registryClient, error) {
	host, err := i.services.Discover(hostname)
	if err != nil {
		return nil, ErrHostUnreachable{
			Hostname: hostname,
			Wrapped:  err,
		}
	}

	url, err := host.ServiceURL("providers.v1")
	switch err := err.(type) {
	case nil:
		// okay! We'll fall through and return below.
	case *disco.ErrServiceNotProvided:
		return nil, ErrHostNoProviders{
			Hostname: hostname,
		}
	case *disco.ErrVersionNotSupported:
		return nil, ErrHostNoProviders{
			Hostname:        hostname,
			HasOtherVersion: true,
		}
	default:
		return nil, ErrHostUnreachable{
			Hostname: hostname,
			Wrapped:  err,
		}
	}

	// Check if we have credentials configured for this hostname.
	creds, err := i.services.CredentialsForHost(hostname)
	if err != nil {
		// This indicates that a credentials helper failed, which means we
		// can't do anything better than just pass through the helper's
		// own error message.
		return nil, fmt.Errorf("failed to retrieve credentials for %s: %s", hostname, err)
	}

	return &registryClient{
		baseURL: url,
		creds:   creds,
	}, nil
}
