package forms

import "net/url"

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}
