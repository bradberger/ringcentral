package ringcentral

import "time"

type MessageStoreEventPayload struct {
	Timestamp      time.Time `json:"timestamp"`
	UUID           string    `json:"uuid"`
	Event          string    `json:"event"`
	SubscriptionID string    `json:"subscriptionId"`
	Body           struct {
		ExtensionID int64     `json:"extensionId"`
		LastUpdated time.Time `json:"lastUpdated"`
		Changes     []struct {
			Type         string `json:"type"`
			UpdatedCount int    `json:"updatedCount"`
			NewCount     int    `json:"newCount"`
		}
	} `json:"body"`
}

type Attachment struct {
	URI         string `json:"uri"`
	ID          string `json:"id"`
	Type        string `json:"type"`
	ContentType string `json:"contentType"`
	Size        int    `json:"size"`
}

type InboundMessageEvent struct {
	UUID           string `json:"uuid"`
	Event          string `json:"event"`
	SubscriptionID string `json:"subscriptionId"`
	Body           struct {
		ID string `json:"id"`
		To []struct {
			PhoneNumber string `json:"phoneNumber"`
			Location    string `json:"location"`
		} `json:"to"`
		From struct {
			PhoneNumber string `json:"phoneNumber"`
			Name        string `json:"name"`
		} `json:"from"`
		Type             string       `json:"type"`
		CreationTime     time.Time    `json:"creationTime"`
		LastModifiedTime time.Time    `json:"lastModifiedTime"`
		ReadStatus       string       `json:"readStatus"`
		Priority         string       `json:"priority"`
		Attachments      []Attachment `json:"attachments"`
		Direction        Direction    `json:"direction"`
		Availability     string       `json:"availability"`
		Subject          string       `json:"subject"`
		MessageStatus    string       `json:"messageStatus"`
		ConversationID   string       `json:"conversationId"`
	} `json:"body"`
}

type TelephonyStatus string

const (
	TelephonyStatusNoCall        TelephonyStatus = "NoCall"
	TelephonyStatusCallConnected TelephonyStatus = "CallConnected"
	TelephonyStatusRinging       TelephonyStatus = "Ringing"
	TelephonyStatusOnHold        TelephonyStatus = "OnHold"
	TelephonyStatusParkedCall    TelephonyStatus = "ParkedCall"
)

type TerminationType string

const (
	TerminationTypeFinal        TerminationType = "final"
	TerminationTypeIntermediate TerminationType = "intermediate"
)

type PresenceStatus string

const (
	PresenceStatusOffline   PresenceStatus = "Offline"
	PresenceStatusBusy      PresenceStatus = "Busy"
	PresenceStatusAvailable PresenceStatus = "Available"
)

type AccountPresenceEvent struct {
	UUID           string        `json:"uuid"`
	Event          string        `json:"event"`
	SubscriptionID string        `json:"subscriptionId"`
	Timestamp      time.Time     `json:"timestamp"`
	Body           PresenceEvent `json:"body"`
}

type DNDStatus string

const (
	DNDStatusTakeAllCalls               DNDStatus = "TakeAllCalls"
	DNDStatusDoNotAcceptAnyCalls        DNDStatus = "DoNotAcceptAnyCalls"
	DNDStatusDoNotAcceptDepartmentCalls DNDStatus = "DoNotAcceptDeparmentCalls"
	DNDStatusTakeDepartmentCallsOnly    DNDStatus = "TakeDepartmentCallsOnly"
)

type PresenceEvent struct {
	ExtensionID         string          `json"extensionId"`
	TelephonyStatus     TelephonyStatus `json:telephonyStatus"`
	TerminationType     TerminationType `json:"terminationType"`
	Sequence            int             `json:"sequence"`
	PresenceStatus      PresenceStatus  `json:"presenceStatus"`
	UserStatus          PresenceStatus  `json:"userStatus"`
	DNDStatus           DNDStatus       `json:"dndStatus"`
	AllowSeeMyPresence  bool            `json:"allowSeeMyPresence"`
	RingOnMonitoredCall bool            `json:"ringOnMonitoredCall"`
	PickUpCallsOnHold   bool            `json:"pickUpCallsOnHold"`
}

type DetailedPresenceEvent struct {
	ExtensionID         string          `json"extensionId"`
	TelephonyStatus     TelephonyStatus `json:telephonyStatus"`
	TerminationType     TerminationType `json:"terminationType"`
	Sequence            int             `json:"sequence"`
	PresenceStatus      PresenceStatus  `json:"presenceStatus"`
	UserStatus          PresenceStatus  `json:"userStatus"`
	DNDStatus           DNDStatus       `json:"dndStatus"`
	AllowSeeMyPresence  bool            `json:"allowSeeMyPresence"`
	RingOnMonitoredCall bool            `json:"ringOnMonitoredCall"`
	PickUpCallsOnHold   bool            `json:"pickUpCallsOnHold"`
	ActiveCalls         []ActiveCall    `json:"activeCalls"`
}

type ExtensionPresenceEvent struct {
	UUID           string        `json:"uuid"`
	Event          string        `json:"event"`
	SubscriptionID string        `json:"subscriptionId"`
	Timestamp      time.Time     `json:"timestamp"`
	Body           PresenceEvent `json:"body"`
}

type DetailedExtensionPresenceEvent struct {
}

type ExtensionPresenceLineEvent struct {
	UUID           string    `json:"uuid"`
	Event          string    `json:"event"`
	SubscriptionID string    `json:"subscriptionId"`
	Timestamp      time.Time `json:"timestamp"`
	Body           PresenceLineEvent
}

type PresenceLineEvent struct {
	Extension ExtensionInfo `json:"extension"`
	Sequence  int           `json:"sequence"`
}

type ActiveCall struct {
	ID              string          `json:"id"`
	Direction       Direction       `json:"direction"`
	From            string          `json:"from"`
	To              string          `json:"to"`
	TelephonyStatus TelephonyStatus `json:"telephonyStatus"`
	SessionID       string          `json:"sessionId"`
}
