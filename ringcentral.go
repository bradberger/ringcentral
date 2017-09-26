package ringcentral

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
)

// Error messages
var (
	ErrNotAuthenticated = errors.New("ringcentral: not authenticated")
	ErrTokenExpired     = errors.New("ringcentral: token expired")
)

const (
	// Endpoint is the production RingCentral endpoint
	Endpoint = "https://platform.ringcentral.com"
	// EndpointTest is the RingCentral dev endpoint
	EndpointTest = "https://platform.devtest.ringcentral.com"
	userAgent    = "GoRingCentral/0.1 (github.com/bradberger/ringcentral)"
)

type CallAction string
type CallResult string
type Direction string
type Type string
type RecordingType string

// Transport is a transport type string
type Transport string

// Transport type definitions
const (
	TransportPSTN Transport = "PSTN"
	TransportVoIP Transport = "VoIP"
)

// RecordingTypes
const (
	RecordingTypeAutomatic RecordingType = "Automatic"
	RecordingTypeOnDemand  RecordingType = "OnDemand"
)

// Directions
const (
	Outbound     Direction = "Outbound"
	Inbound      Direction = "Inbound"
	AnyDirection Direction = "Inbound|Outbound"
)

// Types
const (
	Voice Type = "Voice"
	Fax   Type = "Fax"
)

// Call Actions
const (
	CallActionUnknown       CallAction = "Unknown"
	CallActionPhoneCall     CallAction = "Phone Call"
	CallActionPhoneLogin    CallAction = "Phone Login"
	CallActionIncomingFax   CallAction = "Incoming Fax"
	CallActionAcceptCall    CallAction = "Accept Call"
	CallActionFindMe        CallAction = "FindMe"
	CallActionFollowMe      CallAction = "FollowMe"
	CallActionOutgoingFax   CallAction = "Outgoing Fax"
	CallActionCallReturn    CallAction = "Call Return"
	CallActionCallingCard   CallAction = "Calling Card"
	CallActionRingDirectly  CallAction = "Ring Directly"
	CallActionRingOutWeb    CallAction = "RingOut Web"
	CallActionVoIPCall      CallAction = "VoIP Call"
	CallActionRingOutPC     CallAction = "RingOut PC"
	CallActionRingMe        CallAction = "RingMe"
	CallActionTransfer      CallAction = "Transfer"
	CallAction411Info       CallAction = "411 Info"
	CallActionEmergency     CallAction = "Emergency"
	CallActionE911Update    CallAction = "E911 Update"
	CallActionSupport       CallAction = "Support"
	CallActionRingOutMobile CallAction = "RingOut Mobile"
)

// Call Results
const (
	CallResultUnknown                  CallResult = "Unknown"
	CallResultInProgress               CallResult = "InProgress"
	CallResultMissed                   CallResult = "Missed"
	CallResultCallAccepted             CallResult = "Call accepted"
	CallResultVoicemail                CallResult = "Voicemail"
	CallResultRejected                 CallResult = "Rejected"
	CallResultReply                    CallResult = "Reply"
	CallResultReceived                 CallResult = "Received"
	CallResultReceiveError             CallResult = "Receive Error"
	CallResultFaxOnDemand              CallResult = "Fax on Demand"
	CallResultPartialReceive           CallResult = "Partial Receive"
	CallResultBlocked                  CallResult = "Blocked"
	CallResultCallConnected            CallResult = "Call connected"
	CallResultNoAnswer                 CallResult = "No Answer"
	CallResultInternationalDisabled    CallResult = "International Disabled"
	CallResultBusy                     CallResult = "Busy"
	CallResultSendError                CallResult = "Send Error"
	CallResultSent                     CallResult = "Sent"
	CallResultNoFaxMachine             CallResult = "No fax machine"
	CallResultResultEmpty              CallResult = "ResultEmpty"
	CallResultAccount                  CallResult = "Account"
	CallResultSuspended                CallResult = "Suspended"
	CallResultCallFailed               CallResult = "Call Failed"
	CallResultCallFailure              CallResult = "Call Failure"
	CallResultInternalError            CallResult = "Internal Error"
	CallResultIPPhoneOffline           CallResult = "IP Phone offline"
	CallResultRestrictedNumber         CallResult = "Restricted Number"
	CallResultWrongNumber              CallResult = "Wrong Number"
	CallResultStopped                  CallResult = "Stopped"
	CallResultHangUp                   CallResult = "Hang up"
	CallResultPoorLineQuality          CallResult = "Poor Line Quality"
	CallResultPartiallySent            CallResult = "Partially Sent"
	CallResultInternationalRestriction CallResult = "International Restriction"
	CallResultAbandoned                CallResult = "Abandoned"
	CallResultDeclined                 CallResult = "Declined"
	CallResultFaxReceiptError          CallResult = "Fax Receipt Error"
	CallResultFaxSendError             CallResult = "Fax Send Error"
)

var (
	TestMode, _           = strconv.ParseBool(os.Getenv("RINGCENTRAL_TEST_MODE"))
	defaultRequestTimeout = time.Second * 10

	api *API
)

type URI struct {
	URI string `json:"uri"`
}

func (u URI) Parse() (*url.URL, error) {
	return url.Parse(u.URI)
}

type CallerInfo struct {
	PhoneNumber     string     `json:"phoneNumber"`
	ExtensionNumber string     `json:"extensionNumber"`
	Location        string     `json:"location"`
	Name            string     `json:"name"`
	Device          DeviceInfo `json:"device"`
}

type DeviceInfo struct {
	ID  string `json:"id"`
	URI string `json:"uri"`
}

type VoicemailMessageInfo struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	URI  string `json:"uri"`
}

type BillingInfo struct {
	CostIncluded  int64 `json:"costIncluded"`
	CostPurchased int64 `json:"costPurchased"`
}

type RecordingInfo struct {
	ID   string        `json:"id"`
	URI  string        `json:"uri"`
	Type RecordingType `json:"type"`
}

type API struct {
	TestMode                    bool
	Timeout                     time.Duration
	Token                       *Token
	AccountID, AppID, AppSecret string

	lastRequest  *http.Request
	lastResponse *http.Response
	sync.RWMutex
}

type Response struct {
	URI string `json:"uri"`
}

type Paging struct {
	Page          int `json:"page"`
	PerPage       int `json:"perPage"`
	PageStart     int `json:"pageStart"`
	PageEnd       int `json:"pageEnd"`
	TotalPages    int `json:"totalPages"`
	TotalElements int `json:"totalElements"`
}

type Records struct {
	Action CallAction `json:"action"`
}

type LegInfo struct {
	Action    CallAction    `json:"action"`
	Direction Direction     `json:"direction"`
	Duration  int           `json:"duration"`
	Extension ExtensionInfo `json:"extension"`
	LegType   string        `json:"legType"`
	StartTime time.Time     `json:"startTime"`
	Type      Type          `json:"type"`
	Result    CallResult    `json:"result"`
	From      CallerInfo    `json:"from"`
	To        CallerInfo    `json:"to"`
	Transport Transport     `json:"transport"`
	Billing   BillingInfo   `json:"billing"`
	Recording RecordingInfo `json:"recording"`
}

type ExtensionInfo struct {
	ID               int64                         `json:"id"`
	URI              string                        `json:"uri"`
	Contact          ContactInfo                   `json:"contact"`
	Departments      DepartmentInfo                `json:"department"`
	ExtensionNumber  string                        `json:"extensionNumber"`
	Account          AccountInfo                   `json:"account"`
	Name             string                        `json:"name"`
	PartnerID        string                        `json:"partnerId"`
	Permissions      ExtendedPermissions           `json:"permissions"`
	ProfileImage     ProfileImageInfo              `json:"profileImage"`
	References       []ReferenceInfo               `json:"references"`
	RegionalSettings RegionalSettings              `json:"regionalSettings"`
	ServiceFeatures  []ExtensionServiceFeatureInfo `json:"serviceFeatures"`
	SetupWizardState SetupWizardState              `json:"setupWizardState"`
	Status           ExtensionStatus               `json:"status"`
	StatusInfo       StatusInfo                    `json:"statusInfo"`
	Type             ExtensionType                 `json:"type"`
	CallQueueInfo    CallQueueInfo                 `json:"callQueueInfo"`
}

type StatusInfo struct {
	Comment string `json:"comment"`
	Reason  string `json:"reason"`
}

type ExtensionType string

const (
	ExtensionTypeUser                 ExtensionType = "User"
	ExtensionTypeFaxUser              ExtensionType = "FaxUser"
	ExtensionTypeVirtualUser          ExtensionType = "VirtualUser"
	ExtensionTypeDigitalUser          ExtensionType = "DigitalUser"
	ExtensionTypeDepartment           ExtensionType = "Department"
	ExtensionTypeAnnouncement         ExtensionType = "Announcement"
	ExtensionTypeVoicemail            ExtensionType = "Voicemail"
	ExtensionTypeSharedLinesGroup     ExtensionType = "SharedLinesGroup"
	ExtensionTypePagingOnly           ExtensionType = "PagingOnly"
	ExtensionTypeIvrMenu              ExtensionType = "IvrMenu"
	ExtensionTypeApplicationExtension ExtensionType = "ApplicationExtension"
	ExtensionTypeParkLocation         ExtensionType = "ParkLocation"
	ExtensionTypeLimited              ExtensionType = "Limited"
)

type ExtensionStatus string

const (
	ExtensionStatusEnabled      ExtensionStatus = "Enabled"
	ExtensionStatusDisabled     ExtensionStatus = "Disabled"
	ExtensionStatusNotActivated ExtensionStatus = "NotActivated"
	ExtensionStatusUnassigned   ExtensionStatus = "Unassigned"
)

type SetupWizardState string

const (
	SetupWizardStateNotStarted SetupWizardState = "Not Started"
	SetupWizardStateIncomplete SetupWizardState = "Incomplete"
	SetupWizardStateCompleted  SetupWizardState = "Completed"
)

type ExtensionList struct {
	URI        string          `json:"uri"`
	Records    []ExtensionInfo `json:"records"`
	Navigation Navigation      `json:"navigation"`
	Paging     Paging          `json:"paging"`
}

type ExtensionActiveCalls struct {
	URI        string          `json:"uri"`
	Records    []CallLogRecord `json:"records"`
	Navigation Navigation      `json:"navigation"`
	Paging     Paging          `json:"paging"`
}

type Navigation struct {
	FirstPage    URI `json:"firstPage"`
	NextPage     URI `json:"nextPage"`
	PreviousPage URI `json:"previousPage"`
	LastPage     URI `json:"lastPage"`
}

type CallLogRecord struct {
	ID               string               `json:"id"`
	URI              string               `json:"uri"`
	SessionID        string               `json:"sessionId"`
	From             CallerInfo           `json:"from"`
	To               CallerInfo           `json:"to"`
	Message          VoicemailMessageInfo `json:"message"`
	Type             Type                 `json:"type"`
	Direction        Direction            `json:"direction"`
	Action           CallAction           `json:"action"`
	Result           CallResult           `json:"result"`
	Billing          BillingInfo          `json:"billingInfo"`
	StartTime        time.Time            `json:"startTime"`
	Duration         int                  `json:"duration"`
	Recording        RecordingInfo        `json:"recordingInfo"`
	LastModifiedTime time.Time            `json:"lastModifiedTime"`
	Legs             []LegInfo            `json:"legs"`
}

// Token is an OAuth2 token
type Token struct {
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	ExpiresIn             int64  `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
	OwnerID               string `json:"owner_id"`

	// Expires is calculated by the returned ExpiresIn field
	Expires time.Time
}

func (t *Token) setExpires() {
	t.Expires = time.Now().Add(time.Duration(t.ExpiresIn) * time.Second)
}

func (a *API) getClient(ctx context.Context) *http.Client {
	return getClient(ctx, a.Timeout)
}

func (a *API) getEndpoint() string {
	if a.TestMode {
		return EndpointTest
	}
	return Endpoint
}

func (a *API) makeURL(urlStr string, params url.Values) string {
	urlStr = fmt.Sprintf("%s/%s", a.getEndpoint(), strings.TrimPrefix(urlStr, "/"))
	if params == nil {
		return urlStr
	}
	urlStr = strings.Trim(urlStr, "&")
	enc := params.Encode()
	if strings.Index(urlStr, "?") > -1 {
		urlStr += "&" + enc
		return urlStr
	}
	return urlStr + "?" + enc
}

// Authorized returns true if there's a valid token for the API
func (a *API) Authorized(ctx context.Context) bool {
	return a.Token != nil && a.Token.Expires.After(time.Now())
}

func (a *API) Authorize(ctx context.Context, username, ext, pwd string) (*Token, error) {
	var t Token
	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("username", username)
	form.Add("extension", ext)
	form.Add("password", pwd)
	req, err := http.NewRequest(http.MethodPost, a.makeURL("/restapi/oauth/token", nil), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(a.AppID, a.AppSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if _, err := a.doRequest(ctx, req, &t); err != nil {
		return nil, err
	}
	a.Token = &t
	t.setExpires()
	return &t, nil
}

// PostForm sends a POST request to urlStr with an application/x-www-form-urlencoded body of form, marshaling the response into dstVal
func (a *API) PostForm(ctx context.Context, urlStr string, form url.Values, dstVal interface{}) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, a.makeURL(urlStr, nil), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return a.doRequest(ctx, req, dstVal)
}

func (a *API) doRequest(ctx context.Context, req *http.Request, dstVal interface{}) (*http.Response, error) {
	if ua := req.Header.Get("User-Agent"); ua == "" {
		req.Header.Set("User-Agent", userAgent)
	}

	// Check authentication. If basic auth is set, then use it.
	// Otherwise, check for a valid token, and if it exists set the Auth header
	_, _, ba := req.BasicAuth()

	switch {
	case ba:
	case a.Token == nil:
		return nil, ErrNotAuthenticated
	case a.Token.Expires.Before(time.Now()):
		return nil, ErrTokenExpired
	default:
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.Token.AccessToken))
	}

	client := a.getClient(ctx)
	resp, err := client.Do(req)

	// Set last request, response now. Keep locked for thread-safe access
	a.Lock()
	defer a.Unlock()
	a.lastRequest = req
	a.lastResponse = resp

	switch {
	case err != nil:
		return resp, err
	case resp.StatusCode >= 400:
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, fmt.Errorf("ringcentral: error reading response: %v", err)
		}
		if len(bodyBytes) < 1 {
			return resp, fmt.Errorf("ringcentral: api error: %s", resp.Status)
		}
		// Reset the body so it can be read again (debugging, etc.)
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// If can convert to an API error, then do so and return the message/description
		var e ErrorResponse
		if err := json.Unmarshal(bodyBytes, &e); err == nil {
			return resp, e
		}
		return resp, fmt.Errorf("ringcentral: error: %s", string(bodyBytes))
	case dstVal == nil:
		return resp, nil
	default:

		// Use this for debugging responses.
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, fmt.Errorf("ringcentral: error reading response body: %v", err)
		}
		// Reset the body so it can be read again (debugging, etc.)
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		// Unmarshal the body
		if err := json.Unmarshal(bodyBytes, dstVal); err != nil {
			return resp, fmt.Errorf("ringcentral: error decoding response: %v (response was: %s)", err, string(bodyBytes))
		}
		// Otherwise for better performance, decode using a buffer
		// if err := json.NewDecoder(resp.Body).Decode(dstVal); err != nil {
		// 	return resp, fmt.Errorf("ringcentral: error decoding response: %v (response was: %s)", err)
		// }
		return resp, nil
	}
}

// Post sends JSON encoded data to urlStr and marshals the response into dstVal
func (a *API) Post(ctx context.Context, urlStr string, data, dstVal interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return nil, fmt.Errorf("Could not encode data: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, a.makeURL(urlStr, nil), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return a.doRequest(ctx, req, dstVal)
}

// Put sends a JSON encoded PUT request to urlStr and marshals the response into dstVal
func (a *API) Put(ctx context.Context, urlStr string, data, dstVal interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return nil, fmt.Errorf("Could not encode data: %v", err)
	}
	req, err := http.NewRequest(http.MethodPut, a.makeURL(urlStr, nil), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return a.doRequest(ctx, req, dstVal)
}

// Get sends a GET request to the given urlStr, with optional query string defined in params
func (a *API) Get(ctx context.Context, urlStr string, params url.Values, dstVal interface{}) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, a.makeURL(urlStr, params), nil)
	if err != nil {
		return nil, err
	}
	return a.doRequest(ctx, req, dstVal)
}

// Delete sends a DELETE request to the given urlStr
func (a *API) Delete(ctx context.Context, urlStr string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, a.makeURL(urlStr, nil), nil)
	if err != nil {
		return nil, err
	}
	return a.doRequest(ctx, req, nil)
}

type Error struct {
	ErrorCode   string `json:"errorCode"`
	Code        string `json:"error"`
	Description string `json:"error_description"`
	Message     string `json:"message"`
}

type ErrorResponse struct {
	ErrorCode   string          `json:"errorCode"`
	Code        string          `json:"error"`
	Description string          `json:"error_description"`
	Message     string          `json:"message"`
	Errors      []ErrorResponse `json:"errors"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("[RingCentral API error] %s: %s %s", e.ErrorCode, e.Message, e.Description)
}

// GetExtensionList returns a list of all account extensions
func (a *API) GetExtensionList(ctx context.Context, params url.Values) (*ExtensionList, error) {
	var e ExtensionList
	if _, err := a.Get(ctx, fmt.Sprintf("/restapi/v1.0/account/%s/extension", a.AccountID), params, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

// ActiveCalls returns a list of active calls on the given extension
func (a *API) ActiveCalls(ctx context.Context, ext int64, params url.Values) (*ExtensionActiveCalls, error) {
	urlStr := fmt.Sprintf("/restapi/v1.0/account/%s/extension/%d/active-calls", a.AccountID, ext)
	var active *ExtensionActiveCalls
	if _, err := a.Get(ctx, urlStr, params, &active); err != nil {
		return nil, err
	}
	return active, nil
}

// LastRequest returns the last HTTP request sent via the API client. Use it for debugging.
func (a *API) LastRequest() *http.Request {
	a.RLock()
	defer a.RUnlock()
	return a.lastRequest
}

// LastResponse returns the last HTTP response recieved via the API client. Use it for debugging.
func (a *API) LastResponse() *http.Response {
	a.RLock()
	defer a.RUnlock()
	return a.lastResponse
}

// New creates a new API client
func New(appID, appSecret, accountID string) *API {
	if accountID == "" {
		accountID = "~"
	}
	return &API{AccountID: accountID, AppID: appID, AppSecret: appSecret}
}
