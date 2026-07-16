package transaction

import "errors"

var ErrAmountIsZero = errors.New("amount is zero")
var ErrInvalidCategory = errors.New("invalid category")
var ErrInternalWhileCreateTransaction = errors.New("internal error while create transaction")
var ErrTransactionAlreadyCreated = errors.New("transaction already created")
var ErrInternalWhileGetWalletBalance = errors.New("internal error while get wallet balance")
