# lambda-invoke

Small CLI package to allow you to invoke your Go AWS lambda locally.

## Installing

To use the package
```
go get -u github.com/cshep4/kripto/shared/go/local
```
To use the CLI
```
go install github.com/cshep4/kripto/shared/go/local/cli
```

## Example usage

Run an example lambda function lambdafunction.go on port 8001

```
_LAMBDA_SERVER_PORT=8001 go run ./lambdafunction.go
```

Then use this library in tests or wherever you need it, by calling 

```go
response, err := local.Invoke(local.Input{
    Port:    8001,
    Payload: "payload",
})
```
or
```
lambda-invoke -p=8001 -d=payload
```

Note that `Payload` can be any structure that can be encoded by the `encoding/json` package.
Your lambda function will need to use this structure in its type signature. For use of the CLI a JSON
string can be passed using the `-d` flag.
