package cart

import (
	"fmt"

	"github.com/kapeel-mopkar/ecom/types"
)

func getProductIdsFromCartItems(cartItems []types.CartCheckoutItem) ([]int, error) {
	productIds := make([]int, len(cartItems))

	for _, item := range cartItems {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %d", item.ProductID)
		}

		productIds = append(productIds, item.ProductID)
	}

	return productIds, nil
}

func (h *Handler) createOrder(products []types.Product, cartItems []types.CartCheckoutItem, userID int) (int, float64, error) {
	productsMap := make(map[int]types.Product)
	for _, product := range products {
		productsMap[product.ID] = product
	}

	//check if all products are in stock
	err := checkIfProductsInStock(cartItems, productsMap)
	if err != nil {
		return 0, 0, err
	}

	//calculate total price
	totalPrice := calculateTotalPrice(cartItems, productsMap)

	for _, item := range cartItems {
		product := productsMap[item.ProductID]
		product.Quantity -= item.Quantity
		h.productStore.UpdateProduct(product)
	}

	orderID, err := h.orderStore.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "Some address",
	})
	if err != nil {
		return 0, 0, err
	}
	for _, item := range cartItems {
		err = h.orderStore.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productsMap[item.ProductID].Price * float64(item.Quantity),
		})
		if err != nil {
			return 0, 0, err
		}
	}
	return orderID, totalPrice, nil
}

func checkIfProductsInStock(cartItems []types.CartCheckoutItem, productsMap map[int]types.Product) error {
	for _, item := range cartItems {
		product := productsMap[item.ProductID]
		if product.Quantity <= 0 {
			return fmt.Errorf("product %d is not available in store, please refresh your cart", item.ProductID)
		}
		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is not available in the quantity requested", product.Name)
		}
	}
	return nil
}

func calculateTotalPrice(cartItems []types.CartCheckoutItem, productsMap map[int]types.Product) float64 {
	var totalPrice float64
	for _, item := range cartItems {
		totalPrice += productsMap[item.ProductID].Price * float64(item.Quantity)
	}
	return totalPrice
}
