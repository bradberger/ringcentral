package ringcentral

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"
)

var (
	ctx   = context.Background()
	token *Token
	a     *API
)

func init() {
	a = New(os.Getenv("RINGCENTRAL_APP_ID"), os.Getenv("RINGCENTRAL_APP_SECRET"), "~")
	a.TestMode, _ = strconv.ParseBool(os.Getenv("RINGCENTRAL_TEST_MODE"))
}

func TestMakeURL(t *testing.T) {
	urlStr := a.makeURL("/foo/bar", nil)
	if a.TestMode {
		assert.Equal(t, "https://platform.devtest.ringcentral.com/foo/bar", urlStr)
	} else {
		assert.Equal(t, "https://platform.ringcentral.com/foo/bar", urlStr)
	}
}

// TestAuthorizeUsername should run first to get a valid auth token.
func TestAuthorizeUsername(t *testing.T) {

	var err error
	u := os.Getenv("RINGCENTRAL_PHONE_NUMBER")
	ext := os.Getenv("RINGCENTRAL_EXTENSION")
	pwd := os.Getenv("RINGCENTRAL_PASSWORD")

	token, err = a.Authorize(ctx, u, ext, pwd)
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, a.Token)
	assert.NotNil(t, token)
}

func TestGetExtensionList(t *testing.T) {

	if !assert.NotNil(t, a.Token) {
		t.Fail()
	}

	l, err := a.GetExtensionList(ctx, nil)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, l) {
		return
	}
	if !assert.NotNil(t, l.Records) {
		return
	}
	assert.True(t, len(l.Records) > 0)
}

func TestGetExtensionActiveCalls(t *testing.T) {

	l, err := a.GetExtensionList(ctx, nil)
	if !assert.NoError(t, err) || len(l.Records) < 1 {
		t.Fail()
		return
	}

	// Test the first returned extension
	active, err := a.ActiveCalls(ctx, l.Records[0].ID, nil)
	if !assert.NoError(t, err) {
		return
	}
	if len(active.Records) < 1 {
		t.Logf("No active calls for ext %s", l.Records[0].ExtensionNumber)
		return
	}
	for i := range active.Records {
		t.Logf("Active: %s %+v", active.Records[i])
	}
	return
}
