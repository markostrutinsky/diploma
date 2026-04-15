package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"millog_backend/internal/models"
	"millog_backend/internal/services"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	auditService *services.AuditService
}

// NewAuditHandler створює новий екземпляр хендлера з інжектованим сервісом
func NewAuditHandler(auditService *services.AuditService) *AuditHandler {
	return &AuditHandler{auditService: auditService}
}

// GetLogs обробляє GET-запит для отримання журналу дій
func (h *AuditHandler) GetLogs(c *gin.Context) {
	// 1. Беремо роль як "будь-що" (interface{}) за обома можливими ключами
	roleAny, exists := c.Get("user_role")
	if !exists {
		roleAny, _ = c.Get("role") // Запасний варіант
	}

	// 2. Жорстко конвертуємо в рядок, прибираємо зайві пробіли і робимо ВЕЛИКИМИ літерами
	role := strings.TrimSpace(strings.ToUpper(fmt.Sprintf("%v", roleAny)))

	// 3. Друкуємо в термінал Докера для дебагу
	log.Printf("🔥 [АУДИТ] Перевірка ролі: '%s' (Оригінал: %v)", role, roleAny)

	// 4. Перевіряємо
	if role != "ADMIN" {
		// Тепер, якщо буде помилка, ти прямо на екрані побачиш, ЯКУ САМЕ роль бачить сервер
		c.JSON(http.StatusForbidden, gin.H{
			"error": fmt.Sprintf("Доступ заборонено. Сервер бачить вашу роль як: '%s'", role),
		})
		return
	}

	// ... тут далі твій код з limit та викликом сервісу ...
	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	logs, err := h.auditService.GetLogs(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не вдалося завантажити журнал аудиту"})
		return
	}

	if logs == nil {
		logs = []models.AuditLog{}
	}

	c.JSON(http.StatusOK, logs)
}
