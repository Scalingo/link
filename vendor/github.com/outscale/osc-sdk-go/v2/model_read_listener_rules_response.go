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

// ReadListenerRulesResponse struct for ReadListenerRulesResponse
type ReadListenerRulesResponse struct {
	// The list of the rules to describe.
	ListenerRules   *[]ListenerRule  `json:"ListenerRules,omitempty"`
	ResponseContext *ResponseContext `json:"ResponseContext,omitempty"`
}

// NewReadListenerRulesResponse instantiates a new ReadListenerRulesResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewReadListenerRulesResponse() *ReadListenerRulesResponse {
	this := ReadListenerRulesResponse{}
	return &this
}

// NewReadListenerRulesResponseWithDefaults instantiates a new ReadListenerRulesResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewReadListenerRulesResponseWithDefaults() *ReadListenerRulesResponse {
	this := ReadListenerRulesResponse{}
	return &this
}

// GetListenerRules returns the ListenerRules field value if set, zero value otherwise.
func (o *ReadListenerRulesResponse) GetListenerRules() []ListenerRule {
	if o == nil || o.ListenerRules == nil {
		var ret []ListenerRule
		return ret
	}
	return *o.ListenerRules
}

// GetListenerRulesOk returns a tuple with the ListenerRules field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ReadListenerRulesResponse) GetListenerRulesOk() (*[]ListenerRule, bool) {
	if o == nil || o.ListenerRules == nil {
		return nil, false
	}
	return o.ListenerRules, true
}

// HasListenerRules returns a boolean if a field has been set.
func (o *ReadListenerRulesResponse) HasListenerRules() bool {
	if o != nil && o.ListenerRules != nil {
		return true
	}

	return false
}

// SetListenerRules gets a reference to the given []ListenerRule and assigns it to the ListenerRules field.
func (o *ReadListenerRulesResponse) SetListenerRules(v []ListenerRule) {
	o.ListenerRules = &v
}

// GetResponseContext returns the ResponseContext field value if set, zero value otherwise.
func (o *ReadListenerRulesResponse) GetResponseContext() ResponseContext {
	if o == nil || o.ResponseContext == nil {
		var ret ResponseContext
		return ret
	}
	return *o.ResponseContext
}

// GetResponseContextOk returns a tuple with the ResponseContext field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ReadListenerRulesResponse) GetResponseContextOk() (*ResponseContext, bool) {
	if o == nil || o.ResponseContext == nil {
		return nil, false
	}
	return o.ResponseContext, true
}

// HasResponseContext returns a boolean if a field has been set.
func (o *ReadListenerRulesResponse) HasResponseContext() bool {
	if o != nil && o.ResponseContext != nil {
		return true
	}

	return false
}

// SetResponseContext gets a reference to the given ResponseContext and assigns it to the ResponseContext field.
func (o *ReadListenerRulesResponse) SetResponseContext(v ResponseContext) {
	o.ResponseContext = &v
}

func (o ReadListenerRulesResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.ListenerRules != nil {
		toSerialize["ListenerRules"] = o.ListenerRules
	}
	if o.ResponseContext != nil {
		toSerialize["ResponseContext"] = o.ResponseContext
	}
	return json.Marshal(toSerialize)
}

type NullableReadListenerRulesResponse struct {
	value *ReadListenerRulesResponse
	isSet bool
}

func (v NullableReadListenerRulesResponse) Get() *ReadListenerRulesResponse {
	return v.value
}

func (v *NullableReadListenerRulesResponse) Set(val *ReadListenerRulesResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableReadListenerRulesResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableReadListenerRulesResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableReadListenerRulesResponse(val *ReadListenerRulesResponse) *NullableReadListenerRulesResponse {
	return &NullableReadListenerRulesResponse{value: val, isSet: true}
}

func (v NullableReadListenerRulesResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableReadListenerRulesResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
