package supertest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// MainApp matches signature for main()
// Typical usage is to pass a wrapper function that calls main()
type MainApp func()

// Resp stores the response information in which we are interested
type Resp struct {
	StatusCode string
	Headers    http.Header
	Body       []byte
}

// AppRunner holds onto all errors for a single GET or POST
// Because of this, you cannot concurrently test multiple endpoints
// without spinning up multiple instances of AppRunner.
type AppRunner struct {
	Addr        string
	SetupBuffer time.Duration
	MainWrapper MainApp
	Errors      map[string]string
	Resp
}

// NewAppRunner sets a default SetupBuffer. This gives the main app
// enough time to start up before we need to test against it.
func NewAppRunner(addr string, main MainApp) *AppRunner {
	return &AppRunner{
		Addr:        addr,
		SetupBuffer: 100 * time.Millisecond,
		MainWrapper: main, // represents the main(). This is why everything works.
	}
}

// Start kicks off the applications main() method
func (a *AppRunner) Start() {
	go a.MainWrapper()
	time.Sleep(a.SetupBuffer)
}

// Get expects to have at least one chained method of Expects*. The chain ends with End().
func (a *AppRunner) Get(route string) *AppRunner {
	// clear out the any pre-existing errors
	a.Errors = make(map[string]string)

	url := "http://" + a.Addr + route
	res, err := http.Get(url)
	if err != nil {
		a.Errors["Error GET "+url] = err.Error()
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		a.Errors["Error Reading Body"] = err.Error()
	}
	response := Resp{}
	response.Body = body
	response.StatusCode = res.Status
	response.Headers = res.Header

	a.Resp = response

	return a
}

func (a *AppRunner) ExpectHeader(k, v string) *AppRunner {

	actual := a.Headers.Get(k)

	if actual == "" {
		a.Errors["Header Missing"] = k + " not found"
	} else if v != actual {
		a.Errors["Header Value Error: "+k] = actual + " is not " + v
	}
	return a
}

// ExpectStatusCode should be renamed ExpectStatus (as we return more than just the code)
func (a *AppRunner) ExpectStatusCode(code string) *AppRunner {
	if a.StatusCode != code {
		a.Errors["Status Code Error"] = a.StatusCode + " is not " + code
	}
	return a
}

func (a *AppRunner) ExpectContent(content []byte) *AppRunner {
	if string(a.Body) != string(content) {
		a.Errors["Body Content Error"] = string(a.Body) + "\nis not\n" + string(content)
	}
	return a
}

func (a *AppRunner) End() *AppRunnerError {
	if len(a.Errors) == 0 {
		return nil
	}

	err := &AppRunnerError{}

	for k, v := range a.Errors {
		err.S += "\n" + k + " -> \n" + v + "\n"
	}
	return err
}

type AppRunnerError struct {
	S string
}

func (e *AppRunnerError) Error() string {
	return e.S
}
