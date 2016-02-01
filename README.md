## beter

Don't get lost going up the stack! Works with custom errors, and doesn't require too much modification around your existing code.

# Example

```Go
package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	"github.com/c3mb0/beter"
)

type custError struct {
	err string
}

func (e *custError) Error() string {
	return e.err
}

type inputs struct {
	value1 string
	value2 string
}

func (i *inputs) login() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, b.E(err)
	}
	client := &http.Client{Jar: jar}

	form := url.Values{}
	form.Add("username", i.value1)
	form.Add("password", i.value2)

	login, err := client.PostForm(os.Args[3], form)
	if err != nil {
		return nil, b.E(err) // line 39
	}
	defer login.Body.Close()
	if login.StatusCode != 200 {
		return nil, b.E(&custError{"login not successful"}) // line 43
	}

	return client, nil
}

func main() {
	if err := run(); err != nil {
		details := err.(*b.B)
		fmt.Printf("%s\n%s:%d\n", details.Err, details.Fn, details.Line)
		os.Exit(1)
	}
}

func run() error {
	values := &inputs{value1: os.Args[1], value2: os.Args[2]}
	_, err := values.login()
	if err != nil {
		return b.E(err)
	}
	return nil
}
```

Given `abc efg "https://dummy"`:

```
Post https://dummy: dial tcp: lookup dummy: getaddrinfow: No such host is known.
main.(*inputs).login:39
exit status 1
```

Given `abc efg "https://google.com"`:

```
login not successful
main.(*inputs).login:43
exit status 1
```
