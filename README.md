[![jacobconley](https://circleci.com/gh/jacobconley/habitat.svg?style=shield)](https://app.circleci.com/pipelines/github/jacobconley/habitat)
[![GitHub tag](https://img.shields.io/github/tag/jacobconley/habitat.svg)](https://github.com/jacobconley/habitat/tags)

Habitat
======

Habitat is a declarative, open-ended, holistic web framework for the Go language.  It wraps the [Gorilla](https://github.com/gorilla) toolkit and provides support functionality to server-side Go code, 
    and also integrates with the Node.js ecosystem to facilitate seamless development and deployment of client-side code. 


More to come!

Features
----

 - Integrated with many parts of the [Gorilla](https://github.com/gorilla) toolkit, including its powerful routing and middleware configurations
 - Declarative HTML metadata and dependency configuration, with powerful layout and templating functionality
 - Auxiliary toolchains:  Automatic SCSS and Webpack builds 
 - Customizable error handling, with configurable renderers for HTML and JSON types 
 - [Zerolog](https://github.com/rs/zerolog) integration
 - Asset fingerprinting



Development and testing
=======

Unit tests
-----

Everything that's not in the `./test-e2e` directory is a unit test, so it does not have any dependencies.  The `build/test-unit` script runs all of the unit tests (packages that aren't `e2e`).
Since `go test` doesn't seem to have an easy way to exclude file patterns, we have to maintain a list of packages in that script, or eventually replace it with a solution that invokes `go list ./...`. 

End-to-end Tests
-----

For the web framework itself, the end-to-end test environment is located at `test-fixtures/userland`, which supports the tests located in `test-e2e`.  There's space here to support test environments for middleware / plugins, specific use cases, or whatever else we could need that's beyond the scope of unit tests.  

The bash script `build/test-e2e` runs all end-to-end tests; it's the entry point for the the end-to-end testing in the CI environment and should run fine on any machine that has bash and all of the test environment dependencies listed below.  I'll try my best to keep this up to date but you can always inspect the script for a detailed up-to-date implementation.  In general, successfully running the script once which should perform any required setup, and afterwards any test should complete successfully as long as the test server at `test-fixtures/userland/main.go` is running. 

### Test environment dependencies
 - **Node + yarn**, and the package dependencies (just run `yarn install` in `test-fixtures/userland`)