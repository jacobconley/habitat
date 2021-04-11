Habitat
======

Habitat is more-or-less a web framework for the Go programming language, 
though that feels like an overstatement.
Mainly, it is a simple configuration of the [Buffalo](https://github.com/gobuffalo) and [Gorilla](https://github.com/gorilla) ecosystem, 
along with some other crucial packages for a modern web development environment (who will all be properly thanked in a SHOULDERS.md soon),
and a few original features to tie the whole thing together.  

More to come!


Building and testing
=======

Running End-to-end Tests
-----

For the web framework itself, the end-to-end test environment is located at `test-fixtures/userland`, which supports the tests located in `test/e2e`.  There's space here to support test environments for middleware / plugins, specific use cases, or whatever else we could need that's beyond the scope of unit tests.  

The bash script `build/test-e2e` runs all end-to-end tests; it's the entry point for the the end-to-end testing in the CI environment and should run fine on any machine that has bash and all of the test environment dependencies listed below.  I'll try my best to keep this up to date but you can always inspect the script for a detailed up-to-date implementation.  In general, successfully running the script once which should perform any required setup, and afterwards any test should complete successfully as long as the test server at `test-fixtures/userland/main.go` is running. 

### Test environment dependencies
 - **Node + yarn**, and the package dependencies (just run `yarn install` in `test-fixtures/userland`)
 - A running **PostgreSQL** instance, with the appropriate user as is configured in `test-fixtures/userland/habitat.toml`