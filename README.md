# errortrace

**errortrace** is a Go package for enhanced error tracing. It wraps standard Go errors, providing context like error location, metadata, and detailed error information to make debugging easier.

## Motivation
From my experience working on various projects, handling errors required manually adding context, such as the package and function, which was tedious and error-prone. ErrorTrace was created to automate this process, ensuring that context is always included with each error. Additionally, the library provides a way to define user-friendly error messages, like titles and details, which can be used in the presentation layer to offer clear, helpful responses to the client.

## Features

- Add error metadata and status codes
- Track error origins with `runtime.Caller`
- Easy to wrap and chain errors

## Installation

```bash
go get github.com/techforge-lat/errortrace
```

## Usage Example

```go
import "github.com/techforge-lat/errortrace"

err := someFunc()
trace := errortrace.Wrap(err).SetTitle("Database Error").SetDetail("Failed to connect")
fmt.Println(trace.Error())
```

## Methods Overview

- `Wrap(err error)`: Wraps an error with tracing.
- `SetTitle(title string)`: Adds a title to the error.
- `SetDetail(detail string)`: Adds detailed info about the error.
- `AddMetadata(key string, value any)`: Adds metadata.
- `Error()`: Outputs error details and trace.

## License

This project is licensed under the MIT License.

---

Let me know if you'd like more details or changes!
