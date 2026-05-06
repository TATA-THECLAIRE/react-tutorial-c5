package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	"piggy.com/internal/models"
	"piggy.com/internal/piggyservice"
)

type Handler struct {
	service *piggyservice.Service
}

func NewHandler(service *piggyservice.Service) *Handler {
	return &Handler{service: service}
}

func getUserID(c *gin.Context) (pgtype.UUID, bool) {
	val, exists := c.Get("userID")
	if !exists {
		return pgtype.UUID{}, false
	}
	bytes, ok := val.([16]byte)
	if !ok {
		return pgtype.UUID{}, false
	}
	return pgtype.UUID{Bytes: bytes, Valid: true}, true
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	
	var payload models.CreateTransactionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.service.CreateTransaction(c, payload, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}


func (h *Handler) GetTransactions(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	txnType := c.Query("type") // "saving", "withdrawal", or ""
	size    := c.Query("size")  // e.g. "5" or ""

	transactions, err := h.service.GetTransactions(c, userID,txnType, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}


func (h *Handler) GetDashboardStats(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	stats, err := h.service.GetDashboardStats(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
