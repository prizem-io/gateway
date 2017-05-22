package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Auditable struct {
	CreatedBy    string    `json:"createdBy" yaml:"createdBy" msgpack:"createdBy" valid:"required"`
	CreatedDate  time.Time `json:"createdDate" yaml:"createdDate" msgpack:"createdDate" valid:"required"`
	ModifiedBy   string    `json:"modifiedBy" yaml:"modifiedBy" msgpack:"modifiedBy" valid:"required"`
	ModifiedDate time.Time `json:"modifiedDate" yaml:"modifiedDate" msgpack:"modifiedDate" valid:"required"`
}

type Authentication struct {
	Username string `json:"username" yaml:"username" msgpack:"username" valid:"required"`
	Password string `json:"password" yaml:"password" msgpack:"password" valid:"required"`
}

type ClaimEntry struct {
	Type  string       `json:"type" yaml:"type" msgpack:"type" valid:"required"`
	Value *interface{} `json:"value" yaml:"value" msgpack:"value"`
}

type Client struct {
	Entity       `msgpack:",inline" mapstructure:",squash"`
	ClientUpdate `msgpack:",inline" mapstructure:",squash"`
	Auditable    `msgpack:",inline" mapstructure:",squash"`
}

type ClientUpdate struct {
	ConsumerID          string                            `json:"consumerId" yaml:"consumerId" msgpack:"consumerId" valid:"required"`
	Label               string                            `json:"label" yaml:"label" msgpack:"label" valid:"required"`
	Enabled             bool                              `json:"enabled" yaml:"enabled" msgpack:"enabled" valid:"required"`
	Status              *string                           `json:"status" yaml:"status" msgpack:"status"`
	APIKey              string                            `json:"apiKey" yaml:"apiKey" msgpack:"apiKey" valid:"required"`
	SharedSecret        *string                           `json:"sharedSecret" yaml:"sharedSecret" msgpack:"sharedSecret"`
	Authenticators      map[string]map[string]interface{} `json:"authenticators" yaml:"authenticators" msgpack:"authenticators" valid:"required"`
	Filters             []PluginConfig                    `json:"filters" yaml:"filters" msgpack:"filters" valid:"required"`
	ClientPermissionIDs []string                          `json:"clientPermissionIds" yaml:"clientPermissionIds" msgpack:"clientPermissionIds" valid:"required"`
	PermissionSets      map[string]PermissionSet          `json:"permissionSets" yaml:"permissionSets" msgpack:"permissionSets" valid:"required"`
	Extended            map[string]interface{}            `json:"extended" yaml:"extended" msgpack:"extended" valid:"required"`
	ExternalID          *string                           `json:"externalId" yaml:"externalId" msgpack:"externalId"`
}

type ConfigurableDescriptor struct {
	Name       string               `json:"name" yaml:"name" msgpack:"name" valid:"required"`
	Properties []PropertyDescriptor `json:"properties" yaml:"properties" msgpack:"properties"`
}

type Consumer struct {
	Entity         `msgpack:",inline" mapstructure:",squash"`
	ConsumerUpdate `msgpack:",inline" mapstructure:",squash"`
	Auditable      `msgpack:",inline" mapstructure:",squash"`
}

type ConsumerUpdate struct {
	Name                string                 `json:"name" yaml:"name" msgpack:"name" valid:"required"`
	Description         *string                `json:"description" yaml:"description" msgpack:"description"`
	ApplicationURL      string                 `json:"applicationUrl" yaml:"applicationUrl" msgpack:"applicationUrl" valid:"required"`
	ApplicationImageURL *string                `json:"applicationImageUrl" yaml:"applicationImageUrl" msgpack:"applicationImageUrl"`
	CompanyName         string                 `json:"companyName" yaml:"companyName" msgpack:"companyName" valid:"required"`
	CompanyURL          string                 `json:"companyUrl" yaml:"companyUrl" msgpack:"companyUrl" valid:"required"`
	CompanyImageURL     *string                `json:"companyImageUrl" yaml:"companyImageUrl" msgpack:"companyImageUrl"`
	PermissionIDs       []string               `json:"permissionIds" yaml:"permissionIds" msgpack:"permissionIds" valid:"required"`
	Filters             []PluginConfig         `json:"filters" yaml:"filters" msgpack:"filters" valid:"required"`
	Tags                []string               `json:"tags" yaml:"tags" msgpack:"tags" valid:"required"`
	PlanID              *string                `json:"planId" yaml:"planId" msgpack:"planId"`
	Extended            map[string]interface{} `json:"extended" yaml:"extended" msgpack:"extended" valid:"required"`
	ExternalID          *string                `json:"externalId" yaml:"externalId" msgpack:"externalId"`
}

type Credential struct {
	Entity           `msgpack:",inline" mapstructure:",squash"`
	Auditable        `msgpack:",inline" mapstructure:",squash"`
	CredentialUpdate `msgpack:",inline" mapstructure:",squash"`
}

type CredentialUpdate struct {
	Name        string `json:"name" yaml:"name" msgpack:"name" valid:"required"`
	Enabled     bool   `json:"enabled" yaml:"enabled" msgpack:"enabled" valid:"required"`
	Type        string `json:"type" yaml:"type" msgpack:"type" valid:"required"`
	SubjectID   string `json:"subjectId" yaml:"subjectId" msgpack:"subjectId" valid:"required"`
	SubjectType string `json:"subjectType" yaml:"subjectType" msgpack:"subjectType" valid:"required"`
}

type Developer struct {
	Entity          `msgpack:",inline" mapstructure:",squash"`
	DeveloperUpdate `msgpack:",inline" mapstructure:",squash"`
	Auditable       `msgpack:",inline" mapstructure:",squash"`
}

type DeveloperUpdate struct {
	Username       string                 `json:"username" yaml:"username" msgpack:"username" valid:"required"`
	FirstName      string                 `json:"firstName" yaml:"firstName" msgpack:"firstName" valid:"required"`
	LastName       string                 `json:"lastName" yaml:"lastName" msgpack:"lastName" valid:"required"`
	Roles          []string               `json:"roles" yaml:"roles" msgpack:"roles" valid:"required"`
	Company        *string                `json:"company" yaml:"company" msgpack:"company"`
	Title          *string                `json:"title" yaml:"title" msgpack:"title"`
	Email          *string                `json:"email" yaml:"email" msgpack:"email"`
	Phone          *string                `json:"phone" yaml:"phone" msgpack:"phone"`
	Mobile         *string                `json:"mobile" yaml:"mobile" msgpack:"mobile"`
	Address1       *string                `json:"address1" yaml:"address1" msgpack:"address1"`
	Address2       *string                `json:"address2" yaml:"address2" msgpack:"address2"`
	Locality       *string                `json:"locality" yaml:"locality" msgpack:"locality"`
	Region         *string                `json:"region" yaml:"region" msgpack:"region"`
	PostalCode     *string                `json:"postalCode" yaml:"postalCode" msgpack:"postalCode"`
	CountryCode    *string                `json:"countryCode" yaml:"countryCode" msgpack:"countryCode"`
	RegistrationIP *string                `json:"registrationIp" yaml:"registrationIp" msgpack:"registrationIp"`
	Extended       map[string]interface{} `json:"extended" yaml:"extended" msgpack:"extended" valid:"required"`
	ExternalID     *string                `json:"externalId" yaml:"externalId" msgpack:"externalId"`
}

type Entity struct {
	ID string `json:"id" yaml:"id" msgpack:"id" valid:"required"`
}

type Environment struct {
	Entity            `msgpack:",inline" mapstructure:",squash"`
	EnvironmentUpdate `msgpack:",inline" mapstructure:",squash"`
	Auditable         `msgpack:",inline" mapstructure:",squash"`
}

type EnvironmentUpdate struct {
	Name              string  `json:"name" yaml:"name" msgpack:"name" valid:"required"`
	Description       *string `json:"description" yaml:"description" msgpack:"description"`
	SystemDatabase    string  `json:"systemDatabase" yaml:"systemDatabase" msgpack:"systemDatabase" valid:"required"`
	AnalyticsDatabase string  `json:"analyticsDatabase" yaml:"analyticsDatabase" msgpack:"analyticsDatabase" valid:"required"`
}

type InternalActor struct {
	Permissions  []string          `json:"permissions" yaml:"permissions" msgpack:"permissions" valid:"required"`
	AccessLevels map[string]string `json:"accessLevels" yaml:"accessLevels" msgpack:"accessLevels" valid:"required"`
}

type LogEntry struct {
	Entity         `msgpack:",inline" mapstructure:",squash"`
	LogEntryUpdate `msgpack:",inline" mapstructure:",squash"`
	Auditable      `msgpack:",inline" mapstructure:",squash"`
}

type LogEntryUpdate struct {
	Level     string    `json:"level" yaml:"level" msgpack:"level" valid:"required"`
	Timestamp time.Time `json:"timestamp" yaml:"timestamp" msgpack:"timestamp" valid:"required"`
	Message   string    `json:"message" yaml:"message" msgpack:"message" valid:"required"`
}

type Message struct {
	Entity        `msgpack:",inline" mapstructure:",squash"`
	MessageUpdate `msgpack:",inline" mapstructure:",squash"`
	Auditable     `msgpack:",inline" mapstructure:",squash"`
}

type MessageContent struct {
	Format  string `json:"format" yaml:"format" msgpack:"format" valid:"required"`
	Content string `json:"content" yaml:"content" msgpack:"content" valid:"required"`
}

type MessageUpdate struct {
	ParentID     *string                   `json:"parentId" yaml:"parentId" msgpack:"parentId"`
	Key          string                    `json:"key" yaml:"key" msgpack:"key" valid:"required"`
	Locales      map[string]MessageContent `json:"locales" yaml:"locales" msgpack:"locales" valid:"required"`
	DisplayOrder int32                     `json:"displayOrder" yaml:"displayOrder" msgpack:"displayOrder" valid:"required"`
}

type Operation struct {
	Name          string         `json:"name" yaml:"name" msgpack:"name" valid:"required"`
	Method        Method         `json:"method" yaml:"method" msgpack:"method" valid:"required"`
	URIPattern    string         `json:"uriPattern" yaml:"uriPattern" msgpack:"uriPattern" valid:"required"`
	PermissionIDs []string       `json:"permissionIds" yaml:"permissionIds" msgpack:"permissionIds" valid:"required"`
	Claims        []ClaimEntry   `json:"claims" yaml:"claims" msgpack:"claims" valid:"required"`
	Filters       []PluginConfig `json:"filters" yaml:"filters" msgpack:"filters" valid:"required"`
	Backend       *PluginConfig  `json:"backend" yaml:"backend" msgpack:"backend"`
	BackendConfig interface{}
}

type Permission struct {
	Entity           `msgpack:",inline" mapstructure:",squash"`
	PermissionUpdate `msgpack:",inline" mapstructure:",squash"`
	Auditable        `msgpack:",inline" mapstructure:",squash"`
}

type PermissionSet struct {
	Enabled       bool     `json:"enabled" yaml:"enabled" msgpack:"enabled" valid:"required"`
	Global        bool     `json:"global" yaml:"global" msgpack:"global" valid:"required"`
	Expiration    *int64   `json:"expiration" yaml:"expiration" msgpack:"expiration"`
	Lifespan      Lifespan `json:"lifespan" yaml:"lifespan" msgpack:"lifespan" valid:"required"`
	Refreshable   bool     `json:"refreshable" yaml:"refreshable" msgpack:"refreshable" valid:"required"`
	PermissionIDs []string `json:"permissionIds" yaml:"permissionIds" msgpack:"permissionIds" valid:"required"`
	AutoAuthorize bool     `json:"autoAuthorize" yaml:"autoAuthorize" msgpack:"autoAuthorize" valid:"required"`
	AllowWebView  bool     `json:"allowWebView" yaml:"allowWebView" msgpack:"allowWebView" valid:"required"`
	AllowPopup    bool     `json:"allowPopup" yaml:"allowPopup" msgpack:"allowPopup" valid:"required"`
}

type PermissionUpdate struct {
	Name        string      `json:"name" yaml:"name" msgpack:"name" valid:"required"`
	Description *string     `json:"description" yaml:"description" msgpack:"description"`
	Type        Type        `json:"type" yaml:"type" msgpack:"type" valid:"required"`
	Scope       Scope       `json:"scope" yaml:"scope" msgpack:"scope" valid:"required"`
	MessageID   string      `json:"messageId" yaml:"messageId" msgpack:"messageId" valid:"required"`
	ClaimPath   []string    `json:"claimPath" yaml:"claimPath" msgpack:"claimPath" valid:"required"`
	ClaimValue  interface{} `json:"claimValue" yaml:"claimValue" msgpack:"claimValue" valid:"required"`
	Flags       []string    `json:"flags" yaml:"flags" msgpack:"flags" valid:"required"`
}

type Plan struct {
	Entity     `msgpack:",inline" mapstructure:",squash"`
	PlanUpdate `msgpack:",inline" mapstructure:",squash"`
	Auditable  `msgpack:",inline" mapstructure:",squash"`
}

type PlanUpdate struct {
	Name          string         `json:"name" yaml:"name" msgpack:"name" valid:"required"`
	PriceAmount   float32        `json:"priceAmount" yaml:"priceAmount" msgpack:"priceAmount" valid:"required"`
	PriceCurrency string         `json:"priceCurrency" yaml:"priceCurrency" msgpack:"priceCurrency" valid:"required"`
	Filters       []PluginConfig `json:"filters" yaml:"filters" msgpack:"filters" valid:"required"`
	Quotas        []Quota        `json:"quotas" yaml:"quotas" msgpack:"quotas" valid:"required"`
}

type Plugin struct {
	Entity       `msgpack:",inline" mapstructure:",squash"`
	PluginConfig `msgpack:",inline" mapstructure:",squash"`
	Auditable    `msgpack:",inline" mapstructure:",squash"`
}

type PluginConfig struct {
	Name       string                 `json:"name" yaml:"name" msgpack:"name" valid:"required"`
	Properties map[string]interface{} `json:"properties" yaml:"properties" msgpack:"properties" valid:"required"`
	Config     interface{}
}

type PrincipalClaims struct {
	Entity                `msgpack:",inline" mapstructure:",squash"`
	PrincipalClaimsUpdate `msgpack:",inline" mapstructure:",squash"`
	Auditable             `msgpack:",inline" mapstructure:",squash"`
}

type PrincipalClaimsUpdate struct {
	ProfileID string                   `json:"profileId" yaml:"profileId" msgpack:"profileId" valid:"required"`
	Name      string                   `json:"name" yaml:"name" msgpack:"name" valid:"required"`
	Inherits  []string                 `json:"inherits" yaml:"inherits" msgpack:"inherits" valid:"required"`
	Claims    map[string][]interface{} `json:"claims" yaml:"claims" msgpack:"claims" valid:"required"`
}

type PrincipalProfile struct {
	Entity                 `msgpack:",inline" mapstructure:",squash"`
	PrincipalProfileUpdate `msgpack:",inline" mapstructure:",squash"`
	Auditable              `msgpack:",inline" mapstructure:",squash"`
}

type PrincipalProfileUpdate struct {
	Name string `json:"name" yaml:"name" msgpack:"name" valid:"required"`
}

type PropertyDescriptor struct {
	PropertyName string       `json:"propertyName" yaml:"propertyName" msgpack:"propertyName" valid:"required"`
	DisplayName  string       `json:"displayName" yaml:"displayName" msgpack:"displayName" valid:"required"`
	PropertyType PropertyType `json:"propertyType" yaml:"propertyType" msgpack:"propertyType" valid:"required"`
	Required     bool         `json:"required" yaml:"required" msgpack:"required" valid:"required"`
	Multi        bool         `json:"multi" yaml:"multi" msgpack:"multi" valid:"required"`
}

type Provider struct {
	Entity         `msgpack:",inline" mapstructure:",squash"`
	ProviderUpdate `msgpack:",inline" mapstructure:",squash"`
	Auditable      `msgpack:",inline" mapstructure:",squash"`
}

type ProviderUpdate struct {
	InternalActor      `msgpack:",inline" mapstructure:",squash"`
	Name               string                 `json:"name" yaml:"name" msgpack:"name"`
	Enabled            bool                   `json:"enabled" yaml:"enabled" msgpack:"enabled"`
	BehindReverseProxy bool                   `json:"behindReverseProxy" yaml:"behindReverseProxy" msgpack:"behindReverseProxy"`
	Extended           map[string]interface{} `json:"extended" yaml:"extended" msgpack:"extended"`
}

type Quota struct {
	RequestCount int32    `json:"requestCount" yaml:"requestCount" msgpack:"requestCount" valid:"required"`
	TimeUnit     TimeUnit `json:"timeUnit" yaml:"timeUnit" msgpack:"timeUnit" valid:"required"`
}

type Result struct {
	Result string `json:"result" yaml:"result" msgpack:"result" valid:"required"`
}

type Role struct {
	Entity     `msgpack:",inline" mapstructure:",squash"`
	RoleUpdate `msgpack:",inline" mapstructure:",squash"`
	Auditable  `msgpack:",inline" mapstructure:",squash"`
}

type RoleUpdate struct {
	InternalActor `msgpack:",inline" mapstructure:",squash"`
	Name          string `json:"name" yaml:"name" msgpack:"name"`
	DisplayName   string `json:"displayName" yaml:"displayName" msgpack:"displayName"`
	Description   string `json:"description" yaml:"description" msgpack:"description"`
}

type Service struct {
	Entity        `msgpack:",inline" mapstructure:",squash"`
	ServiceUpdate `msgpack:",inline" mapstructure:",squash"`
	Auditable     `msgpack:",inline" mapstructure:",squash"`
}

type ServiceUpdate struct {
	Name                 string                 `json:"name" yaml:"name" msgpack:"name" valid:"required"`
	Type                 *string                `json:"type" yaml:"type" msgpack:"type"`
	Description          *string                `json:"description" yaml:"description" msgpack:"description"`
	Hostnames            []string               `json:"hostnames" yaml:"hostnames" msgpack:"hostnames" valid:"required"`
	URIPrefix            *string                `json:"uriPrefix" yaml:"uriPrefix" msgpack:"uriPrefix"`
	VersionLocation      *string                `json:"versionLocation" yaml:"versionLocation" msgpack:"versionLocation"`
	DefaultVersion       string                 `json:"defaultVersion" yaml:"defaultVersion" msgpack:"defaultVersion" valid:"required"`
	Scheme               *string                `json:"scheme" yaml:"scheme" msgpack:"scheme"`
	ContextRoot          *string                `json:"contextRoot" yaml:"contextRoot" msgpack:"contextRoot"`
	RequestWeights       map[string]int32       `json:"requestWeights" yaml:"requestWeights" msgpack:"requestWeights" valid:"required"`
	AuthenticationType   AuthenticationType     `json:"authenticationType" yaml:"authenticationType" msgpack:"authenticationType" valid:"required"`
	GlobalClaims         []ClaimEntry           `json:"globalClaims" yaml:"globalClaims" msgpack:"globalClaims" valid:"required"`
	AccessControlEnabled bool                   `json:"accessControlEnabled" yaml:"accessControlEnabled" msgpack:"accessControlEnabled" valid:"required"`
	Operations           []Operation            `json:"operations" yaml:"operations" msgpack:"operations" valid:"required"`
	Filters              []PluginConfig         `json:"filters" yaml:"filters" msgpack:"filters" valid:"required"`
	Tags                 []string               `json:"tags" yaml:"tags" msgpack:"tags" valid:"required"`
	Backend              *PluginConfig          `json:"backend" yaml:"backend" msgpack:"backend" valid:"required"`
	Extended             map[string]interface{} `json:"extended" yaml:"extended" msgpack:"extended" valid:"required"`
}

type SetPassword struct {
	Password string `json:"password" yaml:"password" msgpack:"password" valid:"required"`
}

type Summary struct {
	Name       string `json:"name" yaml:"name" msgpack:"name" valid:"required"`
	Consumers  int64  `json:"consumers" yaml:"consumers" msgpack:"consumers" valid:"required"`
	Developers int64  `json:"developers" yaml:"developers" msgpack:"developers" valid:"required"`
	Services   int64  `json:"services" yaml:"services" msgpack:"services" valid:"required"`
	Plans      int64  `json:"plans" yaml:"plans" msgpack:"plans" valid:"required"`
	Providers  int64  `json:"providers" yaml:"providers" msgpack:"providers" valid:"required"`
	Users      int64  `json:"users" yaml:"users" msgpack:"users" valid:"required"`
	Roles      int64  `json:"roles" yaml:"roles" msgpack:"roles" valid:"required"`
}

type Token struct {
	Entity        `msgpack:",inline" mapstructure:",squash"`
	CredentialID  string                 `json:"credentialId" yaml:"credentialId" msgpack:"credentialId"`
	Scheme        string                 `json:"scheme" yaml:"scheme" msgpack:"scheme"`
	TokenType     string                 `json:"tokenType" yaml:"tokenType" msgpack:"tokenType"`
	GrantType     string                 `json:"grantType" yaml:"grantType" msgpack:"grantType"`
	Issuer        string                 `json:"issuer" yaml:"issuer" msgpack:"issuer"`
	Subject       string                 `json:"subject" yaml:"subject" msgpack:"subject"`
	Audience      []string               `json:"audience" yaml:"audience" msgpack:"audience"`
	IssuedAt      time.Time              `json:"issuedAt" yaml:"issuedAt" msgpack:"issuedAt"`
	Expiry        int64                  `json:"expiry" yaml:"expiry" msgpack:"expiry"`
	Lifespan      Lifespan               `json:"lifespan" yaml:"lifespan" msgpack:"lifespan"`
	FromToken     string                 `json:"fromToken" yaml:"fromToken" msgpack:"fromToken"`
	RefreshToken  string                 `json:"refreshToken" yaml:"refreshToken" msgpack:"refreshToken"`
	PermissionIds []string               `json:"permissionIds" yaml:"permissionIds" msgpack:"permissionIds"`
	Claims        map[string]interface{} `json:"claims" yaml:"claims" msgpack:"claims"`
	State         string                 `json:"state" yaml:"state" msgpack:"state"`
	URI           string                 `json:"uri" yaml:"uri" msgpack:"uri"`
	Extended      map[string]interface{} `json:"extended" yaml:"extended" msgpack:"extended"`
	ExternalID    string                 `json:"externalId" yaml:"externalId" msgpack:"externalId"`
}

type Upstream struct {
	Type     string                            `json:"type" yaml:"type" msgpack:"type" valid:"required"`
	Versions map[string]map[string]interface{} `json:"versions" yaml:"versions" msgpack:"versions" valid:"required"`
}

type User struct {
	Entity       `msgpack:",inline" mapstructure:",squash"`
	UserUpdate   `msgpack:",inline" mapstructure:",squash"`
	UserReadOnly `msgpack:",inline" mapstructure:",squash"`
	Auditable    `msgpack:",inline" mapstructure:",squash"`
}

type UserPermissions struct {
	Name          string            `json:"name" yaml:"name" msgpack:"name" valid:"required"`
	Administrator bool              `json:"administrator" yaml:"administrator" msgpack:"administrator" valid:"required"`
	Permissions   []string          `json:"permissions" yaml:"permissions" msgpack:"permissions" valid:"required"`
	AccessLevels  map[string]string `json:"accessLevels" yaml:"accessLevels" msgpack:"accessLevels" valid:"required"`
	UsersLocked   bool              `json:"usersLocked" yaml:"usersLocked" msgpack:"usersLocked" valid:"required"`
}

type UserReadOnly struct {
	Administrator bool  `json:"administrator" yaml:"administrator" msgpack:"administrator" valid:"required"`
	Role          *Role `json:"role" yaml:"role" msgpack:"role"`
}

type UserUpdate struct {
	UserName   string  `json:"userName" yaml:"userName" msgpack:"userName" valid:"required"`
	FirstName  string  `json:"firstName" yaml:"firstName" msgpack:"firstName" valid:"required"`
	LastName   string  `json:"lastName" yaml:"lastName" msgpack:"lastName" valid:"required"`
	RoleID     *string `json:"roleId" yaml:"roleId" msgpack:"roleId"`
	ExternalID *string `json:"externalId" yaml:"externalId" msgpack:"externalId"`
}

// AuthenticationType is an enum.
// It is serialized as a string in JSON.
type AuthenticationType string

// AuthenticationType enum values.
var (
	AuthenticationTypeNone        AuthenticationType = "none"
	AuthenticationTypeTwoLegged   AuthenticationType = "two-legged"
	AuthenticationTypeThreeLegged AuthenticationType = "three-legged"

	AuthenticationTypeMap = map[string]int{
		string(AuthenticationTypeNone):        1,
		string(AuthenticationTypeTwoLegged):   2,
		string(AuthenticationTypeThreeLegged): 3,
	}

	AuthenticationTypeIntMap = map[int]string{
		1: string(AuthenticationTypeNone),
		2: string(AuthenticationTypeTwoLegged),
		3: string(AuthenticationTypeThreeLegged),
	}
)

// Int converts this to a nullable integer.
func (s *AuthenticationType) Int() sql.NullInt64 {
	ni := sql.NullInt64{}
	if s == nil {
		ni.Int64, ni.Valid = 0, false
	} else {
		num, _ := AuthenticationTypeMap[string(*s)]
		ni.Int64, ni.Valid = int64(num), true
	}
	return ni
}

// Known says whether or not this value is a known enum value.
func (s *AuthenticationType) Known() bool {
	if s == nil {
		return false
	}
	_, ok := AuthenticationTypeMap[string(*s)]
	return ok
}

// UnmarshalJSON satisfies the json.Unmarshaler
func (s *AuthenticationType) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*s = AuthenticationType(str)
	if !s.Known() {
		return fmt.Errorf("Unknown AuthenticationType enum value: %s", s)
	}
	return nil
}

// String is for the standard stringer interface.
func (s *AuthenticationType) String() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

// Lifespan is an enum.
// It is serialized as a string in JSON.
type Lifespan string

// Lifespan enum values.
var (
	LifespanFinite  Lifespan = "finite"
	LifespanSession Lifespan = "session"

	LifespanMap = map[string]int{
		string(LifespanFinite):  1,
		string(LifespanSession): 2,
	}

	LifespanIntMap = map[int]string{
		1: string(LifespanFinite),
		2: string(LifespanSession),
	}
)

// Int converts this to a nullable integer.
func (s *Lifespan) Int() sql.NullInt64 {
	ni := sql.NullInt64{}
	if s == nil {
		ni.Int64, ni.Valid = 0, false
	} else {
		num, _ := LifespanMap[string(*s)]
		ni.Int64, ni.Valid = int64(num), true
	}
	return ni
}

// Known says whether or not this value is a known enum value.
func (s *Lifespan) Known() bool {
	if s == nil {
		return false
	}
	_, ok := LifespanMap[string(*s)]
	return ok
}

// UnmarshalJSON satisfies the json.Unmarshaler
func (s *Lifespan) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*s = Lifespan(str)
	if !s.Known() {
		return fmt.Errorf("Unknown Lifespan enum value: %s", s)
	}
	return nil
}

// String is for the standard stringer interface.
func (s *Lifespan) String() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

// Method is an enum.
// It is serialized as a string in JSON.
type Method string

// Method enum values.
var (
	MethodGET     Method = "GET"
	MethodPOST    Method = "POST"
	MethodPUT     Method = "PUT"
	MethodPATCH   Method = "PATCH"
	MethodDELETE  Method = "DELETE"
	MethodHEAD    Method = "HEAD"
	MethodOPTIONS Method = "OPTIONS"

	MethodMap = map[string]int{
		string(MethodGET):     1,
		string(MethodPOST):    2,
		string(MethodPUT):     3,
		string(MethodPATCH):   4,
		string(MethodDELETE):  5,
		string(MethodHEAD):    6,
		string(MethodOPTIONS): 7,
	}

	MethodIntMap = map[int]string{
		1: string(MethodGET),
		2: string(MethodPOST),
		3: string(MethodPUT),
		4: string(MethodPATCH),
		5: string(MethodDELETE),
		6: string(MethodHEAD),
		7: string(MethodOPTIONS),
	}
)

// Int converts this to a nullable integer.
func (s *Method) Int() sql.NullInt64 {
	ni := sql.NullInt64{}
	if s == nil {
		ni.Int64, ni.Valid = 0, false
	} else {
		num, _ := MethodMap[string(*s)]
		ni.Int64, ni.Valid = int64(num), true
	}
	return ni
}

// Known says whether or not this value is a known enum value.
func (s *Method) Known() bool {
	if s == nil {
		return false
	}
	_, ok := MethodMap[string(*s)]
	return ok
}

// UnmarshalJSON satisfies the json.Unmarshaler
func (s *Method) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*s = Method(str)
	if !s.Known() {
		return fmt.Errorf("Unknown Method enum value: %s", s)
	}
	return nil
}

// String is for the standard stringer interface.
func (s *Method) String() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

// PropertyType is an enum.
// It is serialized as a string in JSON.
type PropertyType string

// PropertyType enum values.
var (
	PropertyTypeString   PropertyType = "string"
	PropertyTypeInteger  PropertyType = "integer"
	PropertyTypeDecimal  PropertyType = "decimal"
	PropertyTypeBooleanV PropertyType = "booleanV"
	PropertyTypeDatetime PropertyType = "datetime"

	PropertyTypeMap = map[string]int{
		string(PropertyTypeString):   1,
		string(PropertyTypeInteger):  2,
		string(PropertyTypeDecimal):  3,
		string(PropertyTypeBooleanV): 4,
		string(PropertyTypeDatetime): 5,
	}

	PropertyTypeIntMap = map[int]string{
		1: string(PropertyTypeString),
		2: string(PropertyTypeInteger),
		3: string(PropertyTypeDecimal),
		4: string(PropertyTypeBooleanV),
		5: string(PropertyTypeDatetime),
	}
)

// Int converts this to a nullable integer.
func (s *PropertyType) Int() sql.NullInt64 {
	ni := sql.NullInt64{}
	if s == nil {
		ni.Int64, ni.Valid = 0, false
	} else {
		num, _ := PropertyTypeMap[string(*s)]
		ni.Int64, ni.Valid = int64(num), true
	}
	return ni
}

// Known says whether or not this value is a known enum value.
func (s *PropertyType) Known() bool {
	if s == nil {
		return false
	}
	_, ok := PropertyTypeMap[string(*s)]
	return ok
}

// UnmarshalJSON satisfies the json.Unmarshaler
func (s *PropertyType) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*s = PropertyType(str)
	if !s.Known() {
		return fmt.Errorf("Unknown PropertyType enum value: %s", s)
	}
	return nil
}

// String is for the standard stringer interface.
func (s *PropertyType) String() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

// Scope is an enum.
// It is serialized as a string in JSON.
type Scope string

// Scope enum values.
var (
	ScopeUser     Scope = "user"
	ScopeConsumer Scope = "consumer"
	ScopeBoth     Scope = "both"

	ScopeMap = map[string]int{
		string(ScopeUser):     1,
		string(ScopeConsumer): 2,
		string(ScopeBoth):     3,
	}

	ScopeIntMap = map[int]string{
		1: string(ScopeUser),
		2: string(ScopeConsumer),
		3: string(ScopeBoth),
	}
)

// Int converts this to a nullable integer.
func (s *Scope) Int() sql.NullInt64 {
	ni := sql.NullInt64{}
	if s == nil {
		ni.Int64, ni.Valid = 0, false
	} else {
		num, _ := ScopeMap[string(*s)]
		ni.Int64, ni.Valid = int64(num), true
	}
	return ni
}

// Known says whether or not this value is a known enum value.
func (s *Scope) Known() bool {
	if s == nil {
		return false
	}
	_, ok := ScopeMap[string(*s)]
	return ok
}

// UnmarshalJSON satisfies the json.Unmarshaler
func (s *Scope) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*s = Scope(str)
	if !s.Known() {
		return fmt.Errorf("Unknown Scope enum value: %s", s)
	}
	return nil
}

// String is for the standard stringer interface.
func (s *Scope) String() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

// TimeUnit is an enum.
// It is serialized as a string in JSON.
type TimeUnit string

// TimeUnit enum values.
var (
	TimeUnitMinute TimeUnit = "minute"
	TimeUnitHour   TimeUnit = "hour"
	TimeUnitDay    TimeUnit = "day"
	TimeUnitMonth  TimeUnit = "month"

	TimeUnitMap = map[string]int{
		string(TimeUnitMinute): 1,
		string(TimeUnitHour):   2,
		string(TimeUnitDay):    3,
		string(TimeUnitMonth):  4,
	}

	TimeUnitIntMap = map[int]string{
		1: string(TimeUnitMinute),
		2: string(TimeUnitHour),
		3: string(TimeUnitDay),
		4: string(TimeUnitMonth),
	}
)

// Int converts this to a nullable integer.
func (s *TimeUnit) Int() sql.NullInt64 {
	ni := sql.NullInt64{}
	if s == nil {
		ni.Int64, ni.Valid = 0, false
	} else {
		num, _ := TimeUnitMap[string(*s)]
		ni.Int64, ni.Valid = int64(num), true
	}
	return ni
}

// Known says whether or not this value is a known enum value.
func (s *TimeUnit) Known() bool {
	if s == nil {
		return false
	}
	_, ok := TimeUnitMap[string(*s)]
	return ok
}

// UnmarshalJSON satisfies the json.Unmarshaler
func (s *TimeUnit) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*s = TimeUnit(str)
	if !s.Known() {
		return fmt.Errorf("Unknown TimeUnit enum value: %s", s)
	}
	return nil
}

// String is for the standard stringer interface.
func (s *TimeUnit) String() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

// Type is an enum.
// It is serialized as a string in JSON.
type Type string

// Type enum values.
var (
	TypeAction Type = "action"
	TypeEntity Type = "entity"

	TypeMap = map[string]int{
		string(TypeAction): 1,
		string(TypeEntity): 2,
	}

	TypeIntMap = map[int]string{
		1: string(TypeAction),
		2: string(TypeEntity),
	}
)

// Int converts this to a nullable integer.
func (s *Type) Int() sql.NullInt64 {
	ni := sql.NullInt64{}
	if s == nil {
		ni.Int64, ni.Valid = 0, false
	} else {
		num, _ := TypeMap[string(*s)]
		ni.Int64, ni.Valid = int64(num), true
	}
	return ni
}

// Known says whether or not this value is a known enum value.
func (s *Type) Known() bool {
	if s == nil {
		return false
	}
	_, ok := TypeMap[string(*s)]
	return ok
}

// UnmarshalJSON satisfies the json.Unmarshaler
func (s *Type) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*s = Type(str)
	if !s.Known() {
		return fmt.Errorf("Unknown Type enum value: %s", s)
	}
	return nil
}

// String is for the standard stringer interface.
func (s *Type) String() string {
	if s == nil {
		return ""
	}
	return string(*s)
}
