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

// With The information to display in each returned log.
type With struct {
	// By default or if set to true, the account ID is displayed.
	AccountId *bool `json:"AccountId,omitempty"`
	// By default or if set to true, the duration of the call is displayed.
	CallDuration *bool `json:"CallDuration,omitempty"`
	// By default or if set to true, the access key is displayed.
	QueryAccessKey *bool `json:"QueryAccessKey,omitempty"`
	// By default or if set to true, the name of the API is displayed.
	QueryApiName *bool `json:"QueryApiName,omitempty"`
	// By default or if set to true, the version of the API is displayed.
	QueryApiVersion *bool `json:"QueryApiVersion,omitempty"`
	// By default or if set to true, the name of the call is displayed.
	QueryCallName *bool `json:"QueryCallName,omitempty"`
	// By default or if set to true, the date of the call is displayed.
	QueryDate *bool `json:"QueryDate,omitempty"`
	// By default or if set to true, the raw header of the HTTP request is displayed.
	QueryHeaderRaw *bool `json:"QueryHeaderRaw,omitempty"`
	// By default or if set to true, the size of the raw header of the HTTP request is displayed.
	QueryHeaderSize *bool `json:"QueryHeaderSize,omitempty"`
	// By default or if set to true, the IP is displayed.
	QueryIpAddress *bool `json:"QueryIpAddress,omitempty"`
	// By default or if set to true, the raw payload of the HTTP request is displayed.
	QueryPayloadRaw *bool `json:"QueryPayloadRaw,omitempty"`
	// By default or if set to true, the size of the raw payload of the HTTP request is displayed.
	QueryPayloadSize *bool `json:"QueryPayloadSize,omitempty"`
	// By default or if set to true, the user agent of the HTTP request is displayed.
	QueryUserAgent *bool `json:"QueryUserAgent,omitempty"`
	// By default or if set to true, the request ID is displayed.
	RequestId *bool `json:"RequestId,omitempty"`
	// By default or if set to true, the size of the response is displayed.
	ResponseSize *bool `json:"ResponseSize,omitempty"`
	// By default or if set to true, the HTTP status code of the response is displayed.
	ResponseStatusCode *bool `json:"ResponseStatusCode,omitempty"`
}

// NewWith instantiates a new With object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWith() *With {
	this := With{}
	var accountId bool = true
	this.AccountId = &accountId
	var callDuration bool = true
	this.CallDuration = &callDuration
	var queryAccessKey bool = true
	this.QueryAccessKey = &queryAccessKey
	var queryApiName bool = true
	this.QueryApiName = &queryApiName
	var queryApiVersion bool = true
	this.QueryApiVersion = &queryApiVersion
	var queryCallName bool = true
	this.QueryCallName = &queryCallName
	var queryDate bool = true
	this.QueryDate = &queryDate
	var queryHeaderRaw bool = true
	this.QueryHeaderRaw = &queryHeaderRaw
	var queryHeaderSize bool = true
	this.QueryHeaderSize = &queryHeaderSize
	var queryIpAddress bool = true
	this.QueryIpAddress = &queryIpAddress
	var queryPayloadRaw bool = true
	this.QueryPayloadRaw = &queryPayloadRaw
	var queryPayloadSize bool = true
	this.QueryPayloadSize = &queryPayloadSize
	var queryUserAgent bool = true
	this.QueryUserAgent = &queryUserAgent
	var requestId bool = true
	this.RequestId = &requestId
	var responseSize bool = true
	this.ResponseSize = &responseSize
	var responseStatusCode bool = true
	this.ResponseStatusCode = &responseStatusCode
	return &this
}

// NewWithWithDefaults instantiates a new With object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWithWithDefaults() *With {
	this := With{}
	var accountId bool = true
	this.AccountId = &accountId
	var callDuration bool = true
	this.CallDuration = &callDuration
	var queryAccessKey bool = true
	this.QueryAccessKey = &queryAccessKey
	var queryApiName bool = true
	this.QueryApiName = &queryApiName
	var queryApiVersion bool = true
	this.QueryApiVersion = &queryApiVersion
	var queryCallName bool = true
	this.QueryCallName = &queryCallName
	var queryDate bool = true
	this.QueryDate = &queryDate
	var queryHeaderRaw bool = true
	this.QueryHeaderRaw = &queryHeaderRaw
	var queryHeaderSize bool = true
	this.QueryHeaderSize = &queryHeaderSize
	var queryIpAddress bool = true
	this.QueryIpAddress = &queryIpAddress
	var queryPayloadRaw bool = true
	this.QueryPayloadRaw = &queryPayloadRaw
	var queryPayloadSize bool = true
	this.QueryPayloadSize = &queryPayloadSize
	var queryUserAgent bool = true
	this.QueryUserAgent = &queryUserAgent
	var requestId bool = true
	this.RequestId = &requestId
	var responseSize bool = true
	this.ResponseSize = &responseSize
	var responseStatusCode bool = true
	this.ResponseStatusCode = &responseStatusCode
	return &this
}

// GetAccountId returns the AccountId field value if set, zero value otherwise.
func (o *With) GetAccountId() bool {
	if o == nil || o.AccountId == nil {
		var ret bool
		return ret
	}
	return *o.AccountId
}

// GetAccountIdOk returns a tuple with the AccountId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *With) GetAccountIdOk() (*bool, bool) {
	if o == nil || o.AccountId == nil {
		return nil, false
	}
	return o.AccountId, true
}

// HasAccountId returns a boolean if a field has been set.
func (o *With) HasAccountId() bool {
	if o != nil && o.AccountId != nil {
		return true
	}

	return false
}

// SetAccountId gets a reference to the given bool and assigns it to the AccountId field.
func (o *With) SetAccountId(v bool) {
	o.AccountId = &v
}

// GetCallDuration returns the CallDuration field value if set, zero value otherwise.
func (o *With) GetCallDuration() bool {
	if o == nil || o.CallDuration == nil {
		var ret bool
		return ret
	}
	return *o.CallDuration
}

// GetCallDurationOk returns a tuple with the CallDuration field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *With) GetCallDurationOk() (*bool, bool) {
	if o == nil || o.CallDuration == nil {
		return nil, false
	}
	return o.CallDuration, true
}

// HasCallDuration returns a boolean if a field has been set.
func (o *With) HasCallDuration() bool {
	if o != nil && o.CallDuration != nil {
		return true
	}

	return false
}

// SetCallDuration gets a reference to the given bool and assigns it to the CallDuration field.
func (o *With) SetCallDuration(v bool) {
	o.CallDuration = &v
}

// GetQueryAccessKey returns the QueryAccessKey field value if set, zero value otherwise.
func (o *With) GetQueryAccessKey() bool {
	if o == nil || o.QueryAccessKey == nil {
		var ret bool
		return ret
	}
	return *o.QueryAccessKey
}

// GetQueryAccessKeyOk returns a tuple with the QueryAccessKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *With) GetQueryAccessKeyOk() (*bool, bool) {
	if o == nil || o.QueryAccessKey == nil {
		return nil, false
	}
	return o.QueryAccessKey, true
}

// HasQueryAccessKey returns a boolean if a field has been set.
func (o *With) HasQueryAccessKey() bool {
	if o != nil && o.QueryAccessKey != nil {
		return true
	}

	return false
}

// SetQueryAccessKey gets a reference to the given bool and assigns it to the QueryAccessKey field.
func (o *With) SetQueryAccessKey(v bool) {
	o.QueryAccessKey = &v
}

// GetQueryApiName returns the QueryApiName field value if set, zero value otherwise.
func (o *With) GetQueryApiName() bool {
	if o == nil || o.QueryApiName == nil {
		var ret bool
		return ret
	}
	return *o.QueryApiName
}

// GetQueryApiNameOk returns a tuple with the QueryApiName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *With) GetQueryApiNameOk() (*bool, bool) {
	if o == nil || o.QueryApiName == nil {
		return nil, false
	}
	return o.QueryApiName, true
}

// HasQueryApiName returns a boolean if a field has been set.
func (o *With) HasQueryApiName() bool {
	if o != nil && o.QueryApiName != nil {
		return true
	}

	return false
}

// SetQueryApiName gets a reference to the given bool and assigns it to the QueryApiName field.
func (o *With) SetQueryApiName(v bool) {
	o.QueryApiName = &v
}

// GetQueryApiVersion returns the QueryApiVersion field value if set, zero value otherwise.
func (o *With) GetQueryApiVersion() bool {
	if o == nil || o.QueryApiVersion == nil {
		var ret bool
		return ret
	}
	return *o.QueryApiVersion
}

// GetQueryApiVersionOk returns a tuple with the QueryApiVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *With) GetQueryApiVersionOk() (*bool, bool) {
	if o == nil || o.QueryApiVersion == nil {
		return nil, false
	}
	return o.QueryApiVersion, true
}

// HasQueryApiVersion returns a boolean if a field has been set.
func (o *With) HasQueryApiVersion() bool {
	if o != nil && o.QueryApiVersion != nil {
		return true
	}

	return false
}

// SetQueryApiVersion gets a reference to the given bool and assigns it to the QueryApiVersion field.
func (o *With) SetQueryApiVersion(v bool) {
	o.QueryApiVersion = &v
}

// GetQueryCallName returns the QueryCallName field value if set, zero value otherwise.
func (o *With) GetQueryCallName() bool {
	if o == nil || o.QueryCallName == nil {
		var ret bool
		return ret
	}
	return *o.QueryCallName
}

// GetQueryCallNameOk returns a tuple with the QueryCallName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *With) GetQueryCallNameOk() (*bool, bool) {
	if o == nil || o.QueryCallName == nil {
		return nil, false
	}
	return o.QueryCallName, true
}

// HasQueryCallName returns a boolean if a field has been set.
func (o *With) HasQueryCallName() bool {
	if o != nil && o.QueryCallName != nil {
		return true
	}

	return false
}

// SetQueryCallName gets a reference to the given bool and assigns it to the QueryCallName field.
func (o *With) SetQueryCallName(v bool) {
	o.QueryCallName = &v
}

// GetQueryDate returns the QueryDate field value if set, zero value otherwise.
func (o *With) GetQueryDate() bool {
	if o == nil || o.QueryDate == nil {
		var ret bool
		return ret
	}
	return *o.QueryDate
}

// GetQueryDateOk returns a tuple with the QueryDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *With) GetQueryDateOk() (*bool, bool) {
	if o == nil || o.QueryDate == nil {
		return nil, false
	}
	return o.QueryDate, true
}

// HasQueryDate returns a boolean if a field has been set.
func (o *With) HasQueryDate() bool {
	if o != nil && o.QueryDate != nil {
		return true
	}

	return false
}

// SetQueryDate gets a reference to the given bool and assigns it to the QueryDate field.
func (o *With) SetQueryDate(v bool) {
	o.QueryDate = &v
}

// GetQueryHeaderRaw returns the QueryHeaderRaw field value if set, zero value otherwise.
func (o *With) GetQueryHeaderRaw() bool {
	if o == nil || o.QueryHeaderRaw == nil {
		var ret bool
		return ret
	}
	return *o.QueryHeaderRaw
}

// GetQueryHeaderRawOk returns a tuple with the QueryHeaderRaw field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *With) GetQueryHeaderRawOk() (*bool, bool) {
	if o == nil || o.QueryHeaderRaw == nil {
		return nil, false
	}
	return o.QueryHeaderRaw, true
}

// HasQueryHeaderRaw returns a boolean if a field has been set.
func (o *With) HasQueryHeaderRaw() bool {
	if o != nil && o.QueryHeaderRaw != nil {
		return true
	}

	return false
}

// SetQueryHeaderRaw gets a reference to the given bool and assigns it to the QueryHeaderRaw field.
func (o *With) SetQueryHeaderRaw(v bool) {
	o.QueryHeaderRaw = &v
}

// GetQueryHeaderSize returns the QueryHeaderSize field value if set, zero value otherwise.
func (o *With) GetQueryHeaderSize() bool {
	if o == nil || o.QueryHeaderSize == nil {
		var ret bool
		return ret
	}
	return *o.QueryHeaderSize
}

// GetQueryHeaderSizeOk returns a tuple with the QueryHeaderSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *With) GetQueryHeaderSizeOk() (*bool, bool) {
	if o == nil || o.QueryHeaderSize == nil {
		return nil, false
	}
	return o.QueryHeaderSize, true
}

// HasQueryHeaderSize returns a boolean if a field has been set.
func (o *With) HasQueryHeaderSize() bool {
	if o != nil && o.QueryHeaderSize != nil {
		return true
	}

	return false
}

// SetQueryHeaderSize gets a reference to the given bool and assigns it to the QueryHeaderSize field.
func (o *With) SetQueryHeaderSize(v bool) {
	o.QueryHeaderSize = &v
}

// GetQueryIpAddress returns the QueryIpAddress field value if set, zero value otherwise.
func (o *With) GetQueryIpAddress() bool {
	if o == nil || o.QueryIpAddress == nil {
		var ret bool
		return ret
	}
	return *o.QueryIpAddress
}

// GetQueryIpAddressOk returns a tuple with the QueryIpAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *With) GetQueryIpAddressOk() (*bool, bool) {
	if o == nil || o.QueryIpAddress == nil {
		return nil, false
	}
	return o.QueryIpAddress, true
}

// HasQueryIpAddress returns a boolean if a field has been set.
func (o *With) HasQueryIpAddress() bool {
	if o != nil && o.QueryIpAddress != nil {
		return true
	}

	return false
}

// SetQueryIpAddress gets a reference to the given bool and assigns it to the QueryIpAddress field.
func (o *With) SetQueryIpAddress(v bool) {
	o.QueryIpAddress = &v
}

// GetQueryPayloadRaw returns the QueryPayloadRaw field value if set, zero value otherwise.
func (o *With) GetQueryPayloadRaw() bool {
	if o == nil || o.QueryPayloadRaw == nil {
		var ret bool
		return ret
	}
	return *o.QueryPayloadRaw
}

// GetQueryPayloadRawOk returns a tuple with the QueryPayloadRaw field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *With) GetQueryPayloadRawOk() (*bool, bool) {
	if o == nil || o.QueryPayloadRaw == nil {
		return nil, false
	}
	return o.QueryPayloadRaw, true
}

// HasQueryPayloadRaw returns a boolean if a field has been set.
func (o *With) HasQueryPayloadRaw() bool {
	if o != nil && o.QueryPayloadRaw != nil {
		return true
	}

	return false
}

// SetQueryPayloadRaw gets a reference to the given bool and assigns it to the QueryPayloadRaw field.
func (o *With) SetQueryPayloadRaw(v bool) {
	o.QueryPayloadRaw = &v
}

// GetQueryPayloadSize returns the QueryPayloadSize field value if set, zero value otherwise.
func (o *With) GetQueryPayloadSize() bool {
	if o == nil || o.QueryPayloadSize == nil {
		var ret bool
		return ret
	}
	return *o.QueryPayloadSize
}

// GetQueryPayloadSizeOk returns a tuple with the QueryPayloadSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *With) GetQueryPayloadSizeOk() (*bool, bool) {
	if o == nil || o.QueryPayloadSize == nil {
		return nil, false
	}
	return o.QueryPayloadSize, true
}

// HasQueryPayloadSize returns a boolean if a field has been set.
func (o *With) HasQueryPayloadSize() bool {
	if o != nil && o.QueryPayloadSize != nil {
		return true
	}

	return false
}

// SetQueryPayloadSize gets a reference to the given bool and assigns it to the QueryPayloadSize field.
func (o *With) SetQueryPayloadSize(v bool) {
	o.QueryPayloadSize = &v
}

// GetQueryUserAgent returns the QueryUserAgent field value if set, zero value otherwise.
func (o *With) GetQueryUserAgent() bool {
	if o == nil || o.QueryUserAgent == nil {
		var ret bool
		return ret
	}
	return *o.QueryUserAgent
}

// GetQueryUserAgentOk returns a tuple with the QueryUserAgent field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *With) GetQueryUserAgentOk() (*bool, bool) {
	if o == nil || o.QueryUserAgent == nil {
		return nil, false
	}
	return o.QueryUserAgent, true
}

// HasQueryUserAgent returns a boolean if a field has been set.
func (o *With) HasQueryUserAgent() bool {
	if o != nil && o.QueryUserAgent != nil {
		return true
	}

	return false
}

// SetQueryUserAgent gets a reference to the given bool and assigns it to the QueryUserAgent field.
func (o *With) SetQueryUserAgent(v bool) {
	o.QueryUserAgent = &v
}

// GetRequestId returns the RequestId field value if set, zero value otherwise.
func (o *With) GetRequestId() bool {
	if o == nil || o.RequestId == nil {
		var ret bool
		return ret
	}
	return *o.RequestId
}

// GetRequestIdOk returns a tuple with the RequestId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *With) GetRequestIdOk() (*bool, bool) {
	if o == nil || o.RequestId == nil {
		return nil, false
	}
	return o.RequestId, true
}

// HasRequestId returns a boolean if a field has been set.
func (o *With) HasRequestId() bool {
	if o != nil && o.RequestId != nil {
		return true
	}

	return false
}

// SetRequestId gets a reference to the given bool and assigns it to the RequestId field.
func (o *With) SetRequestId(v bool) {
	o.RequestId = &v
}

// GetResponseSize returns the ResponseSize field value if set, zero value otherwise.
func (o *With) GetResponseSize() bool {
	if o == nil || o.ResponseSize == nil {
		var ret bool
		return ret
	}
	return *o.ResponseSize
}

// GetResponseSizeOk returns a tuple with the ResponseSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *With) GetResponseSizeOk() (*bool, bool) {
	if o == nil || o.ResponseSize == nil {
		return nil, false
	}
	return o.ResponseSize, true
}

// HasResponseSize returns a boolean if a field has been set.
func (o *With) HasResponseSize() bool {
	if o != nil && o.ResponseSize != nil {
		return true
	}

	return false
}

// SetResponseSize gets a reference to the given bool and assigns it to the ResponseSize field.
func (o *With) SetResponseSize(v bool) {
	o.ResponseSize = &v
}

// GetResponseStatusCode returns the ResponseStatusCode field value if set, zero value otherwise.
func (o *With) GetResponseStatusCode() bool {
	if o == nil || o.ResponseStatusCode == nil {
		var ret bool
		return ret
	}
	return *o.ResponseStatusCode
}

// GetResponseStatusCodeOk returns a tuple with the ResponseStatusCode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *With) GetResponseStatusCodeOk() (*bool, bool) {
	if o == nil || o.ResponseStatusCode == nil {
		return nil, false
	}
	return o.ResponseStatusCode, true
}

// HasResponseStatusCode returns a boolean if a field has been set.
func (o *With) HasResponseStatusCode() bool {
	if o != nil && o.ResponseStatusCode != nil {
		return true
	}

	return false
}

// SetResponseStatusCode gets a reference to the given bool and assigns it to the ResponseStatusCode field.
func (o *With) SetResponseStatusCode(v bool) {
	o.ResponseStatusCode = &v
}

func (o With) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.AccountId != nil {
		toSerialize["AccountId"] = o.AccountId
	}
	if o.CallDuration != nil {
		toSerialize["CallDuration"] = o.CallDuration
	}
	if o.QueryAccessKey != nil {
		toSerialize["QueryAccessKey"] = o.QueryAccessKey
	}
	if o.QueryApiName != nil {
		toSerialize["QueryApiName"] = o.QueryApiName
	}
	if o.QueryApiVersion != nil {
		toSerialize["QueryApiVersion"] = o.QueryApiVersion
	}
	if o.QueryCallName != nil {
		toSerialize["QueryCallName"] = o.QueryCallName
	}
	if o.QueryDate != nil {
		toSerialize["QueryDate"] = o.QueryDate
	}
	if o.QueryHeaderRaw != nil {
		toSerialize["QueryHeaderRaw"] = o.QueryHeaderRaw
	}
	if o.QueryHeaderSize != nil {
		toSerialize["QueryHeaderSize"] = o.QueryHeaderSize
	}
	if o.QueryIpAddress != nil {
		toSerialize["QueryIpAddress"] = o.QueryIpAddress
	}
	if o.QueryPayloadRaw != nil {
		toSerialize["QueryPayloadRaw"] = o.QueryPayloadRaw
	}
	if o.QueryPayloadSize != nil {
		toSerialize["QueryPayloadSize"] = o.QueryPayloadSize
	}
	if o.QueryUserAgent != nil {
		toSerialize["QueryUserAgent"] = o.QueryUserAgent
	}
	if o.RequestId != nil {
		toSerialize["RequestId"] = o.RequestId
	}
	if o.ResponseSize != nil {
		toSerialize["ResponseSize"] = o.ResponseSize
	}
	if o.ResponseStatusCode != nil {
		toSerialize["ResponseStatusCode"] = o.ResponseStatusCode
	}
	return json.Marshal(toSerialize)
}

type NullableWith struct {
	value *With
	isSet bool
}

func (v NullableWith) Get() *With {
	return v.value
}

func (v *NullableWith) Set(val *With) {
	v.value = val
	v.isSet = true
}

func (v NullableWith) IsSet() bool {
	return v.isSet
}

func (v *NullableWith) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWith(val *With) *NullableWith {
	return &NullableWith{value: val, isSet: true}
}

func (v NullableWith) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWith) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
