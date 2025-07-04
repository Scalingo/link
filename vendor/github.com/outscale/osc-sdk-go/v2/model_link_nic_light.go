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

// LinkNicLight Information about the network interface card (NIC).
type LinkNicLight struct {
	// If true, the NIC is deleted when the VM is terminated.
	DeleteOnVmDeletion *bool `json:"DeleteOnVmDeletion,omitempty"`
	// The device index for the NIC attachment (between `1` and `7`, both included).
	DeviceNumber *int32 `json:"DeviceNumber,omitempty"`
	// The ID of the NIC to attach.
	LinkNicId *string `json:"LinkNicId,omitempty"`
	// The state of the attachment (`attaching` \\| `attached` \\| `detaching` \\| `detached`).
	State *string `json:"State,omitempty"`
}

// NewLinkNicLight instantiates a new LinkNicLight object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLinkNicLight() *LinkNicLight {
	this := LinkNicLight{}
	return &this
}

// NewLinkNicLightWithDefaults instantiates a new LinkNicLight object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLinkNicLightWithDefaults() *LinkNicLight {
	this := LinkNicLight{}
	return &this
}

// GetDeleteOnVmDeletion returns the DeleteOnVmDeletion field value if set, zero value otherwise.
func (o *LinkNicLight) GetDeleteOnVmDeletion() bool {
	if o == nil || o.DeleteOnVmDeletion == nil {
		var ret bool
		return ret
	}
	return *o.DeleteOnVmDeletion
}

// GetDeleteOnVmDeletionOk returns a tuple with the DeleteOnVmDeletion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinkNicLight) GetDeleteOnVmDeletionOk() (*bool, bool) {
	if o == nil || o.DeleteOnVmDeletion == nil {
		return nil, false
	}
	return o.DeleteOnVmDeletion, true
}

// HasDeleteOnVmDeletion returns a boolean if a field has been set.
func (o *LinkNicLight) HasDeleteOnVmDeletion() bool {
	if o != nil && o.DeleteOnVmDeletion != nil {
		return true
	}

	return false
}

// SetDeleteOnVmDeletion gets a reference to the given bool and assigns it to the DeleteOnVmDeletion field.
func (o *LinkNicLight) SetDeleteOnVmDeletion(v bool) {
	o.DeleteOnVmDeletion = &v
}

// GetDeviceNumber returns the DeviceNumber field value if set, zero value otherwise.
func (o *LinkNicLight) GetDeviceNumber() int32 {
	if o == nil || o.DeviceNumber == nil {
		var ret int32
		return ret
	}
	return *o.DeviceNumber
}

// GetDeviceNumberOk returns a tuple with the DeviceNumber field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinkNicLight) GetDeviceNumberOk() (*int32, bool) {
	if o == nil || o.DeviceNumber == nil {
		return nil, false
	}
	return o.DeviceNumber, true
}

// HasDeviceNumber returns a boolean if a field has been set.
func (o *LinkNicLight) HasDeviceNumber() bool {
	if o != nil && o.DeviceNumber != nil {
		return true
	}

	return false
}

// SetDeviceNumber gets a reference to the given int32 and assigns it to the DeviceNumber field.
func (o *LinkNicLight) SetDeviceNumber(v int32) {
	o.DeviceNumber = &v
}

// GetLinkNicId returns the LinkNicId field value if set, zero value otherwise.
func (o *LinkNicLight) GetLinkNicId() string {
	if o == nil || o.LinkNicId == nil {
		var ret string
		return ret
	}
	return *o.LinkNicId
}

// GetLinkNicIdOk returns a tuple with the LinkNicId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinkNicLight) GetLinkNicIdOk() (*string, bool) {
	if o == nil || o.LinkNicId == nil {
		return nil, false
	}
	return o.LinkNicId, true
}

// HasLinkNicId returns a boolean if a field has been set.
func (o *LinkNicLight) HasLinkNicId() bool {
	if o != nil && o.LinkNicId != nil {
		return true
	}

	return false
}

// SetLinkNicId gets a reference to the given string and assigns it to the LinkNicId field.
func (o *LinkNicLight) SetLinkNicId(v string) {
	o.LinkNicId = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *LinkNicLight) GetState() string {
	if o == nil || o.State == nil {
		var ret string
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinkNicLight) GetStateOk() (*string, bool) {
	if o == nil || o.State == nil {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *LinkNicLight) HasState() bool {
	if o != nil && o.State != nil {
		return true
	}

	return false
}

// SetState gets a reference to the given string and assigns it to the State field.
func (o *LinkNicLight) SetState(v string) {
	o.State = &v
}

func (o LinkNicLight) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.DeleteOnVmDeletion != nil {
		toSerialize["DeleteOnVmDeletion"] = o.DeleteOnVmDeletion
	}
	if o.DeviceNumber != nil {
		toSerialize["DeviceNumber"] = o.DeviceNumber
	}
	if o.LinkNicId != nil {
		toSerialize["LinkNicId"] = o.LinkNicId
	}
	if o.State != nil {
		toSerialize["State"] = o.State
	}
	return json.Marshal(toSerialize)
}

type NullableLinkNicLight struct {
	value *LinkNicLight
	isSet bool
}

func (v NullableLinkNicLight) Get() *LinkNicLight {
	return v.value
}

func (v *NullableLinkNicLight) Set(val *LinkNicLight) {
	v.value = val
	v.isSet = true
}

func (v NullableLinkNicLight) IsSet() bool {
	return v.isSet
}

func (v *NullableLinkNicLight) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLinkNicLight(val *LinkNicLight) *NullableLinkNicLight {
	return &NullableLinkNicLight{value: val, isSet: true}
}

func (v NullableLinkNicLight) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLinkNicLight) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
