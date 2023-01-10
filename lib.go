package gitcredential

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
)

type GitCrendential struct {
	Protocol *string
	Host     *string
	Path     *string
	Username *string
	Password *string
}

func Parse(reader io.Reader) (cred *GitCrendential, err error) {
	cred = &GitCrendential{}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break
		}

		if strings.HasPrefix(line, "protocol=") {
			protocol := strings.TrimPrefix(line, "protocol=")
			if cred.Protocol == nil {
				cred.Protocol = &protocol
			} else {
				return nil, fmt.Errorf("conflict protocol: %s and %s", *cred.Protocol, protocol)
			}
			continue
		}

		if strings.HasPrefix(line, "host=") {
			host := strings.TrimPrefix(line, "host=")
			if cred.Host == nil {
				cred.Host = &host
			} else {
				return nil, fmt.Errorf("conflict host: %s and %s", *cred.Host, host)
			}
			continue
		}

		if strings.HasPrefix(line, "path=") {
			path := strings.TrimPrefix(line, "path=")
			if cred.Path == nil {
				cred.Path = &path
			} else {
				return nil, fmt.Errorf("conflict path: %s and %s", *cred.Path, path)
			}
			continue
		}

		if strings.HasPrefix(line, "username=") {
			username := strings.TrimPrefix(line, "username=")
			if cred.Username == nil {
				cred.Username = &username
			} else {
				return nil, fmt.Errorf("conflict username: %s and %s", *cred.Username, username)
			}
			continue
		}

		if strings.HasPrefix(line, "password=") {
			password := strings.TrimPrefix(line, "password=")
			if cred.Password == nil {
				cred.Password = &password
			} else {
				return nil, fmt.Errorf("conflict password: %s and %s", *cred.Password, password)
			}
			continue
		}

		if strings.HasPrefix(line, "url=") {
			url, err := url.Parse(strings.TrimPrefix(line, "url="))
			if err != nil {
				return nil, fmt.Errorf("parse url: %w", err)
			}

			if cred.Protocol == nil {
				cred.Protocol = &url.Scheme
			} else {
				return nil, fmt.Errorf("conflict protocol: %s and %s", *cred.Protocol, url.Scheme)
			}

			if cred.Host == nil {
				cred.Host = &url.Host
			} else {
				return nil, fmt.Errorf("conflict host: %s and %s", *cred.Host, url.Host)
			}

			// I couldn't find an explicit description, but the example below eliminates the leading slash in the path.
			// https://git-scm.com/docs/git-credential#_typical_use_of_git_credential
			var path string
			if strings.HasPrefix(url.Path, "/") {
				path = url.Path[len("/"):len(url.Path)]
			} else {
				path = url.Path
			}
			if cred.Path == nil {
				cred.Path = &path
			} else {
				return nil, fmt.Errorf("conflict path: %s and %s", *cred.Path, path)
			}

			if url.User != nil {
				username := url.User.Username()
				if cred.Username == nil {
					cred.Username = &username
				} else {
					return nil, fmt.Errorf("conflict username: %s and %s", *cred.Username, username)
				}
			}

			if password, ok := url.User.Password(); ok {
				if cred.Password == nil {
					cred.Password = &password
				} else {
					return nil, fmt.Errorf("conflict password: %s and %s", *cred.Password, password)
				}
			}

			continue
		}

		// A line must contain "="
		if !strings.Contains(line, "=") {
			return nil, errors.New("invalid line")
		}

		// Unrecognised attributes are silently discarded.
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	return cred, nil
}

func (g GitCrendential) String() string {
	var buf bytes.Buffer

	if g.Protocol != nil {
		buf.WriteString(fmt.Sprintf("protocol=%s\n", *g.Protocol))
	}

	if g.Host != nil {
		buf.WriteString(fmt.Sprintf("host=%s\n", *g.Host))
	}

	if g.Path != nil {
		buf.WriteString(fmt.Sprintf("path=%s\n", *g.Path))
	}

	if g.Username != nil {
		buf.WriteString(fmt.Sprintf("username=%s\n", *g.Username))
	}

	if g.Password != nil {
		buf.WriteString(fmt.Sprintf("password=%s\n", *g.Password))
	}

	buf.WriteString("\n")

	return buf.String()
}
