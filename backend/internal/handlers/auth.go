package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"

    "piggy.com/internal/auth"
	
)

type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3"`
    Email    string `json:"email"    binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func (h *Handler) Register(ctx *gin.Context) {
    var req RegisterRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
        return
    }

    user, err := h.service.CreateUser(ctx, req.Username, req.Email, string(hashed))
    if err != nil {
        ctx.JSON(http.StatusConflict, gin.H{"error": "username or email already in use"})
        return
    }

    token, err := auth.GenerateToken(user.ID, user.Username)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
        return
    }

    ctx.JSON(http.StatusCreated, gin.H{"token": token, "user": user})
}

func (h *Handler) Login(ctx *gin.Context) {
    var req LoginRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.service.GetUserByUsername(ctx, req.Username)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    token, err := auth.GenerateToken(user.ID, user.Username)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}