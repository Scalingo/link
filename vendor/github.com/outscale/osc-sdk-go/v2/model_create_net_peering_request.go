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

// CreateNetPeeringRequest struct for CreateNetPeeringRequest
type CreateNetPeeringRequest struct {
	// The ID of the Net you want to connect with. <br/ > <br/ > If the Net does not belong to you, you must also specify the `AccepterOwnerId` parameter with the account ID owning the Net you want to connect with.
	AccepterNetId string `json:"AccepterNetId"`
	// The account ID of the owner of the Net you want to connect with. By default, the account ID of the owner of the Net from which the peering request is sent. <br /><br/ > This parameter is required if the Net you want to connect with does not belong to you.
	AccepterOwnerId *string `json:"AccepterOwnerId,omitempty"`
	// If true, checks whether you have the required permissions to perform the action.
	DryRun *bool `json:"DryRun,omitempty"`
	// The ID of the Net you send the peering request from.
	SourceNetId string `json:"SourceNetId"`
}

// NewCreateNetPeeringRequest instantiates a new CreateNetPeeringRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateNetPeeringRequest(accepterNetId string, sourceNetId string) *CreateNetPeeringRequest {
	this := CreateNetPeeringRequest{}
	this.AccepterNetId = accepterNetId
	this.SourceNetId = sourceNetId
	return &this
}

// NewCreateNetPeeringRequestWithDefaults instantiates a new CreateNetPeeringRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateNetPeeringRequestWithDefaults() *CreateNetPeeringRequest {
	this := CreateNetPeeringRequest{}
	return &this
}

// GetAccepterNetId returns the AccepterNetId field value
func (o *CreateNetPeeringRequest) GetAccepterNetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.AccepterNetId
}

// GetAccepterNetIdOk returns a tuple with the AccepterNetId field value
// and a boolean to check if the value has been set.
func (o *CreateNetPeeringRequest) GetAccepterNetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.AccepterNetId, true
}

// SetAccepterNetId sets field value
func (o *CreateNetPeeringRequest) SetAccepterNetId(v string) {
	o.AccepterNetId = v
}

// GetAccepterOwnerId returns the AccepterOwnerId field value if set, zero value otherwise.
func (o *CreateNetPeeringRequest) GetAccepterOwnerId() string {
	if o == nil || o.AccepterOwnerId == nil {
		var ret string
		return ret
	}
	return *o.AccepterOwnerId
}

// GetAccepterOwnerIdOk returns a tuple with the AccepterOwnerId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateNetPeeringRequest) GetAccepterOwnerIdOk() (*string, bool) {
	if o == nil || o.AccepterOwnerId == nil {
		return nil, false
	}
	return o.AccepterOwnerId, true
}

// HasAccepterOwnerId returns a boolean if a field has been set.
func (o *CreateNetPeeringRequest) HasAccepterOwnerId() bool {
	if o != nil && o.AccepterOwnerId != nil {
		return true
	}

	return false
}

// SetAccepterOwnerId gets a reference to the given string and assigns it to the AccepterOwnerId field.
func (o *CreateNetPeeringRequest) SetAccepterOwnerId(v string) {
	o.AccepterOwnerId = &v
}

// GetDryRun returns the DryRun field value if set, zero value otherwise.
func (o *CreateNetPeeringRequest) GetDryRun() bool {
	if o == nil || o.DryRun == nil {
		var ret bool
		return ret
	}
	return *o.DryRun
}

// GetDryRunOk returns a tuple with the DryRun field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateNetPeeringRequest) GetDryRunOk() (*bool, bool) {
	if o == nil || o.DryRun == nil {
		return nil, false
	}
	return o.DryRun, true
}

// HasDryRun returns a boolean if a field has been set.
func (o *CreateNetPeeringRequest) HasDryRun() bool {
	if o != nil && o.DryRun != nil {
		return true
	}

	return false
}

// SetDryRun gets a reference to the given bool and assigns it to the DryRun field.
func (o *CreateNetPeeringRequest) SetDryRun(v bool) {
	o.DryRun = &v
}

// GetSourceNetId returns the SourceNetId field value
func (o *CreateNetPeeringRequest) GetSourceNetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.SourceNetId
}

// GetSourceNetIdOk returns a tuple with the SourceNetId field value
// and a boolean to check if the value has been set.
func (o *CreateNetPeeringRequest) GetSourceNetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SourceNetId, true
}

// SetSourceNetId sets field value
func (o *CreateNetPeeringRequest) SetSourceNetId(v string) {
	o.SourceNetId = v
}

func (o CreateNetPeeringRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["AccepterNetId"] = o.AccepterNetId
	}
	if o.AccepterOwnerId != nil {
		toSerialize["AccepterOwnerId"] = o.AccepterOwnerId
	}
	if o.DryRun != nil {
		toSerialize["DryRun"] = o.DryRun
	}
	if true {
		toSerialize["SourceNetId"] = o.SourceNetId
	}
	return json.Marshal(toSerialize)
}

type NullableCreateNetPeeringRequest struct {
	value *CreateNetPeeringRequest
	isSet bool
}

func (v NullableCreateNetPeeringRequest) Get() *CreateNetPeeringRequest {
	return v.value
}

func (v *NullableCreateNetPeeringRequest) Set(val *CreateNetPeeringRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateNetPeeringRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateNetPeeringRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateNetPeeringRequest(val *CreateNetPeeringRequest) *NullableCreateNetPeeringRequest {
	return &NullableCreateNetPeeringRequest{value: val, isSet: true}
}

func (v NullableCreateNetPeeringRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateNetPeeringRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
