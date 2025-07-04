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

// OsuExportToCreate Information about the OOS export task to create.
type OsuExportToCreate struct {
	// The format of the export disk (`qcow2` \\| `raw`).
	DiskImageFormat string     `json:"DiskImageFormat"`
	OsuApiKey       *OsuApiKey `json:"OsuApiKey,omitempty"`
	// The name of the OOS bucket where you want to export the object.
	OsuBucket string `json:"OsuBucket"`
	// The URL of the manifest file.
	OsuManifestUrl *string `json:"OsuManifestUrl,omitempty"`
	// The prefix for the key of the OOS object.
	OsuPrefix *string `json:"OsuPrefix,omitempty"`
}

// NewOsuExportToCreate instantiates a new OsuExportToCreate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOsuExportToCreate(diskImageFormat string, osuBucket string) *OsuExportToCreate {
	this := OsuExportToCreate{}
	this.DiskImageFormat = diskImageFormat
	this.OsuBucket = osuBucket
	return &this
}

// NewOsuExportToCreateWithDefaults instantiates a new OsuExportToCreate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOsuExportToCreateWithDefaults() *OsuExportToCreate {
	this := OsuExportToCreate{}
	return &this
}

// GetDiskImageFormat returns the DiskImageFormat field value
func (o *OsuExportToCreate) GetDiskImageFormat() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DiskImageFormat
}

// GetDiskImageFormatOk returns a tuple with the DiskImageFormat field value
// and a boolean to check if the value has been set.
func (o *OsuExportToCreate) GetDiskImageFormatOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DiskImageFormat, true
}

// SetDiskImageFormat sets field value
func (o *OsuExportToCreate) SetDiskImageFormat(v string) {
	o.DiskImageFormat = v
}

// GetOsuApiKey returns the OsuApiKey field value if set, zero value otherwise.
func (o *OsuExportToCreate) GetOsuApiKey() OsuApiKey {
	if o == nil || o.OsuApiKey == nil {
		var ret OsuApiKey
		return ret
	}
	return *o.OsuApiKey
}

// GetOsuApiKeyOk returns a tuple with the OsuApiKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OsuExportToCreate) GetOsuApiKeyOk() (*OsuApiKey, bool) {
	if o == nil || o.OsuApiKey == nil {
		return nil, false
	}
	return o.OsuApiKey, true
}

// HasOsuApiKey returns a boolean if a field has been set.
func (o *OsuExportToCreate) HasOsuApiKey() bool {
	if o != nil && o.OsuApiKey != nil {
		return true
	}

	return false
}

// SetOsuApiKey gets a reference to the given OsuApiKey and assigns it to the OsuApiKey field.
func (o *OsuExportToCreate) SetOsuApiKey(v OsuApiKey) {
	o.OsuApiKey = &v
}

// GetOsuBucket returns the OsuBucket field value
func (o *OsuExportToCreate) GetOsuBucket() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.OsuBucket
}

// GetOsuBucketOk returns a tuple with the OsuBucket field value
// and a boolean to check if the value has been set.
func (o *OsuExportToCreate) GetOsuBucketOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OsuBucket, true
}

// SetOsuBucket sets field value
func (o *OsuExportToCreate) SetOsuBucket(v string) {
	o.OsuBucket = v
}

// GetOsuManifestUrl returns the OsuManifestUrl field value if set, zero value otherwise.
func (o *OsuExportToCreate) GetOsuManifestUrl() string {
	if o == nil || o.OsuManifestUrl == nil {
		var ret string
		return ret
	}
	return *o.OsuManifestUrl
}

// GetOsuManifestUrlOk returns a tuple with the OsuManifestUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OsuExportToCreate) GetOsuManifestUrlOk() (*string, bool) {
	if o == nil || o.OsuManifestUrl == nil {
		return nil, false
	}
	return o.OsuManifestUrl, true
}

// HasOsuManifestUrl returns a boolean if a field has been set.
func (o *OsuExportToCreate) HasOsuManifestUrl() bool {
	if o != nil && o.OsuManifestUrl != nil {
		return true
	}

	return false
}

// SetOsuManifestUrl gets a reference to the given string and assigns it to the OsuManifestUrl field.
func (o *OsuExportToCreate) SetOsuManifestUrl(v string) {
	o.OsuManifestUrl = &v
}

// GetOsuPrefix returns the OsuPrefix field value if set, zero value otherwise.
func (o *OsuExportToCreate) GetOsuPrefix() string {
	if o == nil || o.OsuPrefix == nil {
		var ret string
		return ret
	}
	return *o.OsuPrefix
}

// GetOsuPrefixOk returns a tuple with the OsuPrefix field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OsuExportToCreate) GetOsuPrefixOk() (*string, bool) {
	if o == nil || o.OsuPrefix == nil {
		return nil, false
	}
	return o.OsuPrefix, true
}

// HasOsuPrefix returns a boolean if a field has been set.
func (o *OsuExportToCreate) HasOsuPrefix() bool {
	if o != nil && o.OsuPrefix != nil {
		return true
	}

	return false
}

// SetOsuPrefix gets a reference to the given string and assigns it to the OsuPrefix field.
func (o *OsuExportToCreate) SetOsuPrefix(v string) {
	o.OsuPrefix = &v
}

func (o OsuExportToCreate) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["DiskImageFormat"] = o.DiskImageFormat
	}
	if o.OsuApiKey != nil {
		toSerialize["OsuApiKey"] = o.OsuApiKey
	}
	if true {
		toSerialize["OsuBucket"] = o.OsuBucket
	}
	if o.OsuManifestUrl != nil {
		toSerialize["OsuManifestUrl"] = o.OsuManifestUrl
	}
	if o.OsuPrefix != nil {
		toSerialize["OsuPrefix"] = o.OsuPrefix
	}
	return json.Marshal(toSerialize)
}

type NullableOsuExportToCreate struct {
	value *OsuExportToCreate
	isSet bool
}

func (v NullableOsuExportToCreate) Get() *OsuExportToCreate {
	return v.value
}

func (v *NullableOsuExportToCreate) Set(val *OsuExportToCreate) {
	v.value = val
	v.isSet = true
}

func (v NullableOsuExportToCreate) IsSet() bool {
	return v.isSet
}

func (v *NullableOsuExportToCreate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableOsuExportToCreate(val *OsuExportToCreate) *NullableOsuExportToCreate {
	return &NullableOsuExportToCreate{value: val, isSet: true}
}

func (v NullableOsuExportToCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableOsuExportToCreate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
