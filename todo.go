package ringcentral

import "time"

const (
	VoiceNameTypeDefault      VoiceNameType = "Default"
	VoiceNameTypeTextToSpeech VoiceNameType = "TextToSpeech"
	VoiceNameTypeRecorded     VoiceNameType = "Recorded"
)

const (
	TimeFormat12H = "12h"
	TimeFormat24H = "24h"
)

type VoiceNameType string

type DepartmentInfo struct {
	ID              string `json:"id"`
	URI             string `json:"uri"`
	ExtensionNumber string `json:"extensionNumber"`
}

type AccountInfo struct {
	ID  string `json:"id"`
	URI string `json:"uri"`
}

type ContactInfo struct {
	FirstName        string             `json:"firstName"`
	LastName         string             `json:"lastName"`
	Company          string             `json:"company"`
	Email            string             `json:"email"`
	BusinessPhone    string             `json:"businessPhone"`
	BusinessAddress  ContactAddressInfo `json:"businessAddress"`
	EmailAsLoginName bool               `json:"emailAsLoginName"`
	PronouncedName   PronouncedNameInfo `json:"pronouncedName"`
	Department       string             `json:"department"`
}

type ContactAddressInfo struct {
	Country string `json:"country"`
	State   string `json:"state"`
	City    string `json:"city"`
	Street  string `json:"street"`
	Zip     string `json:"zip"`
}

type PronouncedNameInfo struct {
	Type VoiceNameType `json:"type"`
	Text string        `json:"text"`
}

type APIVersion struct {
	URI           string    `json:"uri"`
	VersionString string    `json:"versionString"`
	ReleaseDate   time.Time `json:"releaseDate"`
	URIString     string    `json:"uriString"`
}

type ExtendedPermissions struct{}
type ProfileImageInfo struct{}
type ReferenceInfo struct{}

type ExtensionServiceFeatureInfo struct {
	Enabled     bool   `json:"enabled"`
	FeatureName string `json:"featureName"`
	Reason      string `json:"reason"`
}

type CallQueueInfo struct {
	SLAGoal                    int64 `json:"slaGoal"`
	SLAThresholdSeconds        int64 `json:"slaThresholdSeconds"`
	IncludeAbandondedCalls     bool  `json:"includeAbandondedCalls"`
	AbandondedThresholdSeconds int64 `json:"abandondedThresholdSeconds"`
}

type RegionalSettings struct {
	HomeCountry      CountryInfo          `json:"homeCountry"`
	Timezone         TimezoneInfo         `json:"timezone"`
	Language         LanguageInfo         `json:"language"`
	GreetingLanguage GreetingLanguageInfo `json:"greetingLanguage"`
	FormattingLocale FormattingLocaleInfo `json:"formattingLocale"`
	TimeFormat       TimeFormat           `json:"timeFormat"`
}

type CountryInfo struct{}
type TimezoneInfo struct {
	ID          string `json:"id"`
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type LanguageInfo struct {
	ID               string `json:"id"`
	URI              string `json:"uri"`
	Greeting         bool   `json:"greeting"`
	FormattingLocale bool   `json:"formattingLocale"`
	LocaleCode       string `json:"localeCode"`
	Name             string `json:"name"`
	UI               bool   `json:"ui"`
}
type GreetingLanguageInfo struct {
	ID         string `json:"id"`
	LocaleCode string `json:"localeCode"`
	Name       string `json:"name"`
}
type FormattingLocaleInfo struct {
	ID         string `json:"id"`
	LocaleCode string `json:"localeCode"`
	Name       string `json:"name"`
}
type TimeFormat string
