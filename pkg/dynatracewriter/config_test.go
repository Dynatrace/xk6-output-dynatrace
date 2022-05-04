package dynatracewriter

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.k6.io/k6/lib/types"
	"gopkg.in/guregu/null.v3"
)

func TestApply(t *testing.T) {
	t.Parallel()

	fullConfig := Config{
		Url:                   null.StringFrom("some-url"),
		InsecureSkipTLSVerify: null.BoolFrom(false),
		CACert:                null.StringFrom("some-file"),
		ApiToken:                  null.StringFrom("user"),
		FlushPeriod:           types.NullDurationFrom(10 * time.Second),
		Headers: map[string]string{
			"X-Header": "value",
		},
	}

	// Defaults should be overwritten by valid values
	c := NewConfig()
	c = c.Apply(fullConfig)
	assert.Equal(t, fullConfig.Url, c.Url)
	assert.Equal(t, fullConfig.InsecureSkipTLSVerify, c.InsecureSkipTLSVerify)
	assert.Equal(t, fullConfig.CACert, c.CACert)
	assert.Equal(t, fullConfig.ApiToken, c.ApiToken)
	assert.Equal(t, fullConfig.FlushPeriod, c.FlushPeriod)
	assert.Equal(t, fullConfig.Headers, c.Headers)

	// Defaults shouldn't be impacted by invalid values
	c = NewConfig()
	c = c.Apply(Config{
		ApiToken:                  null.NewString("user", false),
		InsecureSkipTLSVerify: null.NewBool(false, false),
	})
	assert.Equal(t, false, c.ApiToken.Valid)
	assert.Equal(t, true, c.InsecureSkipTLSVerify.Valid)
}

func TestConfigParseArg(t *testing.T) {
	t.Parallel()

	c, err := ParseArg("url=https://bix24852.dev.dynatracelabs.com")
	assert.Nil(t, err)
	assert.Equal(t, null.StringFrom("https://bix24852.dev.dynatracelabs.com/"), c.Url)

	c, err = ParseArg("url=https://bix24852.dev.dynatracelabs.com,insecureSkipTLSVerify=false")
	assert.Nil(t, err)
	assert.Equal(t, null.StringFrom("https://bix24852.dev.dynatracelabs.com"), c.Url)
	assert.Equal(t, null.BoolFrom(false), c.InsecureSkipTLSVerify)

	c, err = ParseArg("url=https://bix24852.dev.dynatracelabs.com,caCertFile=f.crt")
	assert.Nil(t, err)
	assert.Equal(t, null.StringFrom("https://bix24852.dev.dynatracelabs.com"), c.Url)
	assert.Equal(t, null.StringFrom("f.crt"), c.CACert)

	c, err = ParseArg("url=https://bix24852.dev.dynatracelabs.com,insecureSkipTLSVerify=false,caCertFile=f.crt,apitoken=dede")
	assert.Nil(t, err)
	assert.Equal(t, null.StringFrom("https://bix24852.dev.dynatracelabs.com"), c.Url)
	assert.Equal(t, null.BoolFrom(false), c.InsecureSkipTLSVerify)
	assert.Equal(t, null.StringFrom("f.crt"), c.CACert)
	assert.Equal(t, null.StringFrom("apitoken"), c.ApiToken)


	c, err = ParseArg("url=https://bix24852.dev.dynatracelabs.com,flushPeriod=2s")
	assert.Nil(t, err)
	assert.Equal(t, null.StringFrom("https://bix24852.dev.dynatracelabs.com"), c.Url)
	assert.Equal(t, types.NullDurationFrom(time.Second*2), c.FlushPeriod)

	c, err = ParseArg("url=https://bix24852.dev.dynatracelabs.com,headers.X-Header=value")
	assert.Nil(t, err)
	assert.Equal(t, null.StringFrom("http://prometheus.remote:3412/write"), c.Url)
	assert.Equal(t, map[string]string{"X-Header": "value"}, c.Headers)
}

// testing both GetConsolidatedConfig and ConstructRemoteConfig here until it's future config refactor takes shape (k6 #883)
func TestConstructRemoteConfig(t *testing.T) {
	u, _ := url.Parse("https://bix24852.dev.dynatracelabs.com")

	t.Parallel()

	testCases := map[string]struct {
		jsonRaw      json.RawMessage
		env          map[string]string
		arg          string
		config       Config
		errString    string
		remoteConfig *remote.ClientConfig
	}{
		"json_success": {
			jsonRaw: json.RawMessage(fmt.Sprintf(`{"url":"%s"}`, u.String())),
			env:     nil,
			arg:     "",
			config: Config{
				Url:                   null.StringFrom(u.String()),
				InsecureSkipTLSVerify: null.BoolFrom(true),
				CACert:                null.NewString("", false),
				ApiToken:              null.NewString("", false),
				FlushPeriod:           types.NullDurationFrom(defaultFlushPeriod),
				KeepTags:              null.BoolFrom(true),
				KeepNameTag:           null.BoolFrom(false),
				KeepUrlTag:            null.BoolFrom(true),
				Headers:               make(map[string]string),
			},
			errString: "",

		},
		"mixed_success": {
			jsonRaw: json.RawMessage(fmt.Sprintf(`{"url":"%s"}`, u.String())),
			env:     map[string]string{"K6_DYNATRACE_INSECURE_SKIP_TLS_VERIFY": "false", "K6_DYNATRACE_APITOKEN": "u"},
			arg:     "apitoken=user",
			config: Config{
				Url:                   null.StringFrom(u.String()),
				InsecureSkipTLSVerify: null.BoolFrom(false),
				CACert:                null.NewString("", false),
                ApiToken:              null.NewString("apitoken", true),
				FlushPeriod:           types.NullDurationFrom(defaultFlushPeriod),
				KeepTags:              null.BoolFrom(true),
				KeepNameTag:           null.BoolFrom(false),
				KeepUrlTag:            null.BoolFrom(true),
				Headers:               make(map[string]string),
			},
			errString: "",

		},
		"invalid_duration": {
			jsonRaw:      json.RawMessage(fmt.Sprintf(`{"url":"%s"}`, u.String())),
			env:          map[string]string{"K6_DYNATRACE_FLUSH_PERIOD": "d"},
			arg:          "",
			config:       Config{},
			errString:    "strconv.ParseInt",
			remoteConfig: nil,
		},
		"invalid_insecureSkipTLSVerify": {
			jsonRaw:      json.RawMessage(fmt.Sprintf(`{"url":"%s"}`, u.String())),
			env:          map[string]string{"K6_DYNATRACE_INSECURE_SKIP_TLS_VERIFY": "d"},
			arg:          "",
			config:       Config{},
			errString:    "strconv.ParseBool",
			remoteConfig: nil,
		},
		"remote_write_with_headers_json": {
			jsonRaw: json.RawMessage(fmt.Sprintf(`{"url":"%s", "headers":{"X-Header":"value"}}`, u.String())),
			env:     nil,
			arg:     "",
			config: Config{
				Url:                   null.StringFrom(u.String()),
				InsecureSkipTLSVerify: null.BoolFrom(true),
				CACert:                null.NewString("", false),
				ApiToken:                  null.NewString("", false),
				FlushPeriod:           types.NullDurationFrom(defaultFlushPeriod),
				KeepTags:              null.BoolFrom(true),
				KeepNameTag:           null.BoolFrom(false),
				KeepUrlTag:            null.BoolFrom(true),
				Headers: map[string]string{
					"X-Header": "value",
				},
			},
			errString: "",

		},
		"remote_write_with_headers_env": {
			jsonRaw: json.RawMessage(fmt.Sprintf(`{"url":"%s", "headers":{"X-Header":"value"}}`, u.String())),
			env: map[string]string{
				"K6_DYNATRACE_HEADER": "value_from_env",
			},
			arg: "",
			config: Config{
				Url:                   null.StringFrom(u.String()),
				InsecureSkipTLSVerify: null.BoolFrom(true),
				CACert:                null.NewString("", false),
				ApiToken:                  null.NewString("", false),
				FlushPeriod:           types.NullDurationFrom(defaultFlushPeriod),
				KeepTags:              null.BoolFrom(true),
				KeepNameTag:           null.BoolFrom(false),
				KeepUrlTag:            null.BoolFrom(true),
				Headers: map[string]string{
					"X-Header": "value_from_env",
				},
			},
			errString: "",

		},
		"remote_write_with_headers_arg": {
			jsonRaw: json.RawMessage(fmt.Sprintf(`{"url":"%s", "headers":{"X-Header":"value"}}`, u.String())),
			env: map[string]string{
				"K6_DYNATRACE_HEADER": "value_from_env",
			},
			arg: "headers.X-Header=value_from_arg",
			config: Config{
				Url:                   null.StringFrom(u.String()),
				InsecureSkipTLSVerify: null.BoolFrom(true),
				CACert:                null.NewString("", false),
				ApiToken:                  null.NewString("", false),
				FlushPeriod:           types.NullDurationFrom(defaultFlushPeriod),
				KeepTags:              null.BoolFrom(true),
				KeepNameTag:           null.BoolFrom(false),
				KeepUrlTag:            null.BoolFrom(true),
				Headers: map[string]string{
					"X-Header": "value_from_arg",
				},
			},
			errString: "",

			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			c, err := GetConsolidatedConfig(testCase.jsonRaw, testCase.env, testCase.arg)
			if len(testCase.errString) > 0 {
				assert.Contains(t, err.Error(), testCase.errString)
				return
			}
			assertConfig(t, c, testCase.config)

			// there can be error only on url.Parse at the moment so skipping that
			remoteConfig, _ := c.ConstructRemoteConfig()
			assertRemoteConfig(t, remoteConfig, testCase.remoteConfig)
		})
	}
}

func assertConfig(t *testing.T, actual, expected Config) {
	assert.Equal(t, expected.Mapping, actual.Mapping)
	assert.Equal(t, expected.Url, actual.Url)
	assert.Equal(t, expected.InsecureSkipTLSVerify, actual.InsecureSkipTLSVerify)
	assert.Equal(t, expected.CACert, actual.CACert)
	assert.Equal(t, expected.ApiToken, actual.ApiToken)
	assert.Equal(t, expected.FlushPeriod, actual.FlushPeriod)
	assert.Equal(t, expected.KeepTags, actual.KeepTags)
	assert.Equal(t, expected.KeepNameTag, expected.KeepNameTag)
	assert.Equal(t, expected.KeepUrlTag, expected.KeepUrlTag)
	assert.Equal(t, expected.Headers, actual.Headers)
}

