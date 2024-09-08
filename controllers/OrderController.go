package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hacktiv8-go-assignment-2/database"
	"hacktiv8-go-assignment-2/models"
	"net/http"
	"time"
)

type CreateOrderRequest struct {
	OrderedAt    string     `json:"ordered_at" binding:"required"`
	CustomerName string     `json:"customer_name" binding:"required"`
	Items        []ItemData `json:"items" binding:"required,dive"`
}
type UpdateOrderRequest struct {
	OrderedAt    string        `json:"ordered_at" binding:"required"`
	CustomerName string        `json:"customer_name" binding:"required"`
	Items        []models.Item `json:"items" binding:"required,dive"`
}

type ItemData struct {
	ItemCode    string `json:"item_code" binding:"required"`
	Description string `json:"description" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required"`
}

func CreateOrder(c *gin.Context) {
	var newOrder CreateOrderRequest

	if err := c.ShouldBindJSON(&newOrder); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	orderedAt, err := time.Parse(time.RFC3339, newOrder.OrderedAt)
	if err != nil {
		fmt.Println("Invalid date format for OrderedAt:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date format for OrderedAt",
		})
		return
	}

	order := models.Order{
		OrderedAt:    orderedAt,
		CustomerName: newOrder.CustomerName,
		Items:        make([]models.Item, len(newOrder.Items)),
	}

	for i, ItemData := range newOrder.Items {
		order.Items[i] = models.Item{
			ItemCode:    ItemData.ItemCode,
			Description: ItemData.Description,
			Quantity:    uint(ItemData.Quantity),
		}
	}
	db := database.GetDB()

	err = db.Create(&order).Error

	if err != nil {
		fmt.Println("Error creating create order:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating create order",
		})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func GetOrders(c *gin.Context) {
	var newOrders []models.Order

	db := database.GetDB()

	err := db.Preload("Items").Find(&newOrders).Error
	if err != nil {
		fmt.Println("Error getting Orders:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error getting Orders",
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"data": newOrders,
	})
}

func GetOrderById(c *gin.Context) {
	orderID := c.Param("id")
	var newOrder models.Order

	db := database.GetDB()

	err := db.Preload("Items").First(&newOrder, "id = ?", orderID).Error
	if err != nil {
		fmt.Println("Error order not found:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error order not found",
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"data": newOrder,
	})
}
func UpdateOrderById(c *gin.Context) {
	orderID := c.Param("id")
	var newOrder UpdateOrderRequest

	if err := c.ShouldBindJSON(&newOrder); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	orderedAt, err := time.Parse(time.RFC3339, newOrder.OrderedAt)
	if err != nil {
		fmt.Println("Invalid date format for OrderedAt:", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date format for OrderedAt",
		})
		return
	}
	db := database.GetDB()
	tx := db.Begin()

	// find the existing order
	var order models.Order

	if err := tx.First(&order, orderID).Error; err != nil {
		tx.Rollback()
		fmt.Println("Error Order not found:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// update order field
	order.CustomerName = newOrder.CustomerName
	order.OrderedAt = orderedAt

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		fmt.Println("Error updating order:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating order"})
		return
	}
	// Clear old items
	if err := tx.Where("order_id = ?", orderID).Delete(&models.Item{}).Error; err != nil {
		tx.Rollback()
		fmt.Println("Error clearing old items", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error clearing old items"})
		return
	}
	// handle items
	for _, ItemData := range newOrder.Items {
		item := models.Item{
			ID:          ItemData.ID,
			ItemCode:    ItemData.ItemCode,
			Description: ItemData.Description,
			Quantity:    uint(ItemData.Quantity),
			OrderID:     order.ID,
		}
		if err := tx.Save(&item).Error; err != nil {
			tx.Rollback()
			fmt.Println("Error updating item", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating item"})
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		fmt.Println("Error committing transaction", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error committing transaction"})
		return
	}

	// Reload the order with associated items
	if err := db.Preload("Items").First(&order, orderID).Error; err != nil {
		fmt.Println("Error reloading order with items", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reloading order with items"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func DeleteOrderById(c *gin.Context) {
	orderID := c.Param("id")

	db := database.GetDB()
	tx := db.Begin()

	// find the existing order
	var order models.Order

	if err := tx.First(&order, orderID).Error; err != nil {
		tx.Rollback()
		fmt.Println("Error Order not found:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// delete items before order
	if err := tx.Where("order_id = ?", orderID).Delete(&models.Item{}).Error; err != nil {
		tx.Rollback()
		fmt.Println("Error deleting order item", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Error deleting order item"})
		return
	}

	// delete order
	if err := tx.Delete(&order).Error; err != nil {
		tx.Rollback()
		fmt.Println("Error deleting order", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Error deleting order"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		fmt.Println("Error committing transaction", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error committing transaction"})
		return
	}

	c.SecureJSON(http.StatusNoContent, nil)
}
