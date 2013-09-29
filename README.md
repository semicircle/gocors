gocors
======

A tiny toolkit for Go CORS support.

### Explaination:

CORS means "Cross-Origin Resource Sharing", which your can find more info here[http://enable-cors.org/]

### Usage:

```go 
	cors := gocors.New()
	http.Handle("/", cors.Handler(x))
```

For any `http.Handler` x, just replace it with `cors.Handler(x)`, then the handler will have the ability of handling CORS request.

### Installation:

Install with `go get` command:

```
go get github.com/semicircle/gocors
```
The Go distribution is the only dependency.

### More:

All the features described in this article[www.html5rocks.com/en/tutorials/cors/] and a TL;DR version[http://www.html5rocks.com/static/images/cors_server_flowchart.png]

So, all the 'Access-Control' parameters can be set by this:

```go
	c := gocors.New()

	// for the 'Access-Control-xxx' headers.
	c.SetAllowOrigin("*")

	c.SetAllowMethods([]string{"PUT", "POST", "DELETE"})

	//important: must contain a 'origin', for I found Chrome's request contains this header.
	c.SetAllowHeaders([]string{"Custom-Headers", "origin"})

	...
```

Important: SetAllowHeaders () must be called with a 'origin' header.

### How I test: (The state of the project)

In fact, this code is unstable now, I have only tested with www.test-cors.org.

It works for me.

### Security:

This piece of code designed to work with a "nice" browser. A 'nice' browser means it follows the flow-chart of the CORS in the communication.So, if any request that didn't obey the Access-Control rules, they will NOT be denied. 

In other words, Gocors just tell a 'nice' browser : you can make a Cross Domain request, and is NOT responsible of any security issue.

This is because checking if the request is valid is meaningless, and any program, spider or so, can make a "non Cross Domain" request directly.


### Any pull requests are welcome

This is a quick and dirty implement. A 'just work' implement. 

So, any pull requests are welcome.

