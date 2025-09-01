# Change Log

**ATTN**: This project uses [semantic versioning](http://semver.org/).

## [3.1.1] - [2024-06-04]

### Fixed

- Writing an HTTP 1xx status codes no longer results in a 200 being sent in
  addition given the expectation is that a follow-up status code will be written
  later. The exception is `101 Switching Protocols` since this is terminal. This
  matches `net/http` behavior.

## [3.1.0] - [2023-10-07]

### Fixed

- `WriteHeader` can again be used in the `ResponseWriter.Before` callback to
  set HTTP status codes

### Added

- Satisfy `http.ResponseController` (Go 1.20+)

## [3.0.0] - [2022-09-18]

### Fixed

- Replace multiple slashes at the beginning of a path with a single one to avoid
  open redirects
- Avoid updating `ResponseWriter.Status()` if the status has already been
  written when `ResponseWriter.WriteHeader()` is called twice

### Changed

- `ResponseWriter` now only implements `http` interfaces (`Flusher`, `Hijacker`,
  `CloseNotifier`) if the wrapped `http.ResponseWriter` does. This is a breaking
  change to make `ResponseWriter`'s interface support more accurate

### Added

- `ResponseWriter` now implements `io.ReaderFrom` to more efficiently send
  responses via `io.Copy`

## [2.0.2] - 2020-07-17

### Fixed

- Fixed Go module name for v2

## [2.0.1] - 2020-05-25

### Fixed

- Recovery middleware now checks that `Formatter` is not `nil` before calling

## [2.0.0] - 2020-05-25

### Changed

- `Recovery.PrintStack`, when false, now also supresses the panic message in
  addition to supressing the stack trace

### Fixed

- `Negroni.With()` now copies handlers to avoid mutating the original `Negroni`
  instance if `Use` is called on the new `Negroni` instance

### Added

- `Recovery.LogStack` was added to control whether the stacktrace is logged for
  panics

### Changed

## [1.0.0] - 2018-09-01

### Fixed
- `Logger` middleware now correctly handles paths containing a `%` instead of trying to treat it as a format specifier

## [0.3.0] - 2017-11-11
### Added
- `With()` helper for building a new `Negroni` struct chaining handlers from
  existing `Negroni` structs
- Format log output in `Logger` middleware via a configurable `text/template`
  string injectable via `.SetFormat`. Added `LoggerDefaultFormat` and
  `LoggerDefaultDateFormat` to configure the default template and date format
  used by the `Logger` middleware.
- Support for HTTP/2 pusher support via `http.Pusher` interface for Go 1.8+.
- `WrapFunc` to convert `http.HandlerFunc` into a `negroni.Handler`
- `Formatter` field added to `Recovery` middleware to allow configuring how
  `panic`s are output. Default of `TextFormatter` (how it was output in
  `0.2.0`) used. `HTMLPanicFormatter` also added to allow easy outputing of
  `panic`s as HTML.

### Fixed
- `Written()` correct returns `false` if no response header has been written
- Only implement `http.CloseNotifier` with the `negroni.ResponseWriter` if the
  underlying `http.ResponseWriter` implements it (previously would always
  implement it and panic if the underlying `http.ResponseWriter` did not.

### Changed
- Set default status to `0` in the case that no handler writes status -- was
  previously `200` (in 0.2.0, before that it was `0` so this reestablishes that
  behavior)
- Catch `panic`s thrown by callbacks provided to the `Recovery` handler
- Recovery middleware will set `text/plain` content-type if none is set
- `ALogger` interface to allow custom logger outputs to be used with the
  `Logger` middleware. Changes embeded field in `negroni.Logger` from `Logger`
  to `ALogger`.
- Default `Logger` middleware output changed to be more structure and verbose
  (also now configurable, see `Added`)
- Automatically bind to port specified in `$PORT` in `.Run()` if an address is
  not passed in. Fall back to binding to `:8080` if no address specified
  (configuable via `DefaultAddress`).
- `PanicHandlerFunc` added to `Recovery` middleware to enhance custom handling
  of `panic`s by providing additional information to the handler including the
  stack and the `http.Request`. `Recovery.ErrorHandlerFunc` was also added, but
  deprecated in favor of the new `PanicHandlerFunc`.

## [0.2.0] - 2016-05-10
### Added
- Support for variadic handlers in `New()`
- Added `Negroni.Handlers()` to fetch all of the handlers for a given chain
- Allowed size in `Recovery` handler was bumped to 8k
- `Negroni.UseFunc` to push another handler onto the chain

### Changed
- Set the status before calling `beforeFuncs` so the information is available to them
- Set default status to `200` in the case that no handler writes status -- was previously `0`
- Panic if `nil` handler is given to `negroni.Use`

## 0.1.0 - 2013-07-22
### Added
- Initial implementation.

[Unreleased]: https://github.com/urfave/negroni/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/urfave/negroni/compare/v0.1.0...v0.2.0
