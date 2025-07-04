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

// DirectLinkInterface Information about the DirectLink interface.
type DirectLinkInterface struct {
	// The BGP (Border Gateway Protocol) ASN (Autonomous System Number) on the customer's side of the DirectLink interface. <br/> This number must be between `1` and `4294967295`, except `50624`, `53306`, and `132418`. <br/> If you do not have an ASN, you can choose one between `64512` and `65534` (both included), or between `4200000000` and `4294967295` (both included).
	BgpAsn int32 `json:"BgpAsn"`
	// The BGP authentication key.
	BgpKey *string `json:"BgpKey,omitempty"`
	// The IP on the customer's side of the DirectLink interface.
	ClientPrivateIp *string `json:"ClientPrivateIp,omitempty"`
	// The name of the DirectLink interface.
	DirectLinkInterfaceName string `json:"DirectLinkInterfaceName"`
	// The IP on the OUTSCALE side of the DirectLink interface.
	OutscalePrivateIp *string `json:"OutscalePrivateIp,omitempty"`
	// The ID of the target virtual gateway.
	VirtualGatewayId string `json:"VirtualGatewayId"`
	// The VLAN number associated with the DirectLink interface. This number must be unique and be between `2` and `4094`.
	Vlan int32 `json:"Vlan"`
}

// NewDirectLinkInterface instantiates a new DirectLinkInterface object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDirectLinkInterface(bgpAsn int32, directLinkInterfaceName string, virtualGatewayId string, vlan int32) *DirectLinkInterface {
	this := DirectLinkInterface{}
	this.BgpAsn = bgpAsn
	this.DirectLinkInterfaceName = directLinkInterfaceName
	this.VirtualGatewayId = virtualGatewayId
	this.Vlan = vlan
	return &this
}

// NewDirectLinkInterfaceWithDefaults instantiates a new DirectLinkInterface object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDirectLinkInterfaceWithDefaults() *DirectLinkInterface {
	this := DirectLinkInterface{}
	return &this
}

// GetBgpAsn returns the BgpAsn field value
func (o *DirectLinkInterface) GetBgpAsn() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.BgpAsn
}

// GetBgpAsnOk returns a tuple with the BgpAsn field value
// and a boolean to check if the value has been set.
func (o *DirectLinkInterface) GetBgpAsnOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.BgpAsn, true
}

// SetBgpAsn sets field value
func (o *DirectLinkInterface) SetBgpAsn(v int32) {
	o.BgpAsn = v
}

// GetBgpKey returns the BgpKey field value if set, zero value otherwise.
func (o *DirectLinkInterface) GetBgpKey() string {
	if o == nil || o.BgpKey == nil {
		var ret string
		return ret
	}
	return *o.BgpKey
}

// GetBgpKeyOk returns a tuple with the BgpKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DirectLinkInterface) GetBgpKeyOk() (*string, bool) {
	if o == nil || o.BgpKey == nil {
		return nil, false
	}
	return o.BgpKey, true
}

// HasBgpKey returns a boolean if a field has been set.
func (o *DirectLinkInterface) HasBgpKey() bool {
	if o != nil && o.BgpKey != nil {
		return true
	}

	return false
}

// SetBgpKey gets a reference to the given string and assigns it to the BgpKey field.
func (o *DirectLinkInterface) SetBgpKey(v string) {
	o.BgpKey = &v
}

// GetClientPrivateIp returns the ClientPrivateIp field value if set, zero value otherwise.
func (o *DirectLinkInterface) GetClientPrivateIp() string {
	if o == nil || o.ClientPrivateIp == nil {
		var ret string
		return ret
	}
	return *o.ClientPrivateIp
}

// GetClientPrivateIpOk returns a tuple with the ClientPrivateIp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DirectLinkInterface) GetClientPrivateIpOk() (*string, bool) {
	if o == nil || o.ClientPrivateIp == nil {
		return nil, false
	}
	return o.ClientPrivateIp, true
}

// HasClientPrivateIp returns a boolean if a field has been set.
func (o *DirectLinkInterface) HasClientPrivateIp() bool {
	if o != nil && o.ClientPrivateIp != nil {
		return true
	}

	return false
}

// SetClientPrivateIp gets a reference to the given string and assigns it to the ClientPrivateIp field.
func (o *DirectLinkInterface) SetClientPrivateIp(v string) {
	o.ClientPrivateIp = &v
}

// GetDirectLinkInterfaceName returns the DirectLinkInterfaceName field value
func (o *DirectLinkInterface) GetDirectLinkInterfaceName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DirectLinkInterfaceName
}

// GetDirectLinkInterfaceNameOk returns a tuple with the DirectLinkInterfaceName field value
// and a boolean to check if the value has been set.
func (o *DirectLinkInterface) GetDirectLinkInterfaceNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DirectLinkInterfaceName, true
}

// SetDirectLinkInterfaceName sets field value
func (o *DirectLinkInterface) SetDirectLinkInterfaceName(v string) {
	o.DirectLinkInterfaceName = v
}

// GetOutscalePrivateIp returns the OutscalePrivateIp field value if set, zero value otherwise.
func (o *DirectLinkInterface) GetOutscalePrivateIp() string {
	if o == nil || o.OutscalePrivateIp == nil {
		var ret string
		return ret
	}
	return *o.OutscalePrivateIp
}

// GetOutscalePrivateIpOk returns a tuple with the OutscalePrivateIp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DirectLinkInterface) GetOutscalePrivateIpOk() (*string, bool) {
	if o == nil || o.OutscalePrivateIp == nil {
		return nil, false
	}
	return o.OutscalePrivateIp, true
}

// HasOutscalePrivateIp returns a boolean if a field has been set.
func (o *DirectLinkInterface) HasOutscalePrivateIp() bool {
	if o != nil && o.OutscalePrivateIp != nil {
		return true
	}

	return false
}

// SetOutscalePrivateIp gets a reference to the given string and assigns it to the OutscalePrivateIp field.
func (o *DirectLinkInterface) SetOutscalePrivateIp(v string) {
	o.OutscalePrivateIp = &v
}

// GetVirtualGatewayId returns the VirtualGatewayId field value
func (o *DirectLinkInterface) GetVirtualGatewayId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.VirtualGatewayId
}

// GetVirtualGatewayIdOk returns a tuple with the VirtualGatewayId field value
// and a boolean to check if the value has been set.
func (o *DirectLinkInterface) GetVirtualGatewayIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.VirtualGatewayId, true
}

// SetVirtualGatewayId sets field value
func (o *DirectLinkInterface) SetVirtualGatewayId(v string) {
	o.VirtualGatewayId = v
}

// GetVlan returns the Vlan field value
func (o *DirectLinkInterface) GetVlan() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Vlan
}

// GetVlanOk returns a tuple with the Vlan field value
// and a boolean to check if the value has been set.
func (o *DirectLinkInterface) GetVlanOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Vlan, true
}

// SetVlan sets field value
func (o *DirectLinkInterface) SetVlan(v int32) {
	o.Vlan = v
}

func (o DirectLinkInterface) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["BgpAsn"] = o.BgpAsn
	}
	if o.BgpKey != nil {
		toSerialize["BgpKey"] = o.BgpKey
	}
	if o.ClientPrivateIp != nil {
		toSerialize["ClientPrivateIp"] = o.ClientPrivateIp
	}
	if true {
		toSerialize["DirectLinkInterfaceName"] = o.DirectLinkInterfaceName
	}
	if o.OutscalePrivateIp != nil {
		toSerialize["OutscalePrivateIp"] = o.OutscalePrivateIp
	}
	if true {
		toSerialize["VirtualGatewayId"] = o.VirtualGatewayId
	}
	if true {
		toSerialize["Vlan"] = o.Vlan
	}
	return json.Marshal(toSerialize)
}

type NullableDirectLinkInterface struct {
	value *DirectLinkInterface
	isSet bool
}

func (v NullableDirectLinkInterface) Get() *DirectLinkInterface {
	return v.value
}

func (v *NullableDirectLinkInterface) Set(val *DirectLinkInterface) {
	v.value = val
	v.isSet = true
}

func (v NullableDirectLinkInterface) IsSet() bool {
	return v.isSet
}

func (v *NullableDirectLinkInterface) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDirectLinkInterface(val *DirectLinkInterface) *NullableDirectLinkInterface {
	return &NullableDirectLinkInterface{value: val, isSet: true}
}

func (v NullableDirectLinkInterface) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDirectLinkInterface) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
