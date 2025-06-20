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

// CreateLoadBalancerRequest struct for CreateLoadBalancerRequest
type CreateLoadBalancerRequest struct {
	// If true, checks whether you have the required permissions to perform the action.
	DryRun *bool `json:"DryRun,omitempty"`
	// One or more listeners to create.
	Listeners []ListenerForCreation `json:"Listeners"`
	// The unique name of the load balancer, with a maximum length of 32 alphanumeric characters and dashes (`-`). This name must not start or end with a dash.
	LoadBalancerName string `json:"LoadBalancerName"`
	// The type of load balancer: `internet-facing` or `internal`. Use this parameter only for load balancers in a Net.
	LoadBalancerType *string `json:"LoadBalancerType,omitempty"`
	// (internet-facing only) The public IP you want to associate with the load balancer. If not specified, a public IP owned by 3DS OUTSCALE is associated.
	PublicIp *string `json:"PublicIp,omitempty"`
	// (Net only) One or more IDs of security groups you want to assign to the load balancer. If not specified, the default security group of the Net is assigned to the load balancer.
	SecurityGroups *[]string `json:"SecurityGroups,omitempty"`
	// (Net only) The ID of the Subnet in which you want to create the load balancer. Regardless of this Subnet, the load balancer can distribute traffic to all Subnets. This parameter is required in a Net.
	Subnets *[]string `json:"Subnets,omitempty"`
	// (public Cloud only) The Subregion in which you want to create the load balancer. Regardless of this Subregion, the load balancer can distribute traffic to all Subregions. This parameter is required in the public Cloud.
	SubregionNames *[]string `json:"SubregionNames,omitempty"`
	// One or more tags assigned to the load balancer.
	Tags *[]ResourceTag `json:"Tags,omitempty"`
}

// NewCreateLoadBalancerRequest instantiates a new CreateLoadBalancerRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateLoadBalancerRequest(listeners []ListenerForCreation, loadBalancerName string) *CreateLoadBalancerRequest {
	this := CreateLoadBalancerRequest{}
	this.Listeners = listeners
	this.LoadBalancerName = loadBalancerName
	return &this
}

// NewCreateLoadBalancerRequestWithDefaults instantiates a new CreateLoadBalancerRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateLoadBalancerRequestWithDefaults() *CreateLoadBalancerRequest {
	this := CreateLoadBalancerRequest{}
	return &this
}

// GetDryRun returns the DryRun field value if set, zero value otherwise.
func (o *CreateLoadBalancerRequest) GetDryRun() bool {
	if o == nil || o.DryRun == nil {
		var ret bool
		return ret
	}
	return *o.DryRun
}

// GetDryRunOk returns a tuple with the DryRun field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateLoadBalancerRequest) GetDryRunOk() (*bool, bool) {
	if o == nil || o.DryRun == nil {
		return nil, false
	}
	return o.DryRun, true
}

// HasDryRun returns a boolean if a field has been set.
func (o *CreateLoadBalancerRequest) HasDryRun() bool {
	if o != nil && o.DryRun != nil {
		return true
	}

	return false
}

// SetDryRun gets a reference to the given bool and assigns it to the DryRun field.
func (o *CreateLoadBalancerRequest) SetDryRun(v bool) {
	o.DryRun = &v
}

// GetListeners returns the Listeners field value
func (o *CreateLoadBalancerRequest) GetListeners() []ListenerForCreation {
	if o == nil {
		var ret []ListenerForCreation
		return ret
	}

	return o.Listeners
}

// GetListenersOk returns a tuple with the Listeners field value
// and a boolean to check if the value has been set.
func (o *CreateLoadBalancerRequest) GetListenersOk() (*[]ListenerForCreation, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Listeners, true
}

// SetListeners sets field value
func (o *CreateLoadBalancerRequest) SetListeners(v []ListenerForCreation) {
	o.Listeners = v
}

// GetLoadBalancerName returns the LoadBalancerName field value
func (o *CreateLoadBalancerRequest) GetLoadBalancerName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.LoadBalancerName
}

// GetLoadBalancerNameOk returns a tuple with the LoadBalancerName field value
// and a boolean to check if the value has been set.
func (o *CreateLoadBalancerRequest) GetLoadBalancerNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.LoadBalancerName, true
}

// SetLoadBalancerName sets field value
func (o *CreateLoadBalancerRequest) SetLoadBalancerName(v string) {
	o.LoadBalancerName = v
}

// GetLoadBalancerType returns the LoadBalancerType field value if set, zero value otherwise.
func (o *CreateLoadBalancerRequest) GetLoadBalancerType() string {
	if o == nil || o.LoadBalancerType == nil {
		var ret string
		return ret
	}
	return *o.LoadBalancerType
}

// GetLoadBalancerTypeOk returns a tuple with the LoadBalancerType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateLoadBalancerRequest) GetLoadBalancerTypeOk() (*string, bool) {
	if o == nil || o.LoadBalancerType == nil {
		return nil, false
	}
	return o.LoadBalancerType, true
}

// HasLoadBalancerType returns a boolean if a field has been set.
func (o *CreateLoadBalancerRequest) HasLoadBalancerType() bool {
	if o != nil && o.LoadBalancerType != nil {
		return true
	}

	return false
}

// SetLoadBalancerType gets a reference to the given string and assigns it to the LoadBalancerType field.
func (o *CreateLoadBalancerRequest) SetLoadBalancerType(v string) {
	o.LoadBalancerType = &v
}

// GetPublicIp returns the PublicIp field value if set, zero value otherwise.
func (o *CreateLoadBalancerRequest) GetPublicIp() string {
	if o == nil || o.PublicIp == nil {
		var ret string
		return ret
	}
	return *o.PublicIp
}

// GetPublicIpOk returns a tuple with the PublicIp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateLoadBalancerRequest) GetPublicIpOk() (*string, bool) {
	if o == nil || o.PublicIp == nil {
		return nil, false
	}
	return o.PublicIp, true
}

// HasPublicIp returns a boolean if a field has been set.
func (o *CreateLoadBalancerRequest) HasPublicIp() bool {
	if o != nil && o.PublicIp != nil {
		return true
	}

	return false
}

// SetPublicIp gets a reference to the given string and assigns it to the PublicIp field.
func (o *CreateLoadBalancerRequest) SetPublicIp(v string) {
	o.PublicIp = &v
}

// GetSecurityGroups returns the SecurityGroups field value if set, zero value otherwise.
func (o *CreateLoadBalancerRequest) GetSecurityGroups() []string {
	if o == nil || o.SecurityGroups == nil {
		var ret []string
		return ret
	}
	return *o.SecurityGroups
}

// GetSecurityGroupsOk returns a tuple with the SecurityGroups field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateLoadBalancerRequest) GetSecurityGroupsOk() (*[]string, bool) {
	if o == nil || o.SecurityGroups == nil {
		return nil, false
	}
	return o.SecurityGroups, true
}

// HasSecurityGroups returns a boolean if a field has been set.
func (o *CreateLoadBalancerRequest) HasSecurityGroups() bool {
	if o != nil && o.SecurityGroups != nil {
		return true
	}

	return false
}

// SetSecurityGroups gets a reference to the given []string and assigns it to the SecurityGroups field.
func (o *CreateLoadBalancerRequest) SetSecurityGroups(v []string) {
	o.SecurityGroups = &v
}

// GetSubnets returns the Subnets field value if set, zero value otherwise.
func (o *CreateLoadBalancerRequest) GetSubnets() []string {
	if o == nil || o.Subnets == nil {
		var ret []string
		return ret
	}
	return *o.Subnets
}

// GetSubnetsOk returns a tuple with the Subnets field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateLoadBalancerRequest) GetSubnetsOk() (*[]string, bool) {
	if o == nil || o.Subnets == nil {
		return nil, false
	}
	return o.Subnets, true
}

// HasSubnets returns a boolean if a field has been set.
func (o *CreateLoadBalancerRequest) HasSubnets() bool {
	if o != nil && o.Subnets != nil {
		return true
	}

	return false
}

// SetSubnets gets a reference to the given []string and assigns it to the Subnets field.
func (o *CreateLoadBalancerRequest) SetSubnets(v []string) {
	o.Subnets = &v
}

// GetSubregionNames returns the SubregionNames field value if set, zero value otherwise.
func (o *CreateLoadBalancerRequest) GetSubregionNames() []string {
	if o == nil || o.SubregionNames == nil {
		var ret []string
		return ret
	}
	return *o.SubregionNames
}

// GetSubregionNamesOk returns a tuple with the SubregionNames field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateLoadBalancerRequest) GetSubregionNamesOk() (*[]string, bool) {
	if o == nil || o.SubregionNames == nil {
		return nil, false
	}
	return o.SubregionNames, true
}

// HasSubregionNames returns a boolean if a field has been set.
func (o *CreateLoadBalancerRequest) HasSubregionNames() bool {
	if o != nil && o.SubregionNames != nil {
		return true
	}

	return false
}

// SetSubregionNames gets a reference to the given []string and assigns it to the SubregionNames field.
func (o *CreateLoadBalancerRequest) SetSubregionNames(v []string) {
	o.SubregionNames = &v
}

// GetTags returns the Tags field value if set, zero value otherwise.
func (o *CreateLoadBalancerRequest) GetTags() []ResourceTag {
	if o == nil || o.Tags == nil {
		var ret []ResourceTag
		return ret
	}
	return *o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateLoadBalancerRequest) GetTagsOk() (*[]ResourceTag, bool) {
	if o == nil || o.Tags == nil {
		return nil, false
	}
	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *CreateLoadBalancerRequest) HasTags() bool {
	if o != nil && o.Tags != nil {
		return true
	}

	return false
}

// SetTags gets a reference to the given []ResourceTag and assigns it to the Tags field.
func (o *CreateLoadBalancerRequest) SetTags(v []ResourceTag) {
	o.Tags = &v
}

func (o CreateLoadBalancerRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.DryRun != nil {
		toSerialize["DryRun"] = o.DryRun
	}
	if true {
		toSerialize["Listeners"] = o.Listeners
	}
	if true {
		toSerialize["LoadBalancerName"] = o.LoadBalancerName
	}
	if o.LoadBalancerType != nil {
		toSerialize["LoadBalancerType"] = o.LoadBalancerType
	}
	if o.PublicIp != nil {
		toSerialize["PublicIp"] = o.PublicIp
	}
	if o.SecurityGroups != nil {
		toSerialize["SecurityGroups"] = o.SecurityGroups
	}
	if o.Subnets != nil {
		toSerialize["Subnets"] = o.Subnets
	}
	if o.SubregionNames != nil {
		toSerialize["SubregionNames"] = o.SubregionNames
	}
	if o.Tags != nil {
		toSerialize["Tags"] = o.Tags
	}
	return json.Marshal(toSerialize)
}

type NullableCreateLoadBalancerRequest struct {
	value *CreateLoadBalancerRequest
	isSet bool
}

func (v NullableCreateLoadBalancerRequest) Get() *CreateLoadBalancerRequest {
	return v.value
}

func (v *NullableCreateLoadBalancerRequest) Set(val *CreateLoadBalancerRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateLoadBalancerRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateLoadBalancerRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateLoadBalancerRequest(val *CreateLoadBalancerRequest) *NullableCreateLoadBalancerRequest {
	return &NullableCreateLoadBalancerRequest{value: val, isSet: true}
}

func (v NullableCreateLoadBalancerRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateLoadBalancerRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
