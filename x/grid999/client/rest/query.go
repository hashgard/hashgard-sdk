package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/grid999/internal/types"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	List      = "list"
	Detail    = "detail"
	DappId    = "dappId"
	ID        = "id"
	numLatest = "numLatest"
	limit     = "limit"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	r.HandleFunc(fmt.Sprintf("/%s/%s", queryRoute, types.QueryParams), queryParams(queryRoute, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/dapp/%s/{%s}/{%s}", queryRoute, List, numLatest, limit), queryGrid999DappList(queryRoute, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/dapp/%s/{%s}", queryRoute, Detail, ID), queryGrid999Dapp(queryRoute, cliCtx)).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/%s/grid/%s/{%s}/{%s}/{%s}", queryRoute, List, DappId, numLatest, limit), queryGrid999List(queryRoute, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/grid/%s/{%s}/{%s}", queryRoute, Detail, DappId, ID), queryGrid999(queryRoute, cliCtx)).Methods("GET")
}
func queryParams(queryRoute string, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryParams), nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryGrid999DappList(queryRoute string, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s", queryRoute, types.QueryDappList,
			vars[numLatest], vars[limit]), nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
func queryGrid999Dapp(queryRoute string, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryDapp,
			vars[ID]), nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryGrid999List(queryRoute string, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s/%s", queryRoute, types.QueryGrid999List,
			vars[DappId], vars[numLatest], vars[limit]), nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
func queryGrid999(queryRoute string, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s", queryRoute, types.QueryGrid999,
			vars[DappId], vars[ID]), nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
