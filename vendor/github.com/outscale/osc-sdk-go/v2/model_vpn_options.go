/*
 * 3DS OUTSCALE API
 *
 * Welcome to the OUTSCALE API documentation.<br /> The OUTSCALE API enables you to manage your resources in the OUTSCALE Cloud. This documentation describes the different actions available along with code examples.<br /><br /> Throttling: To protect against overloads, the number of identical requests allowed in a given time period is limited.<br /> Brute force: To protect against brute force attacks, the number of failed authentication attempts in a given time period is limited.<br /><br /> Note that the OUTSCALE Cloud is compatible with Amazon Web Services (AWS) APIs, but there are [differences in resource names](https://docs.outscale.com/en/userguide/About-the-APIs.html) between AWS and the OUTSCALE API.<br /> You can also manage your resources using the [Cockpit](https://docs.outscale.com/en/userguide/About-Cockpit.html) web interface.<br /><br /> An OpenAPI description of the OUTSCALE API is also available in this [GitHub repository](https://github.com/outscale/osc-api).<br /> # Authentication Schemes ### Access Key/Secret Key The main way to authenticate your requests to the OUTSCALE API is to use an access key and a secret key.<br /> The mechanism behind this is based on AWS Signature Version 4, whose technical implementation details are described in [Signature of API Requests](https://docs.outscale.com/en/userguide/Signature-of-API-Requests.html).<br /><br /> In practice, the way to specify your access key and secret key depends on the tool or SDK you want to use to interact with the API.<br />  > For example, if you use OSC CLI: > 1. You need to create an **~/.osc/config.json** file to specify your access key, secret key, and the Region of your account. > 2. You then specify the `--profile` option when executing OSC CLI commands. > > For more information, see [Installing and Configuring OSC CLI](https://docs.outscale.com/en/userguide/Installing-and-Configuring-OSC-CLI.html).  See the code samples in each section of this documentation for specific examples in different programming languages.<br /> For more information about access keys, see [About Access Keys](https://docs.outscale.com/en/userguide/About-Access-Keys.html).  > If you try to sign requests with an invalid access key four times in a row, further authentication attempts will be prevented for 1 minute. This lockout time increases 1 minute every four failed attempts, for up to 10 minutes.  ### Login/Password For certain API actions, you can also use basic authentication with the login (email address) and password of your TINA account.<br /> This is useful only in special circumstances, for example if you do not know your access key/secret key and want to retrieve them programmatically.<br /> In most cases, however, you can use the Cockpit web interface to retrieve them.<br />  > For example, if you use OSC CLI: > 1. You need to create an **~/.osc/config.json** file to specify the Region of your account, but you leave the access key value and secret key value empty (`&quot;&quot;`). > 2. You then specify the `--profile`, `--authentication-method`, `--login`, and `--password` options when executing OSC CLI commands.  See the code samples in each section of this documentation for specific examples in different programming languages.  > If you try to sign requests with an invalid password four times in a row, further authentication attempts will be prevented for 1 minute. This lockout time increases 1 minute every four failed attempts, for up to 10 minutes.  ### No Authentication A few API actions do not require any authentication. They are indicated as such in this documentation.<br /> ### Other Security Mechanisms In parallel with the authentication schemes, you can add other security mechanisms to your OUTSCALE account, for example to restrict API requests by IP or other criteria.<br /> For more information, see [Managing Your API Accesses](https://docs.outscale.com/en/userguide/Managing-Your-API-Accesses.html). # Pagination Tutorial You can learn more about the pagination methods for read calls in the dedicated [pagination tutorial](https://docs.outscale.com/en/userguide/Tutorial-Paginating-an-API-Request.html). # Error Codes Reference You can learn more about errors returned by the API in the dedicated [errors page](api-errors.html).
 *
 * API version: 1.35.3
 * Contact: support@outscale.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package osc

import (
	"encoding/json"
)

// VpnOptions Information about the VPN options.
type VpnOptions struct {
	Phase1Options *Phase1Options `json:"Phase1Options,omitempty"`
	Phase2Options *Phase2Options `json:"Phase2Options,omitempty"`
	// The range of inside IPs for the tunnel. This must be a /30 CIDR block from the 169.254.254.0/24 range.
	TunnelInsideIpRange *string `json:"TunnelInsideIpRange,omitempty"`
}

// NewVpnOptions instantiates a new VpnOptions object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewVpnOptions() *VpnOptions {
	this := VpnOptions{}
	return &this
}

// NewVpnOptionsWithDefaults instantiates a new VpnOptions object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewVpnOptionsWithDefaults() *VpnOptions {
	this := VpnOptions{}
	return &this
}

// GetPhase1Options returns the Phase1Options field value if set, zero value otherwise.
func (o *VpnOptions) GetPhase1Options() Phase1Options {
	if o == nil || o.Phase1Options == nil {
		var ret Phase1Options
		return ret
	}
	return *o.Phase1Options
}

// GetPhase1OptionsOk returns a tuple with the Phase1Options field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VpnOptions) GetPhase1OptionsOk() (*Phase1Options, bool) {
	if o == nil || o.Phase1Options == nil {
		return nil, false
	}
	return o.Phase1Options, true
}

// HasPhase1Options returns a boolean if a field has been set.
func (o *VpnOptions) HasPhase1Options() bool {
	if o != nil && o.Phase1Options != nil {
		return true
	}

	return false
}

// SetPhase1Options gets a reference to the given Phase1Options and assigns it to the Phase1Options field.
func (o *VpnOptions) SetPhase1Options(v Phase1Options) {
	o.Phase1Options = &v
}

// GetPhase2Options returns the Phase2Options field value if set, zero value otherwise.
func (o *VpnOptions) GetPhase2Options() Phase2Options {
	if o == nil || o.Phase2Options == nil {
		var ret Phase2Options
		return ret
	}
	return *o.Phase2Options
}

// GetPhase2OptionsOk returns a tuple with the Phase2Options field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VpnOptions) GetPhase2OptionsOk() (*Phase2Options, bool) {
	if o == nil || o.Phase2Options == nil {
		return nil, false
	}
	return o.Phase2Options, true
}

// HasPhase2Options returns a boolean if a field has been set.
func (o *VpnOptions) HasPhase2Options() bool {
	if o != nil && o.Phase2Options != nil {
		return true
	}

	return false
}

// SetPhase2Options gets a reference to the given Phase2Options and assigns it to the Phase2Options field.
func (o *VpnOptions) SetPhase2Options(v Phase2Options) {
	o.Phase2Options = &v
}

// GetTunnelInsideIpRange returns the TunnelInsideIpRange field value if set, zero value otherwise.
func (o *VpnOptions) GetTunnelInsideIpRange() string {
	if o == nil || o.TunnelInsideIpRange == nil {
		var ret string
		return ret
	}
	return *o.TunnelInsideIpRange
}

// GetTunnelInsideIpRangeOk returns a tuple with the TunnelInsideIpRange field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VpnOptions) GetTunnelInsideIpRangeOk() (*string, bool) {
	if o == nil || o.TunnelInsideIpRange == nil {
		return nil, false
	}
	return o.TunnelInsideIpRange, true
}

// HasTunnelInsideIpRange returns a boolean if a field has been set.
func (o *VpnOptions) HasTunnelInsideIpRange() bool {
	if o != nil && o.TunnelInsideIpRange != nil {
		return true
	}

	return false
}

// SetTunnelInsideIpRange gets a reference to the given string and assigns it to the TunnelInsideIpRange field.
func (o *VpnOptions) SetTunnelInsideIpRange(v string) {
	o.TunnelInsideIpRange = &v
}

func (o VpnOptions) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Phase1Options != nil {
		toSerialize["Phase1Options"] = o.Phase1Options
	}
	if o.Phase2Options != nil {
		toSerialize["Phase2Options"] = o.Phase2Options
	}
	if o.TunnelInsideIpRange != nil {
		toSerialize["TunnelInsideIpRange"] = o.TunnelInsideIpRange
	}
	return json.Marshal(toSerialize)
}

type NullableVpnOptions struct {
	value *VpnOptions
	isSet bool
}

func (v NullableVpnOptions) Get() *VpnOptions {
	return v.value
}

func (v *NullableVpnOptions) Set(val *VpnOptions) {
	v.value = val
	v.isSet = true
}

func (v NullableVpnOptions) IsSet() bool {
	return v.isSet
}

func (v *NullableVpnOptions) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableVpnOptions(val *VpnOptions) *NullableVpnOptions {
	return &NullableVpnOptions{value: val, isSet: true}
}

func (v NullableVpnOptions) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableVpnOptions) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
