package opax

import (
	"context"
	"github.com/w6d-io/x/errorx"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Conn is the struct variable for connect to a opa server
type Conn struct {
	// Address is the address to the kratos micro service
	Protocol string `json:"protocol" mapstructure:"protocol"`
	// Address is the address to the kratos micro service
	Address string `json:"address" mapstructure:"address"`
	// Port is the port of the uri to the kratos micro service
	Port string `json:"port" mapstructure:"port"`
	// Verbose state to call to the kratos micro service
	Verbose string `default:"false" json:"verbose" mapstructure:"verbose"`
}

// Helper
//
// GetAuthorizationFromHttp is used to check if the request to query is authorized or unauthorized
// and also return opa decision
// if params query is not set, return a nil query with StatusBadRequest and error
// if OPA is unreachable or an other issues, return nil deccision with statusCode of the call and error-go
//
// GetAuthorizationFromGRPCCtx is used to check if the request to query is authorized or unauthorized
// and also return opa decision
// It checks if opa configuration on context is present
// if params query is not set, return a nil query with StatusBadRequest and error
// if OPA is unreachable or an other issues, return nil deccision with statusCode of the call and error-go
type Helper interface {
	GetAuthorizationFromHttp(context.Context, interface{}) (string, error)

	GetAuthorizationFromGRPCCtx(context.Context) (string, error)
}

var (
	_ Helper = &auth{}
)

type auth struct {
	Conn
}

var (
	Opax Helper
)

// callOpaServer concat and format the svc and port from Conn variable
func (k auth) callOpaServer(uri string, body string) (string, int, string, error) {

		address := k.Protocol + k.Address + k.Port + uri

		req, err := http.NewRequest("POST", address, strings.NewReader(body))
		if err != nil {
			return "", 0, "", errorx.Wrap(err, "decode address failed")
		}
		req.Header.Add("content-type","application/json")
		req.Header.Add("cache-control", "no-cache")

		// TODO: check err
		client := &http.Client{}

		res, err := client.Do(req)
		if err != nil {
			return "", 0, "", errorx.Wrap(err, "decode address failed with address " + address)
		}
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", 0, "", errorx.Wrap(err, "decode address failed")
		}

		return res.Status, res.StatusCode, string(data), err
	}

// getVerboseState return the verbose state from Conn variable
func (k Conn) getVerboseState() bool {

	// checking if null string
	if k.Verbose == "" {
		k.Verbose = "false"
	}

	state, err := strconv.ParseBool(k.Verbose)
	if err != nil {
		panic(err)
	}

	return state
}

// SetOpaxDetails ip or uri and set port with verbose state. Default port is nil and default verbose is false.
// In production mode is not necessary to set a verbose state in the ci configuration file
//
// TEST UNITARY Opax
//
// Before Test Opax prepare environment with mock OPA binary with command lines :
// << make opa >> and after << make run >>
//
// For stop opa server run command line : << make stop >>
//
// Run test with command line : << make test >>
func SetOpaxDetails(https bool, address string, verbose bool, port ...int64) {
	var p string
	var h string

	if len(port) > 0 {
		p = ":" + strconv.Itoa(int(port[0]))
	}

	if !https {
		h = "http://"
	} else {
		h = "https://"
	}

	Opax = &auth{Conn{Protocol:h, Address: address, Port: p, Verbose: strconv.FormatBool(verbose)}}
}