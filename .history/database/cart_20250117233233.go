package database

import "errors"

var (
	ErrCantFindProduct    = errors.New("can't find the product")
	ErrCantDecodeProduct  = errors.New("can't find the product")
	ErrUserIdIsNotValid   = errors.New("this user is not valid")
	ErrCantUpdateUser     = errors.New("cannot add this product to the cart")
	ErrCantRemoveCartItem = errors.New("cannot remove this item from cart")
	ErrCantGetItem        = errors.New("unable to get the item from the cart")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
)

func AddProductToCart() {}

func RemoveCartItem() {}

func BuyItemFromCart() {}

func InstantBuyer()
