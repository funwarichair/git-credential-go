package gitcredential_test

import (
	"reflect"
	"testing"

	. "github.com/funwarichair/git-credential-go"
	util "github.com/funwarichair/git-credential-go/util"
)

func intoPtr(base string) *string {
	return &base
}

func TestParse(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected GitCrendential
	}{
		{
			"empty",
			"\n\n",
			GitCrendential{
				Protocol: nil,
				Host:     nil,
				Path:     nil,
				Username: nil,
				Password: nil,
			},
		},
		{
			"simple protocol",
			"protocol=https\n\n",
			GitCrendential{
				Protocol: intoPtr("https"),
				Host:     nil,
				Path:     nil,
				Username: nil,
				Password: nil,
			},
		},
		{
			"simple host",
			"host=example.com\n\n",
			GitCrendential{
				Protocol: nil,
				Host:     intoPtr("example.com"),
				Path:     nil,
				Username: nil,
				Password: nil,
			},
		},
		{
			"simple path",
			"path=foo.git\n\n",
			GitCrendential{
				Protocol: nil,
				Host:     nil,
				Path:     intoPtr("foo.git"),
				Username: nil,
				Password: nil,
			},
		},
		{
			"simple username",
			"username=bob\n\n",
			GitCrendential{
				Protocol: nil,
				Host:     nil,
				Path:     nil,
				Username: intoPtr("bob"),
				Password: nil,
			},
		},
		{
			"simple password",
			"password=secr3t\n\n",
			GitCrendential{
				Protocol: nil,
				Host:     nil,
				Path:     nil,
				Username: nil,
				Password: intoPtr("secr3t"),
			},
		},
		{
			"simple all except url",
			"protocol=https\nhost=example.com\npath=foo.git\nusername=bob\npassword=secr3t\n\n",
			GitCrendential{
				Protocol: intoPtr("https"),
				Host:     intoPtr("example.com"),
				Path:     intoPtr("foo.git"),
				Username: intoPtr("bob"),
				Password: intoPtr("secr3t"),
			},
		},
		{
			"simple url",
			"url=https://bob:secr3t@example.com/foo.git\n\n",
			GitCrendential{
				Protocol: intoPtr("https"),
				Host:     intoPtr("example.com"),
				Path:     intoPtr("foo.git"),
				Username: intoPtr("bob"),
				Password: intoPtr("secr3t"),
			},
		},
		{
			"simple url without password",
			"url=https://bob@example.com/foo.git\n\n",
			GitCrendential{
				Protocol: intoPtr("https"),
				Host:     intoPtr("example.com"),
				Path:     intoPtr("foo.git"),
				Username: intoPtr("bob"),
				Password: nil,
			},
		},
		{
			"simple url without user",
			"url=https://example.com/foo.git\n\n",
			GitCrendential{
				Protocol: intoPtr("https"),
				Host:     intoPtr("example.com"),
				Path:     intoPtr("foo.git"),
				Username: nil,
				Password: nil,
			},
		},
		{
			"url and password",
			"password=secr3t\nurl=https://bob@example.com/foo.git\n\n",
			GitCrendential{
				Protocol: intoPtr("https"),
				Host:     intoPtr("example.com"),
				Path:     intoPtr("foo.git"),
				Username: intoPtr("bob"),
				Password: intoPtr("secr3t"),
			},
		},
		{
			"url and username",
			"username=bob\nurl=https://example.com/foo.git\n\n",
			GitCrendential{
				Protocol: intoPtr("https"),
				Host:     intoPtr("example.com"),
				Path:     intoPtr("foo.git"),
				Username: intoPtr("bob"),
				Password: nil,
			},
		},
		{
			"url and username and password",
			"username=bob\npassword=secr3t\nurl=https://example.com/foo.git\n\n",
			GitCrendential{
				Protocol: intoPtr("https"),
				Host:     intoPtr("example.com"),
				Path:     intoPtr("foo.git"),
				Username: intoPtr("bob"),
				Password: intoPtr("secr3t"),
			},
		},
		{
			"blank line",
			"protocol=https\n\nhost=example.com\npath=foo.git\nusername=bob\npassword=secr3t\n\n",
			GitCrendential{
				Protocol: intoPtr("https"),
				Host:     nil,
				Path:     nil,
				Username: nil,
				Password: nil,
			},
		},
		{
			"unrecognised attribute",
			"protocol=https\nunrecognised=attribute\n\n",
			GitCrendential{
				Protocol: intoPtr("https"),
				Host:     nil,
				Path:     nil,
				Username: nil,
				Password: nil,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			actual, err := util.FromString(test.input)
			if err != nil {
				t.Errorf("unexpecred error: %s", err.Error())
			}

			if !reflect.DeepEqual(&test.expected, actual) {
				t.Fatalf("%#v is expected but got %#v", &test.expected, actual)
			}
		})
	}
}

func TestParseError(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			"protocol and url",
			"protocol=https\nurl=https://bob:secr3t@example.com/foo.git\n\n",
		},
		{
			"host and url",
			"host=example.com\nurl=https://bob:secr3t@example.com/foo.git\n\n",
		},
		{
			"path and url",
			"path=foo.git\nurl=https://bob:secr3t@example.com/foo.git\n\n",
		},
		{
			"username and url",
			"username=bob\nurl=https://bob:secr3t@example.com/foo.git\n\n",
		},
		{
			"password and url",
			"password=secr3t\nurl=https://bob:secr3t@example.com/foo.git\n\n",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			_, err := util.FromString(test.input)
			if err == nil {
				t.Fatalf("error is expected but didn't get")
			}

			t.Logf("got error: %s", err.Error())
		})
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		name     string
		input    GitCrendential
		expected string
	}{
		{
			"empty",
			GitCrendential{
				Protocol: nil,
				Host:     nil,
				Path:     nil,
				Username: nil,
				Password: nil,
			},
			"\n",
		},
		{
			"simple protocol",
			GitCrendential{
				Protocol: intoPtr("https"),
				Host:     nil,
				Path:     nil,
				Username: nil,
				Password: nil,
			},
			"protocol=https\n\n",
		},
		{
			"simple host",
			GitCrendential{
				Protocol: nil,
				Host:     intoPtr("example.com"),
				Path:     nil,
				Username: nil,
				Password: nil,
			},
			"host=example.com\n\n",
		},
		{
			"simple path",
			GitCrendential{
				Protocol: nil,
				Host:     nil,
				Path:     intoPtr("foo.git"),
				Username: nil,
				Password: nil,
			},
			"path=foo.git\n\n",
		},
		{
			"simple username",
			GitCrendential{
				Protocol: nil,
				Host:     nil,
				Path:     nil,
				Username: intoPtr("bob"),
				Password: nil,
			},
			"username=bob\n\n",
		},
		{
			"simple password",
			GitCrendential{
				Protocol: nil,
				Host:     nil,
				Path:     nil,
				Username: nil,
				Password: intoPtr("secr3t"),
			},
			"password=secr3t\n\n",
		},
		{
			"simple all",
			GitCrendential{
				Protocol: intoPtr("https"),
				Host:     intoPtr("example.com"),
				Path:     intoPtr("foo.git"),
				Username: intoPtr("bob"),
				Password: intoPtr("secr3t"),
			},
			"protocol=https\nhost=example.com\npath=foo.git\nusername=bob\npassword=secr3t\n\n",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			actual := util.ToString(&test.input)
			if test.expected != actual {
				t.Fatalf("%#v is expected but got %#v", test.expected, actual)
			}
		})
	}
}
