package handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/caplo84/quizz-backend/internal/services"
    "github.com/caplo84/quizz-backend/internal/models"
)

type AdminHandler struct {
    adminService services.AdminService
}

func NewAdminHandler(adminService services.AdminService) *AdminHandler {
    return &AdminHandler{
        adminService: adminService,
    }
}

// CreateQuiz handles POST /admin/quizzes
func (h *AdminHandler) CreateQuiz(c *gin.Context) {
    var quiz models.Quiz
    
    if err := c.ShouldBindJSON(&quiz); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
        })
        return
    }
    
    if err := h.adminService.CreateQuiz(c.Request.Context(), &quiz); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create quiz",
        })
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "data": quiz,
    })
}

// UpdateQuiz handles PUT /admin/quizzes/:id
func (h *AdminHandler) UpdateQuiz(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid quiz ID",
        })
        return
    }
    
    var quiz models.Quiz
    if err := c.ShouldBindJSON(&quiz); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
        })
        return
    }
    
    quiz.ID = uint(id)
    
    if err := h.adminService.UpdateQuiz(c.Request.Context(), &quiz); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to update quiz",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "data": quiz,
    })
}

// DeleteQuiz handles DELETE /admin/quizzes/:id
func (h *AdminHandler) DeleteQuiz(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid quiz ID",
        })
        return
    }
    
    if err := h.adminService.DeleteQuiz(c.Request.Context(), uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to delete quiz",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Quiz deleted successfully",
    })
}