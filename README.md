# opax

## Constants

```golang
const (
    // OpaDataName where is stored the configuration query
    OpaDataName = "opa"
)
```

## Functions

### func [SetOpaxDetails](/opax.go#L111)

`func SetOpaxDetails(https bool, address string, verbose bool, port ...int64)`

SetOpaxDetails ip or uri and set port with verbose state. Default port is nil and default verbose is false.
In production mode is not necessary to set a verbose state in the ci configuration file

HOW RUN UNITS TEST TO Opax ?

Before Test:
- prepare environment to mock OPA binary with command lines :
<< make opa >> and after << make run >>

- for stop opa server run command line : << make stop >>

- run test with command line : << make test >>

## Types

### type [Conn](/opax.go#L13)

`type Conn struct { ... }`

Conn is the struct variable for connect to a opa server

### type [Helper](/opax.go#L36)

`type Helper interface { ... }`

Helper

GetAuthorizationFromHttp is used to check if the request to query is authorized or unauthorized
and also return opa decision
if params query is not set, return a nil query with StatusBadRequest and error
if OPA is unreachable or an other issues, return nil deccision with statusCode of the call and error-go

GetAuthorizationFromGRPCCtx is used to check if the request to query is authorized or unauthorized
and also return opa decision
It checks if opa configuration on context is present
if params query is not set, return a nil query with StatusBadRequest and error
if OPA is unreachable or an other issues, return nil deccision with statusCode of the call and error-go

#### Variables

```golang
var (
    Opax Helper
)
```

### type [Query](/authz.go#L28)

`type Query struct { ... }`
