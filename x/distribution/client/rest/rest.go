package rest

import (
	client2 "github.com/cosmos/cosmos-sdk/x/staking/client"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
)

// RegisterRoutes register distribution REST routes.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	registerQueryRoutes(cliCtx, r, queryRoute)
	registerTxRoutes(cliCtx, r, queryRoute)
}

// ProposalRESTHandler returns a ProposalRESTHandler that exposes the community pool spend REST handler with a given sub-route.
func ProposalRESTHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "community_pool_spend",
		Handler:  postProposalHandlerFn(cliCtx),
	}
}

func postProposalHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CommunityPoolSpendProposalReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		// HashGard
		if err := client2.QueryBondedRatioByClient(cliCtx); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		content := types.NewCommunityPoolSpendProposal(req.Title, req.Description, req.Recipient, req.Amount)

		msg := gov.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// HashGard
// ProposalStakeIssueLockedRESTHandler returns a ProposalRESTHandler that exposes the stake issue locked spend REST handler with a given sub-route.
func ProposalStakeIssueLockedRESTHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "stake_issue_locked_spend",
		Handler:  postProposalStakeIssueLockedHandlerFn(cliCtx),
	}
}

func postProposalStakeIssueLockedHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req StakeIssueLockedSpendProposalReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		content := types.NewStakeIssueLockedSpendProposal(req.Title, req.Description, req.Denom, req.Recipient, req.Amount)

		msg := gov.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
