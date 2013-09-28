gocors
======

A tiny toolkit for Go CORS support.

### Explaination:

CORS means "Cross-Origin Resource Sharing", which your can find more info here[http://enable-cors.org/]

### Usage:

```go 
	c := gocors.New()
	http.Handle("/", c.Handler(x))
}
```

For any `http.Handler` x, just replace it with `c.Handler(x)`, then the handler will have the ability of handling CORS request.

### Installation:

Install with `go get` command:

```
go get github.com/semicircle/gocors
```
The Go distribution is the only dependency.:wq

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
