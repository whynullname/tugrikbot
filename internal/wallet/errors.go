package wallet

import "errors"

var ErrInternalWhileGetUserWalletID = errors.New("internal error while get user wallet id")
var ErrInternalWhileCreateNewWallet = errors.New("internal error while create new wallet")
var ErrUserAlreadyInWallet = errors.New("user already in wallet")
var ErrInvalidInviteCode = errors.New("invalid invite code")
var ErrInternalWhileJoinToWallet = errors.New("internal error while join to wallet")
var ErrInternalWhileGetWallets = errors.New("internal error while get wallets")
var ErrUserNotInWallet = errors.New("user not in wallet")
var ErrInvalidWalletID = errors.New("invalid wallet id")
var ErrInternalWhileSetActiveWallet = errors.New("internal error while set active wallet")
