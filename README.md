go-supertest
============

Inspired by supertest from node, this is my GopherCon hackday project. This was kicked off by a coworker who was frustrated with how to do an integration test for his HTTP Go project. He went back to Node. I built this as a proof of concept. I don't think I would use this in a production environment. YMMV. 


Usage
=====
Let's say you set up a classic Martini sample application with the text, "Hello, World!". You continue to set up unit tests, just like normal. For integration tests, you can run ```go test -tags "integration"``` and the following test will execute:

    // +build integration

    package main

    import (
    	"github.com/sethgrid/go-supertest"

    	"testing"
    )

    // to be passed into supertest
    func wrapper() {
    	go main()
    }

    func TestGetContent(t *testing.T) {

    	app := supertest.NewAppRunner("localhost:3000", wrapper)

        // start your app
    	app.Start()

    	// you can method chain and check headers, content, and status
    	err := app.Get("/").
    		ExpectStatusCode("200 OK").
    		ExpectHeader("Content-Type", "text/plain; charset=utf-8").
    		ExpectHeader("Content-Length", "13").
    		ExpectContent([]byte("Hello, World!")).
    		End()
    	if err != nil {
    		t.Error("Error Getting Root", err)
    	}

	    // or just check something indeed came up.
	    err = app.Get("/foo").ExpectStatusCode("200 OK").ExpectContent([]byte("bar")).End()
	    if err != nil {
		    t.Error("Error getting foo", err)
	    }
	    // or verify that it does not exist.
	    err = app.Get("/nonexistant").ExpectStatusCode("404 Not Found").End()
	    if err != nil {
		    t.Error("Error getting foo", err)
	    }
    }
