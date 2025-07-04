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

// FiltersNatService One or more filters.
type FiltersNatService struct {
	// The idempotency tokens provided when creating the NAT services.
	ClientTokens *[]string `json:"ClientTokens,omitempty"`
	// The IDs of the NAT services.
	NatServiceIds *[]string `json:"NatServiceIds,omitempty"`
	// The IDs of the Nets in which the NAT services are.
	NetIds *[]string `json:"NetIds,omitempty"`
	// The states of the NAT services (`pending` \\| `available` \\| `deleting` \\| `deleted`).
	States *[]string `json:"States,omitempty"`
	// The IDs of the Subnets in which the NAT services are.
	SubnetIds *[]string `json:"SubnetIds,omitempty"`
	// The keys of the tags associated with the NAT services.
	TagKeys *[]string `json:"TagKeys,omitempty"`
	// The values of the tags associated with the NAT services.
	TagValues *[]string `json:"TagValues,omitempty"`
	// The key/value combination of the tags associated with the NAT services, in the following format: &quot;Filters&quot;:{&quot;Tags&quot;:[&quot;TAGKEY=TAGVALUE&quot;]}.
	Tags *[]string `json:"Tags,omitempty"`
}

// NewFiltersNatService instantiates a new FiltersNatService object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFiltersNatService() *FiltersNatService {
	this := FiltersNatService{}
	return &this
}

// NewFiltersNatServiceWithDefaults instantiates a new FiltersNatService object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFiltersNatServiceWithDefaults() *FiltersNatService {
	this := FiltersNatService{}
	return &this
}

// GetClientTokens returns the ClientTokens field value if set, zero value otherwise.
func (o *FiltersNatService) GetClientTokens() []string {
	if o == nil || o.ClientTokens == nil {
		var ret []string
		return ret
	}
	return *o.ClientTokens
}

// GetClientTokensOk returns a tuple with the ClientTokens field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersNatService) GetClientTokensOk() (*[]string, bool) {
	if o == nil || o.ClientTokens == nil {
		return nil, false
	}
	return o.ClientTokens, true
}

// HasClientTokens returns a boolean if a field has been set.
func (o *FiltersNatService) HasClientTokens() bool {
	if o != nil && o.ClientTokens != nil {
		return true
	}

	return false
}

// SetClientTokens gets a reference to the given []string and assigns it to the ClientTokens field.
func (o *FiltersNatService) SetClientTokens(v []string) {
	o.ClientTokens = &v
}

// GetNatServiceIds returns the NatServiceIds field value if set, zero value otherwise.
func (o *FiltersNatService) GetNatServiceIds() []string {
	if o == nil || o.NatServiceIds == nil {
		var ret []string
		return ret
	}
	return *o.NatServiceIds
}

// GetNatServiceIdsOk returns a tuple with the NatServiceIds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersNatService) GetNatServiceIdsOk() (*[]string, bool) {
	if o == nil || o.NatServiceIds == nil {
		return nil, false
	}
	return o.NatServiceIds, true
}

// HasNatServiceIds returns a boolean if a field has been set.
func (o *FiltersNatService) HasNatServiceIds() bool {
	if o != nil && o.NatServiceIds != nil {
		return true
	}

	return false
}

// SetNatServiceIds gets a reference to the given []string and assigns it to the NatServiceIds field.
func (o *FiltersNatService) SetNatServiceIds(v []string) {
	o.NatServiceIds = &v
}

// GetNetIds returns the NetIds field value if set, zero value otherwise.
func (o *FiltersNatService) GetNetIds() []string {
	if o == nil || o.NetIds == nil {
		var ret []string
		return ret
	}
	return *o.NetIds
}

// GetNetIdsOk returns a tuple with the NetIds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersNatService) GetNetIdsOk() (*[]string, bool) {
	if o == nil || o.NetIds == nil {
		return nil, false
	}
	return o.NetIds, true
}

// HasNetIds returns a boolean if a field has been set.
func (o *FiltersNatService) HasNetIds() bool {
	if o != nil && o.NetIds != nil {
		return true
	}

	return false
}

// SetNetIds gets a reference to the given []string and assigns it to the NetIds field.
func (o *FiltersNatService) SetNetIds(v []string) {
	o.NetIds = &v
}

// GetStates returns the States field value if set, zero value otherwise.
func (o *FiltersNatService) GetStates() []string {
	if o == nil || o.States == nil {
		var ret []string
		return ret
	}
	return *o.States
}

// GetStatesOk returns a tuple with the States field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersNatService) GetStatesOk() (*[]string, bool) {
	if o == nil || o.States == nil {
		return nil, false
	}
	return o.States, true
}

// HasStates returns a boolean if a field has been set.
func (o *FiltersNatService) HasStates() bool {
	if o != nil && o.States != nil {
		return true
	}

	return false
}

// SetStates gets a reference to the given []string and assigns it to the States field.
func (o *FiltersNatService) SetStates(v []string) {
	o.States = &v
}

// GetSubnetIds returns the SubnetIds field value if set, zero value otherwise.
func (o *FiltersNatService) GetSubnetIds() []string {
	if o == nil || o.SubnetIds == nil {
		var ret []string
		return ret
	}
	return *o.SubnetIds
}

// GetSubnetIdsOk returns a tuple with the SubnetIds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersNatService) GetSubnetIdsOk() (*[]string, bool) {
	if o == nil || o.SubnetIds == nil {
		return nil, false
	}
	return o.SubnetIds, true
}

// HasSubnetIds returns a boolean if a field has been set.
func (o *FiltersNatService) HasSubnetIds() bool {
	if o != nil && o.SubnetIds != nil {
		return true
	}

	return false
}

// SetSubnetIds gets a reference to the given []string and assigns it to the SubnetIds field.
func (o *FiltersNatService) SetSubnetIds(v []string) {
	o.SubnetIds = &v
}

// GetTagKeys returns the TagKeys field value if set, zero value otherwise.
func (o *FiltersNatService) GetTagKeys() []string {
	if o == nil || o.TagKeys == nil {
		var ret []string
		return ret
	}
	return *o.TagKeys
}

// GetTagKeysOk returns a tuple with the TagKeys field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersNatService) GetTagKeysOk() (*[]string, bool) {
	if o == nil || o.TagKeys == nil {
		return nil, false
	}
	return o.TagKeys, true
}

// HasTagKeys returns a boolean if a field has been set.
func (o *FiltersNatService) HasTagKeys() bool {
	if o != nil && o.TagKeys != nil {
		return true
	}

	return false
}

// SetTagKeys gets a reference to the given []string and assigns it to the TagKeys field.
func (o *FiltersNatService) SetTagKeys(v []string) {
	o.TagKeys = &v
}

// GetTagValues returns the TagValues field value if set, zero value otherwise.
func (o *FiltersNatService) GetTagValues() []string {
	if o == nil || o.TagValues == nil {
		var ret []string
		return ret
	}
	return *o.TagValues
}

// GetTagValuesOk returns a tuple with the TagValues field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersNatService) GetTagValuesOk() (*[]string, bool) {
	if o == nil || o.TagValues == nil {
		return nil, false
	}
	return o.TagValues, true
}

// HasTagValues returns a boolean if a field has been set.
func (o *FiltersNatService) HasTagValues() bool {
	if o != nil && o.TagValues != nil {
		return true
	}

	return false
}

// SetTagValues gets a reference to the given []string and assigns it to the TagValues field.
func (o *FiltersNatService) SetTagValues(v []string) {
	o.TagValues = &v
}

// GetTags returns the Tags field value if set, zero value otherwise.
func (o *FiltersNatService) GetTags() []string {
	if o == nil || o.Tags == nil {
		var ret []string
		return ret
	}
	return *o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersNatService) GetTagsOk() (*[]string, bool) {
	if o == nil || o.Tags == nil {
		return nil, false
	}
	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *FiltersNatService) HasTags() bool {
	if o != nil && o.Tags != nil {
		return true
	}

	return false
}

// SetTags gets a reference to the given []string and assigns it to the Tags field.
func (o *FiltersNatService) SetTags(v []string) {
	o.Tags = &v
}

func (o FiltersNatService) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.ClientTokens != nil {
		toSerialize["ClientTokens"] = o.ClientTokens
	}
	if o.NatServiceIds != nil {
		toSerialize["NatServiceIds"] = o.NatServiceIds
	}
	if o.NetIds != nil {
		toSerialize["NetIds"] = o.NetIds
	}
	if o.States != nil {
		toSerialize["States"] = o.States
	}
	if o.SubnetIds != nil {
		toSerialize["SubnetIds"] = o.SubnetIds
	}
	if o.TagKeys != nil {
		toSerialize["TagKeys"] = o.TagKeys
	}
	if o.TagValues != nil {
		toSerialize["TagValues"] = o.TagValues
	}
	if o.Tags != nil {
		toSerialize["Tags"] = o.Tags
	}
	return json.Marshal(toSerialize)
}

type NullableFiltersNatService struct {
	value *FiltersNatService
	isSet bool
}

func (v NullableFiltersNatService) Get() *FiltersNatService {
	return v.value
}

func (v *NullableFiltersNatService) Set(val *FiltersNatService) {
	v.value = val
	v.isSet = true
}

func (v NullableFiltersNatService) IsSet() bool {
	return v.isSet
}

func (v *NullableFiltersNatService) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFiltersNatService(val *FiltersNatService) *NullableFiltersNatService {
	return &NullableFiltersNatService{value: val, isSet: true}
}

func (v NullableFiltersNatService) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFiltersNatService) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
