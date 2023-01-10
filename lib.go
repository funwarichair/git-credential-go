package gitcredential

import (
	"bufio"
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

func (g *GitCrendential) Write(writer io.Writer) (err error) {
	// Write must return a non-nil error if it returns n < len(p).
	// So conversely the error being nil implies n == len(p).

	if g.Protocol != nil {
		_, err = fmt.Fprintf(writer, "protocol=%s\n", *g.Protocol)
		if err != nil {
			return
		}
	}

	if g.Host != nil {
		_, err = fmt.Fprintf(writer, "host=%s\n", *g.Host)
		if err != nil {
			return
		}
	}

	if g.Path != nil {
		_, err = fmt.Fprintf(writer, "path=%s\n", *g.Path)
		if err != nil {
			return
		}
	}

	if g.Username != nil {
		_, err = fmt.Fprintf(writer, "username=%s\n", *g.Username)
		if err != nil {
			return
		}
	}

	if g.Password != nil {
		_, err = fmt.Fprintf(writer, "password=%s\n", *g.Password)
		if err != nil {
			return
		}
	}

	_, err = fmt.Fprintf(writer, "\n")
	if err != nil {
		return
	}

	return
}
