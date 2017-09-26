package ringcentral

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"
)

var (
	ErrInvalidSubscriptionID = errors.New("subscription id invalid")
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "Active"
	SubscriptionStatusSuspended SubscriptionStatus = "Suspended"
)

type TransportType string

const (
	TransportTypePubNum  TransportType = "PubNum"
	TransportTypeWebHook TransportType = "WebHook"

	SubscriptionMaxExipresIn = 604800
)

type DeliveryMode struct {
	TransportType       TransportType `json:"transportType,omitempty"`
	Encryption          bool          `json:"encryption,omitempty"`
	Address             string        `json:"address,omitempty"`
	SubscriberKey       string        `json:"subscriberKey,omitempty"`
	EncryptionAlgorithm string        `json:"encryptionAlgorithm,omitempty"`
	EncryptionKey       string        `json:"encryptionKey,omitempty"`
	RegistrationID      string        `json:"registrationId,omitempty"`
	CertificateName     string        `json:"certificateName,omitempty"`
}

type SubscriptionListResponse struct {
	URI     string             `json:"uri"`
	Records []SubscriptionInfo `json:"records"`
}

type CreateSubscriptionRequest struct {
	EventFilters []string     `json:"eventFilters"`
	DeliveryMode DeliveryMode `json:"deliveryMode"`
	ExpiresIn    int          `json:"expiresIn"`
}

type SubscriptionInfo struct {
	ID             string             `json:"id"`
	URI            string             `json:"uri"`
	EventFilters   []string           `json:"eventFilters"`
	ExpirationTime time.Time          `json:"expirationTime"`
	ExpiresIn      int                `json:"expiresIn"`
	Status         SubscriptionStatus `json:"status"`
	CreationTime   time.Time          `json:"creationTime"`
	DeliveryMode   DeliveryMode       `json:"deliveryMode"`
}

func (a *API) SubscriptionList(ctx context.Context) (*SubscriptionListResponse, error) {
	var list SubscriptionListResponse
	if _, err := a.Get(ctx, "/restapi/v1.0/subscription", nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// CreateSubscription creates a new subscrption.
func (a *API) CreateSubscription(ctx context.Context, req *CreateSubscriptionRequest) (*SubscriptionInfo, error) {
	if req.ExpiresIn < 1 || req.ExpiresIn > SubscriptionMaxExipresIn {
		req.ExpiresIn = SubscriptionMaxExipresIn
	}
	var s SubscriptionInfo
	if resp, err := a.Post(ctx, "/restapi/v1.0/subscription", req, &s); err != nil {
		fmt.Println(resp)
		return nil, err
	}
	return &s, nil
}

func (a *API) GetSubscription(ctx context.Context, sub string) (*SubscriptionInfo, error) {
	var s SubscriptionInfo
	id, err := getSubscriptionID(sub)
	if err != nil {
		return nil, err
	}
	if _, err := a.Get(ctx, "/restapi/v1.0/subscription/"+id, nil, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (a *API) UpdateSubscription(ctx context.Context, s *SubscriptionInfo, threshold, interval int) (*SubscriptionInfo, error) {
	var result SubscriptionInfo
	form := url.Values{}
	if threshold != 0 {
		form.Add("threshold", fmt.Sprintf("%d", threshold))
	}
	if interval != 0 {
		form.Add("interval", fmt.Sprintf("%d", interval))
	}
	urlStr := fmt.Sprintf("/restapi/v1.0/subscription/%s?%s", s.ID, form.Encode())
	if _, err := a.Put(ctx, urlStr, s, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *API) DeleteSubscription(ctx context.Context, sub interface{}) error {
	id, err := getSubscriptionID(sub)
	if err != nil {
		return err
	}
	_, err = a.Delete(ctx, "/restapi/v1.0/subscription/"+id)
	return err
}

func (a *API) RenewSubscription(ctx context.Context, sub interface{}) (*SubscriptionInfo, error) {
	var s SubscriptionInfo
	id, err := getSubscriptionID(sub)
	if err != nil {
		return nil, err
	}
	if _, err := a.Post(ctx, "/restapi/v1.0/subscription/"+id+"/renew", nil, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func getSubscriptionID(sub interface{}) (id string, err error) {
	switch sub.(type) {
	case string:
		id = sub.(string)
	case *string:
		id = *sub.(*string)
	case *SubscriptionInfo:
		id = sub.(*SubscriptionInfo).ID
	case SubscriptionInfo:
		id = sub.(SubscriptionInfo).ID
	default:
		err = ErrInvalidSubscriptionID
	}
	return
}
