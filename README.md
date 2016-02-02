# beter

Don't get lost going up the stack! Works with custom errors, and doesn't require too much modification around your existing code. Just wrap all returned/passed errors with `b.E()` and assert the received error to `*b.B` when you would like to access details about the stack that returned the first error.

## Examples

```Go
package main

import (
	"fmt"
	"os"

	"github.com/c3mb0/beter"
)

func check(err error) {
	if err != nil {
		details := err.(*b.B)
		fmt.Printf("%s\n%s:%d\n", details.Err, details.Fn, details.Line)
		os.Exit(1)
	}
}

func main() {
	f, err := os.Open(`C:\sshd.log`)
	check(b.E(err))
	_, err = f.Seek(-1, 0)
	check(b.E(err)) // line 22
}
```

Output:

```
seek C:\sshd.log: An attempt was made to move the file pointer before the beginning of the file.
main.main:22
exit status 1
```

A more detailed and alternative use case can be as such:

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
		return nil, b.E(&custError{err: "login not successful"}) // line 43
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
