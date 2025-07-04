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

// FiltersPublicIp One or more filters.
type FiltersPublicIp struct {
	// The IDs representing the associations of public IPs with VMs or NICs.
	LinkPublicIpIds *[]string `json:"LinkPublicIpIds,omitempty"`
	// The account IDs of the owners of the NICs.
	NicAccountIds *[]string `json:"NicAccountIds,omitempty"`
	// The IDs of the NICs.
	NicIds *[]string `json:"NicIds,omitempty"`
	// Whether the public IPs are for use in the public Cloud or in a Net.
	Placements *[]string `json:"Placements,omitempty"`
	// The private IPs associated with the public IPs.
	PrivateIps *[]string `json:"PrivateIps,omitempty"`
	// The IDs of the public IPs.
	PublicIpIds *[]string `json:"PublicIpIds,omitempty"`
	// The public IPs.
	PublicIps *[]string `json:"PublicIps,omitempty"`
	// The keys of the tags associated with the public IPs.
	TagKeys *[]string `json:"TagKeys,omitempty"`
	// The values of the tags associated with the public IPs.
	TagValues *[]string `json:"TagValues,omitempty"`
	// The key/value combination of the tags associated with the public IPs, in the following format: &quot;Filters&quot;:{&quot;Tags&quot;:[&quot;TAGKEY=TAGVALUE&quot;]}.
	Tags *[]string `json:"Tags,omitempty"`
	// The IDs of the VMs.
	VmIds *[]string `json:"VmIds,omitempty"`
}

// NewFiltersPublicIp instantiates a new FiltersPublicIp object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFiltersPublicIp() *FiltersPublicIp {
	this := FiltersPublicIp{}
	return &this
}

// NewFiltersPublicIpWithDefaults instantiates a new FiltersPublicIp object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFiltersPublicIpWithDefaults() *FiltersPublicIp {
	this := FiltersPublicIp{}
	return &this
}

// GetLinkPublicIpIds returns the LinkPublicIpIds field value if set, zero value otherwise.
func (o *FiltersPublicIp) GetLinkPublicIpIds() []string {
	if o == nil || o.LinkPublicIpIds == nil {
		var ret []string
		return ret
	}
	return *o.LinkPublicIpIds
}

// GetLinkPublicIpIdsOk returns a tuple with the LinkPublicIpIds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersPublicIp) GetLinkPublicIpIdsOk() (*[]string, bool) {
	if o == nil || o.LinkPublicIpIds == nil {
		return nil, false
	}
	return o.LinkPublicIpIds, true
}

// HasLinkPublicIpIds returns a boolean if a field has been set.
func (o *FiltersPublicIp) HasLinkPublicIpIds() bool {
	if o != nil && o.LinkPublicIpIds != nil {
		return true
	}

	return false
}

// SetLinkPublicIpIds gets a reference to the given []string and assigns it to the LinkPublicIpIds field.
func (o *FiltersPublicIp) SetLinkPublicIpIds(v []string) {
	o.LinkPublicIpIds = &v
}

// GetNicAccountIds returns the NicAccountIds field value if set, zero value otherwise.
func (o *FiltersPublicIp) GetNicAccountIds() []string {
	if o == nil || o.NicAccountIds == nil {
		var ret []string
		return ret
	}
	return *o.NicAccountIds
}

// GetNicAccountIdsOk returns a tuple with the NicAccountIds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersPublicIp) GetNicAccountIdsOk() (*[]string, bool) {
	if o == nil || o.NicAccountIds == nil {
		return nil, false
	}
	return o.NicAccountIds, true
}

// HasNicAccountIds returns a boolean if a field has been set.
func (o *FiltersPublicIp) HasNicAccountIds() bool {
	if o != nil && o.NicAccountIds != nil {
		return true
	}

	return false
}

// SetNicAccountIds gets a reference to the given []string and assigns it to the NicAccountIds field.
func (o *FiltersPublicIp) SetNicAccountIds(v []string) {
	o.NicAccountIds = &v
}

// GetNicIds returns the NicIds field value if set, zero value otherwise.
func (o *FiltersPublicIp) GetNicIds() []string {
	if o == nil || o.NicIds == nil {
		var ret []string
		return ret
	}
	return *o.NicIds
}

// GetNicIdsOk returns a tuple with the NicIds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersPublicIp) GetNicIdsOk() (*[]string, bool) {
	if o == nil || o.NicIds == nil {
		return nil, false
	}
	return o.NicIds, true
}

// HasNicIds returns a boolean if a field has been set.
func (o *FiltersPublicIp) HasNicIds() bool {
	if o != nil && o.NicIds != nil {
		return true
	}

	return false
}

// SetNicIds gets a reference to the given []string and assigns it to the NicIds field.
func (o *FiltersPublicIp) SetNicIds(v []string) {
	o.NicIds = &v
}

// GetPlacements returns the Placements field value if set, zero value otherwise.
func (o *FiltersPublicIp) GetPlacements() []string {
	if o == nil || o.Placements == nil {
		var ret []string
		return ret
	}
	return *o.Placements
}

// GetPlacementsOk returns a tuple with the Placements field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersPublicIp) GetPlacementsOk() (*[]string, bool) {
	if o == nil || o.Placements == nil {
		return nil, false
	}
	return o.Placements, true
}

// HasPlacements returns a boolean if a field has been set.
func (o *FiltersPublicIp) HasPlacements() bool {
	if o != nil && o.Placements != nil {
		return true
	}

	return false
}

// SetPlacements gets a reference to the given []string and assigns it to the Placements field.
func (o *FiltersPublicIp) SetPlacements(v []string) {
	o.Placements = &v
}

// GetPrivateIps returns the PrivateIps field value if set, zero value otherwise.
func (o *FiltersPublicIp) GetPrivateIps() []string {
	if o == nil || o.PrivateIps == nil {
		var ret []string
		return ret
	}
	return *o.PrivateIps
}

// GetPrivateIpsOk returns a tuple with the PrivateIps field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersPublicIp) GetPrivateIpsOk() (*[]string, bool) {
	if o == nil || o.PrivateIps == nil {
		return nil, false
	}
	return o.PrivateIps, true
}

// HasPrivateIps returns a boolean if a field has been set.
func (o *FiltersPublicIp) HasPrivateIps() bool {
	if o != nil && o.PrivateIps != nil {
		return true
	}

	return false
}

// SetPrivateIps gets a reference to the given []string and assigns it to the PrivateIps field.
func (o *FiltersPublicIp) SetPrivateIps(v []string) {
	o.PrivateIps = &v
}

// GetPublicIpIds returns the PublicIpIds field value if set, zero value otherwise.
func (o *FiltersPublicIp) GetPublicIpIds() []string {
	if o == nil || o.PublicIpIds == nil {
		var ret []string
		return ret
	}
	return *o.PublicIpIds
}

// GetPublicIpIdsOk returns a tuple with the PublicIpIds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersPublicIp) GetPublicIpIdsOk() (*[]string, bool) {
	if o == nil || o.PublicIpIds == nil {
		return nil, false
	}
	return o.PublicIpIds, true
}

// HasPublicIpIds returns a boolean if a field has been set.
func (o *FiltersPublicIp) HasPublicIpIds() bool {
	if o != nil && o.PublicIpIds != nil {
		return true
	}

	return false
}

// SetPublicIpIds gets a reference to the given []string and assigns it to the PublicIpIds field.
func (o *FiltersPublicIp) SetPublicIpIds(v []string) {
	o.PublicIpIds = &v
}

// GetPublicIps returns the PublicIps field value if set, zero value otherwise.
func (o *FiltersPublicIp) GetPublicIps() []string {
	if o == nil || o.PublicIps == nil {
		var ret []string
		return ret
	}
	return *o.PublicIps
}

// GetPublicIpsOk returns a tuple with the PublicIps field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersPublicIp) GetPublicIpsOk() (*[]string, bool) {
	if o == nil || o.PublicIps == nil {
		return nil, false
	}
	return o.PublicIps, true
}

// HasPublicIps returns a boolean if a field has been set.
func (o *FiltersPublicIp) HasPublicIps() bool {
	if o != nil && o.PublicIps != nil {
		return true
	}

	return false
}

// SetPublicIps gets a reference to the given []string and assigns it to the PublicIps field.
func (o *FiltersPublicIp) SetPublicIps(v []string) {
	o.PublicIps = &v
}

// GetTagKeys returns the TagKeys field value if set, zero value otherwise.
func (o *FiltersPublicIp) GetTagKeys() []string {
	if o == nil || o.TagKeys == nil {
		var ret []string
		return ret
	}
	return *o.TagKeys
}

// GetTagKeysOk returns a tuple with the TagKeys field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersPublicIp) GetTagKeysOk() (*[]string, bool) {
	if o == nil || o.TagKeys == nil {
		return nil, false
	}
	return o.TagKeys, true
}

// HasTagKeys returns a boolean if a field has been set.
func (o *FiltersPublicIp) HasTagKeys() bool {
	if o != nil && o.TagKeys != nil {
		return true
	}

	return false
}

// SetTagKeys gets a reference to the given []string and assigns it to the TagKeys field.
func (o *FiltersPublicIp) SetTagKeys(v []string) {
	o.TagKeys = &v
}

// GetTagValues returns the TagValues field value if set, zero value otherwise.
func (o *FiltersPublicIp) GetTagValues() []string {
	if o == nil || o.TagValues == nil {
		var ret []string
		return ret
	}
	return *o.TagValues
}

// GetTagValuesOk returns a tuple with the TagValues field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersPublicIp) GetTagValuesOk() (*[]string, bool) {
	if o == nil || o.TagValues == nil {
		return nil, false
	}
	return o.TagValues, true
}

// HasTagValues returns a boolean if a field has been set.
func (o *FiltersPublicIp) HasTagValues() bool {
	if o != nil && o.TagValues != nil {
		return true
	}

	return false
}

// SetTagValues gets a reference to the given []string and assigns it to the TagValues field.
func (o *FiltersPublicIp) SetTagValues(v []string) {
	o.TagValues = &v
}

// GetTags returns the Tags field value if set, zero value otherwise.
func (o *FiltersPublicIp) GetTags() []string {
	if o == nil || o.Tags == nil {
		var ret []string
		return ret
	}
	return *o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersPublicIp) GetTagsOk() (*[]string, bool) {
	if o == nil || o.Tags == nil {
		return nil, false
	}
	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *FiltersPublicIp) HasTags() bool {
	if o != nil && o.Tags != nil {
		return true
	}

	return false
}

// SetTags gets a reference to the given []string and assigns it to the Tags field.
func (o *FiltersPublicIp) SetTags(v []string) {
	o.Tags = &v
}

// GetVmIds returns the VmIds field value if set, zero value otherwise.
func (o *FiltersPublicIp) GetVmIds() []string {
	if o == nil || o.VmIds == nil {
		var ret []string
		return ret
	}
	return *o.VmIds
}

// GetVmIdsOk returns a tuple with the VmIds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FiltersPublicIp) GetVmIdsOk() (*[]string, bool) {
	if o == nil || o.VmIds == nil {
		return nil, false
	}
	return o.VmIds, true
}

// HasVmIds returns a boolean if a field has been set.
func (o *FiltersPublicIp) HasVmIds() bool {
	if o != nil && o.VmIds != nil {
		return true
	}

	return false
}

// SetVmIds gets a reference to the given []string and assigns it to the VmIds field.
func (o *FiltersPublicIp) SetVmIds(v []string) {
	o.VmIds = &v
}

func (o FiltersPublicIp) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.LinkPublicIpIds != nil {
		toSerialize["LinkPublicIpIds"] = o.LinkPublicIpIds
	}
	if o.NicAccountIds != nil {
		toSerialize["NicAccountIds"] = o.NicAccountIds
	}
	if o.NicIds != nil {
		toSerialize["NicIds"] = o.NicIds
	}
	if o.Placements != nil {
		toSerialize["Placements"] = o.Placements
	}
	if o.PrivateIps != nil {
		toSerialize["PrivateIps"] = o.PrivateIps
	}
	if o.PublicIpIds != nil {
		toSerialize["PublicIpIds"] = o.PublicIpIds
	}
	if o.PublicIps != nil {
		toSerialize["PublicIps"] = o.PublicIps
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
	if o.VmIds != nil {
		toSerialize["VmIds"] = o.VmIds
	}
	return json.Marshal(toSerialize)
}

type NullableFiltersPublicIp struct {
	value *FiltersPublicIp
	isSet bool
}

func (v NullableFiltersPublicIp) Get() *FiltersPublicIp {
	return v.value
}

func (v *NullableFiltersPublicIp) Set(val *FiltersPublicIp) {
	v.value = val
	v.isSet = true
}

func (v NullableFiltersPublicIp) IsSet() bool {
	return v.isSet
}

func (v *NullableFiltersPublicIp) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFiltersPublicIp(val *FiltersPublicIp) *NullableFiltersPublicIp {
	return &NullableFiltersPublicIp{value: val, isSet: true}
}

func (v NullableFiltersPublicIp) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFiltersPublicIp) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
