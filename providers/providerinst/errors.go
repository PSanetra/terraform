package providerinst

import (
	"fmt"

	"github.com/apparentlymart/go-versions/versions"

	svchost "github.com/hashicorp/terraform-svchost"
	"github.com/hashicorp/terraform/addrs"
)

// ErrHostNoProviders is an error type used to indicate that a hostname given
// in a provider address does not support the provider registry protocol.
type ErrHostNoProviders struct {
	Hostname svchost.Hostname

	// HasOtherVersionis set to true if the discovery process detected
	// declarations of services named "providers" whose version numbers did not
	// match any version supported by the current version of Terraform.
	//
	// If this is set, it's helpful to hint to the user in an error message
	// that the provider host may be expecting an older or a newer version
	// of Terraform, rather than that it isn't a provider registry host at all.
	HasOtherVersion bool
}

func (err ErrHostNoProviders) Error() string {
	switch {
	case err.HasOtherVersion:
		return fmt.Sprintf("host %s does not support the provider registry protocol required by this Terraform version, but may be compatible with a different Terraform version", err.Hostname.ForDisplay())
	default:
		return fmt.Sprintf("host %s does not offer a Terraform provider registry", err.Hostname.ForDisplay())
	}
}

// ErrHostUnreachable is an error type used to indicate that a hostname
// given in a provider address did not resolve in DNS, did not respond to an
// HTTPS request for service discovery, or otherwise failed to correctly speak
// the service discovery protocol.
type ErrHostUnreachable struct {
	Hostname svchost.Hostname
	Wrapped  error
}

func (err ErrHostUnreachable) Error() string {
	return fmt.Sprintf("could not connect to %s: %s", err.Hostname.ForDisplay(), err.Wrapped.Error())
}

// Unwrap returns the underlying error that occurred when trying to reach the
// indicated host.
func (err ErrHostUnreachable) Unwrap() error {
	return err.Wrapped
}

// ErrUnauthorized is an error type used to indicate that a hostname
// given in a provider address returned a "401 Unauthorized" or "403 Forbidden"
// error response when we tried to access it.
type ErrUnauthorized struct {
	Hostname svchost.Hostname

	// HaveCredentials is true when the request that failed included some
	// credentials, and thus it seems that those credentials were invalid.
	// Conversely, HaveCredentials is false if the request did not include
	// credentials at all, in which case it seems that credentials must be
	// provided.
	HaveCredentials bool
}

func (err ErrUnauthorized) Error() string {
	switch {
	case err.HaveCredentials:
		return fmt.Sprintf("host %s rejected the given authentication credentials", err.Hostname)
	default:
		return fmt.Sprintf("host %s requires authentication credentials", err.Hostname)
	}
}

// ErrProviderNotKnown is an error type used to indicate that the hostname
// given in a provider address does appear to be a provider registry but that
// registry does not know about the given provider namespace or type.
//
// A caller serving requests from an end-user should recognize this error type
// and use it to produce user-friendly hints for common errors such as failing
// to specify an explicit source for a provider not in the default namespace
// (one not under registry.terraform.io/hashicorp/). The default error message
// for this type is a direct description of the problem with no such hints,
// because we expect that the caller will have better context to decide what
// hints are appropriate, e.g. by looking at the configuration given by the
// user.
type ErrProviderNotKnown struct {
	Provider addrs.Provider
}

func (err ErrProviderNotKnown) Error() string {
	return fmt.Sprintf(
		"provider registry %s does not have a provider named %s",
		err.Provider.Hostname.ForDisplay(),
		err.Provider,
	)
}

// ErrQueryFailed is an error type used to indicate that the hostname given
// in a provider address does appear to be a provider registry but that when
// we queried it for metadata for the given provider the server returned an
// unexpected error.
//
// This is used for any error responses other than "Not Found", which would
// indicate the absense of a provider and is thus reported using
// ErrProviderNotKnown instead.
type ErrQueryFailed struct {
	Provider addrs.Provider
	Wrapped  error
}

func (err ErrQueryFailed) Error() string {
	return fmt.Sprintf(
		"could not query provider registry for %s: %s",
		err.Provider.String(),
		err.Wrapped.Error(),
	)
}

// Unwrap returns the underlying error that occurred when trying to reach the
// indicated host.
func (err ErrQueryFailed) Unwrap() error {
	return err.Wrapped
}

// ErrDownloadFailed is an error type used to indicate that a specific provider
// version was successfully chosen but that there was an error when trying to
// download the distribution archive for that version.
type ErrDownloadFailed struct {
	Provider addrs.Provider
	Version  versions.Version
	Wrapped  error
}

func (err ErrDownloadFailed) Error() string {
	return fmt.Sprintf(
		"failed to download %s %s: %s",
		err.Provider.String(),
		err.Version.String(),
		err.Wrapped.Error(),
	)
}

// Unwrap returns the underlying error that occurred when trying to reach the
// indicated host.
func (err ErrDownloadFailed) Unwrap() error {
	return err.Wrapped
}
