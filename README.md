<p align="right">
  ⭐ &nbsp;&nbsp;<strong>the project to show your appreciation.</strong> :arrow_upper_right:
</p>

<p align="right">
  <a href="http://godoc.org/github.com/rocketlaunchr/unsafe"><img src="http://godoc.org/github.com/rocketlaunchr/unsafe?status.svg" /></a>
  <a href="https://goreportcard.com/report/github.com/rocketlaunchr/unsafe"><img src="https://goreportcard.com/badge/github.com/rocketlaunchr/unsafe" /></a>
</p>

# Unsafe

This package is for "safely" modifying unexported fields of a struct.

## Safety Considerations

Contrary to popular belief, the `unsafe` package is actually safe to use - provided you know what
you are doing.

1. You need to ensure the `newValue` is the **same type** as the field's type.
2. You need to make sure the field is not being modifed by multiple go-routines concurrently.
3. Read [Exploring ‘unsafe’ Features in Go 1.20: A Hands-On Demo](https://medium.com/@bradford_hamilton/exploring-unsafe-features-in-go-1-20-a-hands-on-demo-7149ba82e6e1) and [Modifying Private Variables of a Struct in Go Using unsafe and reflect](https://medium.com/@darshan.na185/modifying-private-variables-of-a-struct-in-go-using-unsafe-and-reflect-5447b3019a80).
  

## Example

```go
type Example struct {
	e string // unexported string 
}
```

Let's modify the `Example` struct's `e` field.

```go
import "github.com/rocketlaunchr/unsafe"

e := Example{}

// Option 1
ptr := unsafe.SetField[string](&e, unsafe.F("e"))
*(*string)(ptr) = "New Value" // Note: If the field type is string, then ptr must be cast to *string

// Option 2
unsafe.SetField(&e, unsafe.F("e"), "New Value 2")
```

Other useful packages
------------

- [awesome-svelte](https://github.com/rocketlaunchr/awesome-svelte) - Resources for killing react
- [dataframe-go](https://github.com/rocketlaunchr/dataframe-go) - Statistics and data manipulation
- [dbq](https://github.com/rocketlaunchr/dbq) - Zero boilerplate database operations for Go
- [electron-alert](https://github.com/rocketlaunchr/electron-alert) - SweetAlert2 for Electron Applications
- [go-pool](https://github.com/rocketlaunchr/go-pool) - A Generic pool
- [igo](https://github.com/rocketlaunchr/igo) - A Go transpiler with cool new syntax such as fordefer (defer for for-loops)
- [mysql-go](https://github.com/rocketlaunchr/mysql-go) - Properly cancel slow MySQL queries
- [react](https://github.com/rocketlaunchr/react) - Build front end applications using Go
- [remember-go](https://github.com/rocketlaunchr/remember-go) - Cache slow database queries
- [testing-go](https://github.com/rocketlaunchr/testing-go) - Testing framework for unit testing