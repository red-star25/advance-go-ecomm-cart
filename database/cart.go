package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/red-star25/advance-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("can't find the product")
	ErrCantDecodeProduct  = errors.New("can't find the product")
	ErrUserIdIsNotValid   = errors.New("this user is not valid")
	ErrCantUpdateUser     = errors.New("cannot add this product to the cart")
	ErrCantRemoveCartItem = errors.New("cannot remove this item from cart")
	ErrCantGetItem        = errors.New("unable to get the item from the cart")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchFromDB, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}
	var productCart []models.ProductUser
	err = searchFromDB.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProduct
	}
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}

func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}

	_, err = prodCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantRemoveCartItem
	}
	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	// fetch the cart of the user
	// find the cart total
	// create an order with the items
	// added order to the user collection
	// added items in the cart to order list
	// empty up the cart
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var getCartItems models.User
	var orderCart models.Order

	orderCart.OrderID = primitive.NewObjectID()
	orderCart.OrderedAt = time.Now()
	orderCart.OrderCart = make([]models.ProductUser, 0)
	orderCart.PaymentMethod.COD = true

	unwind := bson.D{{Key: "$unwind", Value: primitive.E{Key: "path", Value: "$usercart"}}}
	grouping := bson.D{{Key: "$group", Value: primitive.E{Key: "_id", Value: "$_id"}}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}

	currentResults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}

	var getUserCart []bson.M
	if err := currentResults.All(ctx, &getUserCart); err != nil {
		panic(err)
	}

	var totalPrice int32

	for _, userItem := range getUserCart {
		price := userItem["total"]
		totalPrice = price.(int32)
	}

	orderCart.Price = int(totalPrice)

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderCart}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getCartItems)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getCartItems.UserCart}}}

	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}

	userCartEmpty := make([]models.ProductUser, 0)

	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: userCartEmpty}}}}

	_, err = userCollection.UpdateOne(ctx, filter3, update3)
	if err != nil {
		return ErrCantBuyCartItem
	}
	return nil
}

func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, UserID string) error {
	id, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	var product_details models.ProductUser
	var orders_detail models.Order
	orders_detail.OrderID = primitive.NewObjectID()
	orders_detail.OrderedAt = time.Now()
	orders_detail.OrderCart = make([]models.ProductUser, 0)
	orders_detail.PaymentMethod.COD = true
	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&product_details)
	if err != nil {
		log.Println(err)
	}
	orders_detail.Price = product_details.Price
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orders_detail}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": product_details}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}
	return nil
}
