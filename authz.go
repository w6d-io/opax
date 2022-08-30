package opax

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	"google.golang.org/grpc/metadata"

	"github.com/w6d-io/x/errorx"
	"github.com/w6d-io/x/logx"
)

const (
	// OpaDataName where is stored the configuration query
	OpaDataName = "opa"
)

var (
	errQueryNotFound       = errorx.New(nil, OpaDataName+" params request not found")
	errQueryPathNotFound   = errorx.New(nil, OpaDataName+" Path param request not found")
	errQueryInputNotFound  = errorx.New(nil, OpaDataName+" Input param request not found")
	errNoMDFromCtx         = errorx.New(nil, "cannot get metadata from context")
	errDecisionOpaFromCurl = errorx.New(nil, OpaDataName+" Decision from curl not successful")
)

type Query struct {
	Path  string                 `json:"path"`
	Input map[string]interface{} `json:"input"`
}

// GetAuthorizationFromHttp is used to check if the request to query is authorized or unauthorized
// and also return opa decision
// if params query is not set, return a nil query with StatusBadRequest and error
// if OPA is unreachable or an other issues, return nil deccision with statusCode of the call and error-go
func (a auth) GetAuthorizationFromHttp(ctx context.Context, params interface{}) (string, error) {
	log := logx.WithName(ctx, "GetAuthorizationFromHttp")
	var query Query

	m, _ := params.(map[string]interface{})
	b, err := json.Marshal(m)

	err = json.Unmarshal(b, &query)
	if err != nil {
		log.Error(err, "get all params from query failed")
		return "", errorx.NewHTTP(err, http.StatusBadRequest, "fail on params request")
	}

	if len(query.Path) < 1 {
		log.Error(errQueryPathNotFound, "get path param from query failed")
		return "", errorx.NewHTTP(errNoMDFromCtx, http.StatusBadRequest, "fail to get param Path to request")
	}

	if len(query.Input) < 1 {
		log.Error(errQueryInputNotFound, "get input param from query failed")
		return "", errorx.NewHTTP(errNoMDFromCtx, http.StatusBadRequest, "fail to get param Input to request")
	}

	log.V(2).Info("query request FromHttp", "params", query)

	return a.do(ctx, query)
}

// GetAuthorizationFromGRPCCtx is used to check if the request to query is authorized or unauthorized
// and also return opa decision
// It checks if opa configuration on context is present
// if params query is not set, return a nil query with StatusBadRequest and error
// if OPA is unreachable or an other issues, return nil deccision with statusCode of the call and error-go
func (a auth) GetAuthorizationFromGRPCCtx(ctx context.Context) (string, error) {
	log := logx.WithName(ctx, "GetAuthorizationFromGRPCCtx")

	//get metadata from ctx
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		log.Error(errNoMDFromCtx, "metadata boolean from metadata.FromIncomingContext(ctx) = %v", ok)
		return "", errorx.NewHTTP(errNoMDFromCtx, http.StatusNotFound, "fail to get metadata")
	}

	// check if opa configuration is present on our metadata
	if _, ok := md["opa"]; !ok {
		log.Error(errQueryNotFound, "metadata doesn't exist", "opa", OpaDataName)
		return "", errorx.NewHTTP(errQueryNotFound, http.StatusNotFound, "bad metadata")
	}

	// check if we have more than zero value for this key cause MD is map[string][]string
	if len(md["opa"]) == 0 || len(md["opa"][0]) == 0 {
		log.Error(errQueryNotFound, "metadata exist but no value params exist")
		return "", errorx.NewHTTP(errQueryNotFound, http.StatusNotFound, "empty metadata")
	}

	var query Query
	err := json.Unmarshal([]byte(md["opa"][0]), &query)
	if err != nil {
		log.Error(err, "get all params from query failed")
		return "", errorx.NewHTTP(err, http.StatusBadRequest, "fail on params request")
	}

	if len(query.Path) < 1 {
		log.Error(errQueryPathNotFound, "get path param from query failed")
		return "", errorx.NewHTTP(errNoMDFromCtx, http.StatusBadRequest, "fail to get param Path to request")
	}

	if len(query.Input) < 1 {
		log.Error(errQueryInputNotFound, "get input param from query failed")
		return "", errorx.NewHTTP(errNoMDFromCtx, http.StatusBadRequest, "fail to get param Input to request")
	}

	log.V(2).Info("query request FromGRPCCtx", "params", query)

	return a.do(ctx, query)
}

func (a auth) do(ctx context.Context, query Query) (string, error) {
	log := logx.WithName(ctx, "do() (string, error)")

	result, err := a.getOpaDecisionFromCurl(log, query)
	if err != nil {
		log.Error(err, "do() Error !")
		return "", err
	}

	str := fmt.Sprintf("%v", result)
	return str, nil
}

// getOpaDecisionFromCurl call opa by curl
// if OPA is unreachable or an other issues, return nil deccision with statusCode of the call and error-go
func (a auth) getOpaDecisionFromCurl(log logr.Logger, query Query) (string, error) {
	jsonStr, err := json.Marshal(query.Input)
	if err != nil {
		log.Error(err, "get opa param Input failed")
		return "", errorx.NewHTTP(err, http.StatusBadRequest, "get opa param Input failed")
	}

	status, statutcode, resp, err := a.callOpaServer(query.Path, `{"input":`+string(jsonStr)+`}`)
	if err != nil {
		log.Error(err, "get opa decision failed")
		return string(jsonStr), errorx.NewHTTP(err, statutcode, "get opa decision failed")
	}

	if statutcode != 200 {
		return string(jsonStr), errorx.NewHTTP(errDecisionOpaFromCurl, statutcode, "opa server not return a status code 200 but return  : "+status)
	}

	return resp, err
}
