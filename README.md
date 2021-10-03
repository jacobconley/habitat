[![jacobconley](https://circleci.com/gh/jacobconley/habitat.svg?style=shield)](https://app.circleci.com/pipelines/github/jacobconley/habitat)
[![GitHub tag](https://img.shields.io/github/tag/jacobconley/habitat.svg)](https://github.com/jacobconley/habitat/tags)

Habitat
======

Habitat is a declarative, open-ended, holistic web framework written primarily in the Go language.  It wraps the [Gorilla](https://github.com/gorilla) toolkit and provides support functionality to server-side Go code, 
    and also integrates with the Node.js ecosystem to facilitate seamless development and deployment of client-side code. 

More to come!


Philosophy
----

Habitat stands in contrast to Rails and other Rails-inspired frameworks in Go.  I love Rails as much as the next guy, but Habitat exists to serve a slightly different purpose, guided by the below principles: 

 - **Declarative**: While the [convention over configuration](https://en.wikipedia.org/wiki/Convention_over_configuration) paradigm can be convenient, it comes with trade-offs in that it begets untraceable code and can sometimes be restrictive.  Habitat does exactly what you tell it to do, only when you tell it to do so.  And unlike some other frameworks, which are completely encapsulated behind their server mechanism, Habitat leverages Go's powerful standard [`http`](https://pkg.go.dev/net/http) implementations to give the programmer complete control over the entire server lifecycle, not just the lifecycle of a single request boxed in to that request's context.
 - **Open-ended**: Habitat aims to be reasonably unopinionated.  I'm not a fan of Rails-style separation-of-concerns which splits relevant code across many different directories.  While I think the model-view-controller framework is good as a guiding principle in general, I built Habitat to be free of restrictive organizational structures that might not make sense in every scenario.  I believe that one of Go's biggest strengths as a language is in the extensibility of its ecosystem, and that the flexibility it provides while still being robust makes it well worth the trade-off of requiring more explicit declarative configuration. 
 - **Holistic**:  One thing that I think Rails does well is integrating with the Javascript ecosystem.  I have _very many_ criticisms of Javascript and the Node ecosystem that are common in the webdev community as a whole.  Nonetheless, it reigns as the dominant platform for front-end web development, at least for the time being.  Therefore, Habitat will support it and aim to provide a seamless, powerful, integrated development experience to whatever extent it does not conflict with the principles above.  With Habitat I aim to take this integration a step further than other server frameworks, leveraging modern client-side techniques and frameworks to unify the front-end and back-end development experience to create a platform built for the modern era. 


Features
----

 - Integrated with the [Gorilla](https://github.com/gorilla) toolkit, including its powerful routing and middleware configurations
 - Declarative HTML metadata and dependency configuration, with powerful layout and templating functionality
 - Auxiliary toolchains:  Automatic SCSS and Webpack builds 
 - Easy rendering for HTML, JSON, or raw text, with customizable error handling
 - Clean parameter handling, with a JSON-style unmarshalling function for URL parameters that also handles string parsing 
 - Thorough logging with [Zerolog](https://github.com/rs/zerolog)
 - Asset fingerprinting



Coming soon
----

 - Deep React.js integration
 - Task scheduling
 - A whole world of possibilities!  Feel free to file a [proposal issue](https://github.com/jacobconley/habitat/issues/new?labels=proposal) with any suggestions.  


Development and testing
=======

Unit tests
-----

Every test in the source directories, rather than in `./test-e2e`, is a unit test without any special dependencies.  The `build/test-unit` script runs all of the unit tests (packages that aren't `e2e`).
Since `go test` doesn't seem to have an easy way to exclude file patterns, we have to maintain a list of packages in that script, or eventually replace it with a solution that invokes `go list ./...`. 

End-to-end Tests
-----

For the web framework itself, the end-to-end test environment is located at `test-fixtures/userland`, which supports the tests located in `test-e2e`.  There's space here to support test environments for middleware / plugins, specific use cases, or whatever else we could need that's beyond the scope of unit tests.  

The bash script `build/test-e2e` runs all end-to-end tests; it's the entry point for the the end-to-end testing in the CI environment and should run fine on any machine that has bash and all of the test environment dependencies listed below.  I'll try my best to keep this up to date but you can always inspect the script for a detailed up-to-date implementation.  In general, successfully running the script once which should perform any required setup, and afterwards any test should complete successfully as long as the test server at `test-fixtures/userland/main.go` is running. 

### Test environment dependencies
 - **Node + yarn**, and the package dependencies (just run `yarn install` in `test-fixtures/userland`)