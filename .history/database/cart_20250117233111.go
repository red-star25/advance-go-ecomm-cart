package database

import "errors"

var (
	ErrCantFindProduct    = errors.New("")
	ErrCantDecodeProduct  = errors.New("")
	ErrUserIdIsNotValid   = errors.New("")
	ErrCantUpdateUser     = errors.New("")
	ErrCantRemoveCartItem = errors.New("")
	ErrCantGetItem        = errors.New("")
	ErrCantBuyCartItem    = errors.New("")
)

func AddProductToCart() {}

func RemoveCartItem() {}

func BuyItemFromCart() {}

func InstantBuyer()
