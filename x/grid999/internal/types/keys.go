package types

// nolint
const (
	// module name
	ModuleName = "grid999"

	// default paramspace for params keeper
	DefaultParamspace = ModuleName

	// StoreKey is the default store key for mint
	StoreKey = ModuleName

	// RouterKey is the message route for distribution
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = StoreKey
)
const (
	MaxGrid           = 9
	NumberDigit       = 9
	TypeFuturePayment = "future_payment"
	TypeVote          = "vote"
	TypeLock          = "lock"
	TypePrepaid       = "prepaid"
)
