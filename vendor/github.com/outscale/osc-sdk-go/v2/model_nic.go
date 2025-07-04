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

// Nic Information about the NIC.
type Nic struct {
	// The account ID of the owner of the NIC.
	AccountId *string `json:"AccountId,omitempty"`
	// The description of the NIC.
	Description *string `json:"Description,omitempty"`
	// (Net only) If true, the source/destination check is enabled. If false, it is disabled.
	IsSourceDestChecked *bool         `json:"IsSourceDestChecked,omitempty"`
	LinkNic             *LinkNic      `json:"LinkNic,omitempty"`
	LinkPublicIp        *LinkPublicIp `json:"LinkPublicIp,omitempty"`
	// The Media Access Control (MAC) address of the NIC.
	MacAddress *string `json:"MacAddress,omitempty"`
	// The ID of the Net for the NIC.
	NetId *string `json:"NetId,omitempty"`
	// The ID of the NIC.
	NicId *string `json:"NicId,omitempty"`
	// The name of the private DNS.
	PrivateDnsName *string `json:"PrivateDnsName,omitempty"`
	// The private IPs of the NIC.
	PrivateIps *[]PrivateIp `json:"PrivateIps,omitempty"`
	// One or more IDs of security groups for the NIC.
	SecurityGroups *[]SecurityGroupLight `json:"SecurityGroups,omitempty"`
	// The state of the NIC (`available` \\| `attaching` \\| `in-use` \\| `detaching`).
	State *string `json:"State,omitempty"`
	// The ID of the Subnet.
	SubnetId *string `json:"SubnetId,omitempty"`
	// The Subregion in which the NIC is located.
	SubregionName *string `json:"SubregionName,omitempty"`
	// One or more tags associated with the NIC.
	Tags *[]ResourceTag `json:"Tags,omitempty"`
}

// NewNic instantiates a new Nic object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNic() *Nic {
	this := Nic{}
	return &this
}

// NewNicWithDefaults instantiates a new Nic object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNicWithDefaults() *Nic {
	this := Nic{}
	return &this
}

// GetAccountId returns the AccountId field value if set, zero value otherwise.
func (o *Nic) GetAccountId() string {
	if o == nil || o.AccountId == nil {
		var ret string
		return ret
	}
	return *o.AccountId
}

// GetAccountIdOk returns a tuple with the AccountId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nic) GetAccountIdOk() (*string, bool) {
	if o == nil || o.AccountId == nil {
		return nil, false
	}
	return o.AccountId, true
}

// HasAccountId returns a boolean if a field has been set.
func (o *Nic) HasAccountId() bool {
	if o != nil && o.AccountId != nil {
		return true
	}

	return false
}

// SetAccountId gets a reference to the given string and assigns it to the AccountId field.
func (o *Nic) SetAccountId(v string) {
	o.AccountId = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *Nic) GetDescription() string {
	if o == nil || o.Description == nil {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nic) GetDescriptionOk() (*string, bool) {
	if o == nil || o.Description == nil {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *Nic) HasDescription() bool {
	if o != nil && o.Description != nil {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *Nic) SetDescription(v string) {
	o.Description = &v
}

// GetIsSourceDestChecked returns the IsSourceDestChecked field value if set, zero value otherwise.
func (o *Nic) GetIsSourceDestChecked() bool {
	if o == nil || o.IsSourceDestChecked == nil {
		var ret bool
		return ret
	}
	return *o.IsSourceDestChecked
}

// GetIsSourceDestCheckedOk returns a tuple with the IsSourceDestChecked field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nic) GetIsSourceDestCheckedOk() (*bool, bool) {
	if o == nil || o.IsSourceDestChecked == nil {
		return nil, false
	}
	return o.IsSourceDestChecked, true
}

// HasIsSourceDestChecked returns a boolean if a field has been set.
func (o *Nic) HasIsSourceDestChecked() bool {
	if o != nil && o.IsSourceDestChecked != nil {
		return true
	}

	return false
}

// SetIsSourceDestChecked gets a reference to the given bool and assigns it to the IsSourceDestChecked field.
func (o *Nic) SetIsSourceDestChecked(v bool) {
	o.IsSourceDestChecked = &v
}

// GetLinkNic returns the LinkNic field value if set, zero value otherwise.
func (o *Nic) GetLinkNic() LinkNic {
	if o == nil || o.LinkNic == nil {
		var ret LinkNic
		return ret
	}
	return *o.LinkNic
}

// GetLinkNicOk returns a tuple with the LinkNic field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nic) GetLinkNicOk() (*LinkNic, bool) {
	if o == nil || o.LinkNic == nil {
		return nil, false
	}
	return o.LinkNic, true
}

// HasLinkNic returns a boolean if a field has been set.
func (o *Nic) HasLinkNic() bool {
	if o != nil && o.LinkNic != nil {
		return true
	}

	return false
}

// SetLinkNic gets a reference to the given LinkNic and assigns it to the LinkNic field.
func (o *Nic) SetLinkNic(v LinkNic) {
	o.LinkNic = &v
}

// GetLinkPublicIp returns the LinkPublicIp field value if set, zero value otherwise.
func (o *Nic) GetLinkPublicIp() LinkPublicIp {
	if o == nil || o.LinkPublicIp == nil {
		var ret LinkPublicIp
		return ret
	}
	return *o.LinkPublicIp
}

// GetLinkPublicIpOk returns a tuple with the LinkPublicIp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nic) GetLinkPublicIpOk() (*LinkPublicIp, bool) {
	if o == nil || o.LinkPublicIp == nil {
		return nil, false
	}
	return o.LinkPublicIp, true
}

// HasLinkPublicIp returns a boolean if a field has been set.
func (o *Nic) HasLinkPublicIp() bool {
	if o != nil && o.LinkPublicIp != nil {
		return true
	}

	return false
}

// SetLinkPublicIp gets a reference to the given LinkPublicIp and assigns it to the LinkPublicIp field.
func (o *Nic) SetLinkPublicIp(v LinkPublicIp) {
	o.LinkPublicIp = &v
}

// GetMacAddress returns the MacAddress field value if set, zero value otherwise.
func (o *Nic) GetMacAddress() string {
	if o == nil || o.MacAddress == nil {
		var ret string
		return ret
	}
	return *o.MacAddress
}

// GetMacAddressOk returns a tuple with the MacAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nic) GetMacAddressOk() (*string, bool) {
	if o == nil || o.MacAddress == nil {
		return nil, false
	}
	return o.MacAddress, true
}

// HasMacAddress returns a boolean if a field has been set.
func (o *Nic) HasMacAddress() bool {
	if o != nil && o.MacAddress != nil {
		return true
	}

	return false
}

// SetMacAddress gets a reference to the given string and assigns it to the MacAddress field.
func (o *Nic) SetMacAddress(v string) {
	o.MacAddress = &v
}

// GetNetId returns the NetId field value if set, zero value otherwise.
func (o *Nic) GetNetId() string {
	if o == nil || o.NetId == nil {
		var ret string
		return ret
	}
	return *o.NetId
}

// GetNetIdOk returns a tuple with the NetId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nic) GetNetIdOk() (*string, bool) {
	if o == nil || o.NetId == nil {
		return nil, false
	}
	return o.NetId, true
}

// HasNetId returns a boolean if a field has been set.
func (o *Nic) HasNetId() bool {
	if o != nil && o.NetId != nil {
		return true
	}

	return false
}

// SetNetId gets a reference to the given string and assigns it to the NetId field.
func (o *Nic) SetNetId(v string) {
	o.NetId = &v
}

// GetNicId returns the NicId field value if set, zero value otherwise.
func (o *Nic) GetNicId() string {
	if o == nil || o.NicId == nil {
		var ret string
		return ret
	}
	return *o.NicId
}

// GetNicIdOk returns a tuple with the NicId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nic) GetNicIdOk() (*string, bool) {
	if o == nil || o.NicId == nil {
		return nil, false
	}
	return o.NicId, true
}

// HasNicId returns a boolean if a field has been set.
func (o *Nic) HasNicId() bool {
	if o != nil && o.NicId != nil {
		return true
	}

	return false
}

// SetNicId gets a reference to the given string and assigns it to the NicId field.
func (o *Nic) SetNicId(v string) {
	o.NicId = &v
}

// GetPrivateDnsName returns the PrivateDnsName field value if set, zero value otherwise.
func (o *Nic) GetPrivateDnsName() string {
	if o == nil || o.PrivateDnsName == nil {
		var ret string
		return ret
	}
	return *o.PrivateDnsName
}

// GetPrivateDnsNameOk returns a tuple with the PrivateDnsName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nic) GetPrivateDnsNameOk() (*string, bool) {
	if o == nil || o.PrivateDnsName == nil {
		return nil, false
	}
	return o.PrivateDnsName, true
}

// HasPrivateDnsName returns a boolean if a field has been set.
func (o *Nic) HasPrivateDnsName() bool {
	if o != nil && o.PrivateDnsName != nil {
		return true
	}

	return false
}

// SetPrivateDnsName gets a reference to the given string and assigns it to the PrivateDnsName field.
func (o *Nic) SetPrivateDnsName(v string) {
	o.PrivateDnsName = &v
}

// GetPrivateIps returns the PrivateIps field value if set, zero value otherwise.
func (o *Nic) GetPrivateIps() []PrivateIp {
	if o == nil || o.PrivateIps == nil {
		var ret []PrivateIp
		return ret
	}
	return *o.PrivateIps
}

// GetPrivateIpsOk returns a tuple with the PrivateIps field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nic) GetPrivateIpsOk() (*[]PrivateIp, bool) {
	if o == nil || o.PrivateIps == nil {
		return nil, false
	}
	return o.PrivateIps, true
}

// HasPrivateIps returns a boolean if a field has been set.
func (o *Nic) HasPrivateIps() bool {
	if o != nil && o.PrivateIps != nil {
		return true
	}

	return false
}

// SetPrivateIps gets a reference to the given []PrivateIp and assigns it to the PrivateIps field.
func (o *Nic) SetPrivateIps(v []PrivateIp) {
	o.PrivateIps = &v
}

// GetSecurityGroups returns the SecurityGroups field value if set, zero value otherwise.
func (o *Nic) GetSecurityGroups() []SecurityGroupLight {
	if o == nil || o.SecurityGroups == nil {
		var ret []SecurityGroupLight
		return ret
	}
	return *o.SecurityGroups
}

// GetSecurityGroupsOk returns a tuple with the SecurityGroups field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nic) GetSecurityGroupsOk() (*[]SecurityGroupLight, bool) {
	if o == nil || o.SecurityGroups == nil {
		return nil, false
	}
	return o.SecurityGroups, true
}

// HasSecurityGroups returns a boolean if a field has been set.
func (o *Nic) HasSecurityGroups() bool {
	if o != nil && o.SecurityGroups != nil {
		return true
	}

	return false
}

// SetSecurityGroups gets a reference to the given []SecurityGroupLight and assigns it to the SecurityGroups field.
func (o *Nic) SetSecurityGroups(v []SecurityGroupLight) {
	o.SecurityGroups = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *Nic) GetState() string {
	if o == nil || o.State == nil {
		var ret string
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nic) GetStateOk() (*string, bool) {
	if o == nil || o.State == nil {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *Nic) HasState() bool {
	if o != nil && o.State != nil {
		return true
	}

	return false
}

// SetState gets a reference to the given string and assigns it to the State field.
func (o *Nic) SetState(v string) {
	o.State = &v
}

// GetSubnetId returns the SubnetId field value if set, zero value otherwise.
func (o *Nic) GetSubnetId() string {
	if o == nil || o.SubnetId == nil {
		var ret string
		return ret
	}
	return *o.SubnetId
}

// GetSubnetIdOk returns a tuple with the SubnetId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nic) GetSubnetIdOk() (*string, bool) {
	if o == nil || o.SubnetId == nil {
		return nil, false
	}
	return o.SubnetId, true
}

// HasSubnetId returns a boolean if a field has been set.
func (o *Nic) HasSubnetId() bool {
	if o != nil && o.SubnetId != nil {
		return true
	}

	return false
}

// SetSubnetId gets a reference to the given string and assigns it to the SubnetId field.
func (o *Nic) SetSubnetId(v string) {
	o.SubnetId = &v
}

// GetSubregionName returns the SubregionName field value if set, zero value otherwise.
func (o *Nic) GetSubregionName() string {
	if o == nil || o.SubregionName == nil {
		var ret string
		return ret
	}
	return *o.SubregionName
}

// GetSubregionNameOk returns a tuple with the SubregionName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nic) GetSubregionNameOk() (*string, bool) {
	if o == nil || o.SubregionName == nil {
		return nil, false
	}
	return o.SubregionName, true
}

// HasSubregionName returns a boolean if a field has been set.
func (o *Nic) HasSubregionName() bool {
	if o != nil && o.SubregionName != nil {
		return true
	}

	return false
}

// SetSubregionName gets a reference to the given string and assigns it to the SubregionName field.
func (o *Nic) SetSubregionName(v string) {
	o.SubregionName = &v
}

// GetTags returns the Tags field value if set, zero value otherwise.
func (o *Nic) GetTags() []ResourceTag {
	if o == nil || o.Tags == nil {
		var ret []ResourceTag
		return ret
	}
	return *o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nic) GetTagsOk() (*[]ResourceTag, bool) {
	if o == nil || o.Tags == nil {
		return nil, false
	}
	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *Nic) HasTags() bool {
	if o != nil && o.Tags != nil {
		return true
	}

	return false
}

// SetTags gets a reference to the given []ResourceTag and assigns it to the Tags field.
func (o *Nic) SetTags(v []ResourceTag) {
	o.Tags = &v
}

func (o Nic) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.AccountId != nil {
		toSerialize["AccountId"] = o.AccountId
	}
	if o.Description != nil {
		toSerialize["Description"] = o.Description
	}
	if o.IsSourceDestChecked != nil {
		toSerialize["IsSourceDestChecked"] = o.IsSourceDestChecked
	}
	if o.LinkNic != nil {
		toSerialize["LinkNic"] = o.LinkNic
	}
	if o.LinkPublicIp != nil {
		toSerialize["LinkPublicIp"] = o.LinkPublicIp
	}
	if o.MacAddress != nil {
		toSerialize["MacAddress"] = o.MacAddress
	}
	if o.NetId != nil {
		toSerialize["NetId"] = o.NetId
	}
	if o.NicId != nil {
		toSerialize["NicId"] = o.NicId
	}
	if o.PrivateDnsName != nil {
		toSerialize["PrivateDnsName"] = o.PrivateDnsName
	}
	if o.PrivateIps != nil {
		toSerialize["PrivateIps"] = o.PrivateIps
	}
	if o.SecurityGroups != nil {
		toSerialize["SecurityGroups"] = o.SecurityGroups
	}
	if o.State != nil {
		toSerialize["State"] = o.State
	}
	if o.SubnetId != nil {
		toSerialize["SubnetId"] = o.SubnetId
	}
	if o.SubregionName != nil {
		toSerialize["SubregionName"] = o.SubregionName
	}
	if o.Tags != nil {
		toSerialize["Tags"] = o.Tags
	}
	return json.Marshal(toSerialize)
}

type NullableNic struct {
	value *Nic
	isSet bool
}

func (v NullableNic) Get() *Nic {
	return v.value
}

func (v *NullableNic) Set(val *Nic) {
	v.value = val
	v.isSet = true
}

func (v NullableNic) IsSet() bool {
	return v.isSet
}

func (v *NullableNic) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNic(val *Nic) *NullableNic {
	return &NullableNic{value: val, isSet: true}
}

func (v NullableNic) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNic) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
