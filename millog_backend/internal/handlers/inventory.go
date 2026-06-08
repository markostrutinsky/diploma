package handlers

import (
	"Omnilog_backend/internal/middleware"
	"Omnilog_backend/internal/models"
	"Omnilog_backend/internal/services"
	"context"
	"fmt"
	"math"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type InventoryHandler struct {
	invService        *services.InventoryService
	auditService      *services.AuditService
	limitationService *services.LimitationService
	authService       *services.AuthService
}

func NewInventoryHandler(inv *services.InventoryService, audit *services.AuditService, limitation *services.LimitationService, auth *services.AuthService) *InventoryHandler {
	return &InventoryHandler{
		invService:        inv,
		auditService:      audit,
		limitationService: limitation,
		authService:       auth,
	}
}

func (h *InventoryHandler) CreateCategory(c *gin.Context) {
	userID := c.GetString("user_id")
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat, err := h.invService.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func(uID string, entityID string, name string) {
		_ = h.auditService.LogAction(context.Background(), uID, "CREATE", "CATEGORY", entityID, "Створено нову категорію: "+name)
	}(userID, fmt.Sprintf("%v", cat.ID), req.Name)

	c.JSON(http.StatusCreated, cat)
}

func (h *InventoryHandler) ListCategories(c *gin.Context) {
	list, err := h.invService.ListCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *InventoryHandler) CreateResource(c *gin.Context) {
	userID := c.GetString("user_id")
	var req models.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🛡️ Перевіримо ліміт на ресурсах
	if err := h.limitationService.CheckResourceLimit(c.Request.Context(), req.UnitID); err != nil {
		c.JSON(http.StatusPaymentRequired, gin.H{
			"error":       err.Error(),
			"upgrade_url": "/billing?plan=pro",
		})
		return
	}

	res, err := h.invService.CreateResource(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func(uID string, entityID string, name string) {
		_ = h.auditService.LogAction(context.Background(), uID, "CREATE", "RESOURCE", entityID, "Створено нову картку майна: "+name)
	}(userID, fmt.Sprintf("%v", res.ID), req.Name)

	c.JSON(http.StatusCreated, res)
}

func (h *InventoryHandler) GetResource(c *gin.Context) {
	resourceID := c.Param("id")
	if resourceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не вказано ID ресурсу"})
		return
	}

	res, err := h.invService.GetResource(c.Request.Context(), resourceID)
	if err != nil {
		if err.Error() == "repository get failed: resource not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "ресурс не знайдено"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *InventoryHandler) ListResources(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "claims not found"})
		return
	}
	claims := claimsVal.(*middleware.Claims)
	userRole := string(claims.Role)
	userUnitID := claims.UnitID

	var finalUnitID *int64

	var requestedUnitID *int64
	if s := c.Query("unit_id"); s != "" {
		if id, err := strconv.ParseInt(s, 10, 64); err == nil {
			requestedUnitID = &id
		}
	}

	// Перевіряємо чи це власник організації (ADMIN, TENANT_ADMIN, SYSTEM_ADMIN)
	isOwner := userRole == string(models.RoleAdmin) ||
		userRole == string(models.RoleTenantAdmin) ||
		userRole == string(models.RoleSystemAdmin)

	if isOwner {
		finalUnitID = requestedUnitID
	} else {
		if requestedUnitID != nil {
			finalUnitID = requestedUnitID
		} else {
			finalUnitID = &userUnitID
		}
	}

	list, err := h.invService.ListResources(c.Request.Context(), finalUnitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

// GetUniqueResourceNames повертає унікальні назви ресурсів (без дублювання по складах)
// для використання у формі створення заявки
func (h *InventoryHandler) GetUniqueResourceNames(c *gin.Context) {
	fmt.Println("🔍 GetUniqueResourceNames викликано")

	claims, exists := c.Get("claims")
	if !exists {
		fmt.Println("❌ User claims not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userClaims := claims.(*middleware.Claims)
	userRole := string(userClaims.Role)
	userUnitID := userClaims.UnitID

	fmt.Printf("✅ User: role=%s, unitID=%d\n", userRole, userUnitID)

	var finalUnitID *int64
	var requestedUnitID *int64
	if s := c.Query("unit_id"); s != "" {
		if id, err := strconv.ParseInt(s, 10, 64); err == nil {
			requestedUnitID = &id
		}
	}

	isOwner := userRole == string(models.RoleAdmin) ||
		userRole == string(models.RoleTenantAdmin) ||
		userRole == string(models.RoleSystemAdmin)

	if isOwner {
		finalUnitID = requestedUnitID
	} else {
		if requestedUnitID != nil {
			finalUnitID = requestedUnitID
		} else {
			finalUnitID = &userUnitID
		}
	}

	// Отримуємо всі ресурси
	list, err := h.invService.ListResources(c.Request.Context(), finalUnitID)
	if err != nil {
		fmt.Printf("❌ Error loading resources: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("📦 Loaded %d resources from DB\n", len(list))

	// Групуємо по назвах - зберігаємо тільки унікальні назви з їх категоріями
	type UniqueResource struct {
		Name       string `json:"name"`
		CategoryID string `json:"category_id"`
	}

	uniqueMap := make(map[string]UniqueResource)
	for _, res := range list {
		fmt.Printf("  - Resource: %s (category: %s)\n", res.Name, res.CategoryID)
		if _, exists := uniqueMap[res.Name]; !exists {
			uniqueMap[res.Name] = UniqueResource{
				Name:       res.Name,
				CategoryID: res.CategoryID,
			}
		}
	}

	fmt.Printf("✨ Found %d unique resource names\n", len(uniqueMap))

	// Перетворюємо map у slice
	unique := make([]UniqueResource, 0, len(uniqueMap))
	for _, ur := range uniqueMap {
		unique = append(unique, ur)
	}

	c.JSON(http.StatusOK, unique)
}

func (h *InventoryHandler) WriteOff(c *gin.Context) {
	resourceID := c.Param("id")
	userID := c.GetString("user_id")
	var req models.WriteOffResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невірний формат: вкажіть кількість (quantity)"})
		return
	}

	// 🛡️ Ієрархічна перевірка: BRANCH/DEPT рівень не може видавати з регіонального складу.
	// TENANT_ADMIN та SYSTEM_ADMIN (UnitID == 0) мають доступ до всіх ресурсів.
	claimsVal, _ := c.Get("claims")
	if claims, ok := claimsVal.(*middleware.Claims); ok && claims.UnitID != 0 {
		resource, err := h.invService.GetResource(c.Request.Context(), resourceID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ресурс не знайдено"})
			return
		}
		allowed, err := h.invService.IsUnitInSubtree(c.Request.Context(), claims.UnitID, resource.UnitID)
		if err != nil || !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Недостатньо прав: ресурс належить до іншого рівня ієрархії",
			})
			return
		}
	}

	err := h.invService.WriteOff(c.Request.Context(), resourceID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func(uID string, rID string) {
		_ = h.auditService.LogAction(
			context.Background(),
			uID,
			"WRITE_OFF",
			"RESOURCE",
			rID,
			"Списання майна зі складу",
		)
	}(userID, resourceID)

	c.JSON(http.StatusOK, gin.H{"message": "Успішно списано"})
}

func (h *InventoryHandler) UpdateResource(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resource id is required"})
		return
	}

	var req models.UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format: " + err.Error()})
		return
	}

	err := h.invService.UpdateResource(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "repository update failed: resource not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
			return
		}

		fmt.Println("🚨 UPDATE ERROR:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	go func(uID string, rID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "RESOURCE", rID, "Оновлено дані картки майна")
	}(userID, id)

	c.JSON(http.StatusOK, gin.H{"message": "resource successfully updated"})
}

func (h *InventoryHandler) Delete(c *gin.Context) {
	resourceID := c.Param("id")
	userID := c.GetString("user_id")
	if resourceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не вказано ID ресурсу"})
		return
	}

	err := h.invService.Delete(c.Request.Context(), resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	go func(uID string, rID string) {
		_ = h.auditService.LogAction(
			context.Background(),
			uID,
			"DELETE",
			"RESOURCE",
			rID,
			"Безповоротне видалення картки майна",
		)
	}(userID, resourceID)
	c.JSON(http.StatusOK, gin.H{"message": "Ресурс успішно видалено"})
}

func (h *InventoryHandler) Assign(c *gin.Context) {
	userID := c.GetString("user_id")
	resourceID := c.Param("id")
	var req models.AssignResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невірний формат запиту"})
		return
	}

	err := h.invService.Assign(c.Request.Context(), resourceID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func(uID string, rID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "RESOURCE", rID, "Майно видано користувачу")
	}(userID, resourceID)

	c.JSON(http.StatusOK, gin.H{"message": "Майно успішно видано"})
}

func (h *InventoryHandler) GetMyEquipment(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не знайдено токен авторизації"})
		return
	}
	claims := claimsVal.(*middleware.Claims)
	userID := claims.UserID

	items, err := h.invService.GetMyEquipment(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка отримання майна: " + err.Error()})
		return
	}

	if items == nil {
		items = []models.MyEquipmentItem{}
	}

	c.JSON(http.StatusOK, items)
}

func (h *InventoryHandler) IssueResource(c *gin.Context) {
	userID := c.GetString("user_id")
	var req models.IssueResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Невірний формат даних: " + err.Error()})
		return
	}

	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизовано"})
		return
	}
	claims := claimsVal.(*middleware.Claims)

	var unitID *int64
	if claims.UnitID != 0 {
		val := claims.UnitID
		unitID = &val
	}

	err := h.invService.IssueResource(c.Request.Context(), unitID, string(claims.Role), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func(uID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "RESOURCE", "", "Майно видано користувачу")
	}(userID)

	c.JSON(http.StatusOK, gin.H{"message": "Майно успішно видано користувачу"})
}

// ReportEquipment — рапорт користувача про майно (зламано/втрачено/зношено)
func (h *InventoryHandler) ReportEquipment(c *gin.Context) {
	userID := c.GetString("user_id")
	assignmentID := c.Param("id")

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Вкажіть причину запиту"})
		return
	}

	if err := h.invService.ReportEquipment(c.Request.Context(), userID, assignmentID, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func(uID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "REPORT", "RESOURCE_ASSIGNMENT", assignmentID, "Подано рапорт щодо майна: "+req.Reason)
	}(userID)

	c.JSON(http.StatusOK, gin.H{"message": "Рапорт успішно відправлено"})
}

func (h *InventoryHandler) CreateShipment(c *gin.Context) {
	userID := c.GetString("user_id")
	var req models.CreateShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Невірний формат даних: " + err.Error()})
		return
	}

	err := h.invService.CreateShipment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func(uID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "CREATE", "SHIPMENT", "", "Сформовано новий рейс")
	}(userID)

	c.JSON(http.StatusOK, gin.H{"message": "Рейс успішно сформовано, транспорт відправлено"})
}

func (h *InventoryHandler) GetByWarehouse(c *gin.Context) {
	warehouseID := c.Param("id")

	items, err := h.invService.GetByWarehouse(c.Request.Context(), warehouseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if items == nil {
		items = []models.InventoryItem{}
	}

	c.JSON(http.StatusOK, items)
}

func (h *InventoryHandler) ListShipments(c *gin.Context) {
	list, err := h.invService.ListShipments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *InventoryHandler) ListMyShipments(c *gin.Context) {
	userID := c.GetString("user_id")

	list, err := h.invService.ListMyShipments(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *InventoryHandler) ReceiveShipment(c *gin.Context) {
	userID := c.GetString("user_id")
	shipmentID := c.Param("id")

	// Зчитуємо фактичний пробіг і дані GPS для аудиту
	var body struct {
		ActualKm     float64 `json:"actual_km"`
		GpsKm        float64 `json:"gps_km"`        // пробіг по GPS (0 якщо немає)
		RouteStatus  string  `json:"route_status"`  // "on_route" | "deviated" | "unknown"
		DeviationPct float64 `json:"deviation_pct"` // відхилення у %
	}
	_ = c.ShouldBindJSON(&body)

	// Захист від явної фальсифікації: якщо є GPS і введений пробіг менший за 60% від GPS
	if body.GpsKm > 0 && body.ActualKm > 0 {
		minAllowed := body.GpsKm * 0.6
		if body.ActualKm < minAllowed {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Введений пробіг (%.1f км) значно менший за GPS-трек (%.1f км). Мінімально допустиме значення: %.1f км. Якщо є обґрунтування — зверніться до адміністратора.", body.ActualKm, body.GpsKm, minAllowed),
			})
			return
		}
	}

	err := h.invService.ReceiveShipment(c.Request.Context(), shipmentID, body.ActualKm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Формуємо деталі для аудит-логу
	auditNote := fmt.Sprintf("Вантаж прийнято. Пробіг: %.1f км", body.ActualKm)
	if body.GpsKm > 0 {
		auditNote += fmt.Sprintf(", GPS: %.1f км", body.GpsKm)
		if body.RouteStatus == "on_route" {
			auditNote += " ✅ маршрут виконано"
		} else if body.RouteStatus == "deviated" {
			auditNote += fmt.Sprintf(" ⚠️ відхилення %.0f%%", body.DeviationPct)
		}
	} else {
		auditNote += " (GPS не записувався)"
	}

	go func(uID string, sID string, note string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "SHIPMENT", sID, note)
	}(userID, shipmentID, auditNote)

	c.JSON(http.StatusOK, gin.H{"message": "Вантаж успішно прийнято на склад!"})
}

// GetShipmentGPSDistance — розраховує фактичний пробіг рейсу на основі GPS-треку.
// Повертає відстань по GPS (якщо є точки) і планову відстань OSRM для порівняння.
// GET /api/inventory/shipments/:id/gps-distance
func (h *InventoryHandler) GetShipmentGPSDistance(c *gin.Context) {
	shipmentID := c.Param("id")
	ctx := c.Request.Context()

	// Отримуємо інформацію про рейс: vehicle_id, started_at, distance_km
	var vehicleID string
	var startedAt *string
	var plannedKm float64
	err := h.invService.GetDB().QueryRow(ctx,
		`SELECT vehicle_id, TO_CHAR(started_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), COALESCE(distance_km, 0)
		 FROM shipments WHERE id = $1`, shipmentID,
	).Scan(&vehicleID, &startedAt, &plannedKm)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "рейс не знайдено"})
		return
	}

	// Якщо рейс ще не стартував — GPS точок немає
	if startedAt == nil || *startedAt == "" {
		c.JSON(http.StatusOK, gin.H{
			"has_gps":        false,
			"gps_km":         0,
			"points":         0,
			"planned_km":     plannedKm,
			"planned_rt_km":  math.Round(plannedKm*2*10) / 10,
			"suggested_km":   math.Round(plannedKm*2*10) / 10,
			"source":         "osrm",
			"route_status":   "unknown",
			"deviation_pct":  0,
			"min_allowed_km": 0,
		})
		return
	}

	// Запитуємо GPS точки з моменту старту рейсу до зараз
	rows, err := h.invService.GetDB().Query(ctx,
		`SELECT latitude, longitude FROM gps_locations
		 WHERE vehicle_id = $1 AND timestamp >= $2::timestamptz
		 ORDER BY timestamp ASC`, vehicleID, *startedAt,
	)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"has_gps": false, "gps_km": 0, "points": 0, "planned_km": plannedKm, "planned_rt_km": math.Round(plannedKm*2*10) / 10, "suggested_km": math.Round(plannedKm*2*10) / 10, "source": "osrm", "route_status": "unknown", "deviation_pct": 0, "min_allowed_km": 0})
		return
	}
	defer rows.Close()

	type point struct{ lat, lon float64 }
	var points []point
	for rows.Next() {
		var p point
		if rows.Scan(&p.lat, &p.lon) == nil {
			points = append(points, p)
		}
	}

	if len(points) < 2 {
		suggestedKm := math.Round(plannedKm*2*10) / 10
		c.JSON(http.StatusOK, gin.H{
			"has_gps":        false,
			"gps_km":         0,
			"points":         len(points),
			"planned_km":     plannedKm,
			"planned_rt_km":  suggestedKm,
			"suggested_km":   suggestedKm,
			"source":         "osrm",
			"route_status":   "unknown",
			"deviation_pct":  0,
			"min_allowed_km": 0,
		})
		return
	}

	// Haversine по всіх точках треку
	const R = 6371.0
	totalKm := 0.0
	for i := 1; i < len(points); i++ {
		lat1 := points[i-1].lat * math.Pi / 180
		lon1 := points[i-1].lon * math.Pi / 180
		lat2 := points[i].lat * math.Pi / 180
		lon2 := points[i].lon * math.Pi / 180
		dlat := lat2 - lat1
		dlon := lon2 - lon1
		a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
		c2 := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
		totalKm += R * c2
	}
	gpsKm := math.Round(totalKm*10) / 10 // округлення до 0.1 км

	// Порівнюємо GPS з плановим маршрутом (туди і назад = planned * 2)
	plannedRoundTrip := plannedKm * 2
	deviationPct := 0.0
	routeStatus := "unknown"
	if plannedRoundTrip > 0 && gpsKm > 0 {
		deviationPct = math.Abs(gpsKm-plannedRoundTrip) / plannedRoundTrip * 100
		if deviationPct <= 20 {
			routeStatus = "on_route" // ±20% — вважається що маршрут виконаний
		} else {
			routeStatus = "deviated"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"has_gps":        true,
		"gps_km":         gpsKm,
		"points":         len(points),
		"planned_km":     plannedKm,
		"planned_rt_km":  math.Round(plannedRoundTrip*10) / 10,
		"suggested_km":   gpsKm,
		"source":         "gps",
		"route_status":   routeStatus,
		"deviation_pct":  math.Round(deviationPct*10) / 10,
		"min_allowed_km": math.Round(gpsKm*0.6*10) / 10,
	})
}

func (h *InventoryHandler) StartShipment(c *gin.Context) {
	userID := c.GetString("user_id")
	shipmentID := c.Param("id")

	err := h.invService.StartShipment(c.Request.Context(), shipmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func(uID string, sID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "SHIPMENT", sID, "Рейс розпочато (виїзд підтверджено)")
	}(userID, shipmentID)

	c.JSON(http.StatusOK, gin.H{"message": "Рейс розпочато! Гарної дороги 🚚"})
}

// LogShipmentRefuel — водій реєструє дозаправку під час рейсу.
// POST /api/inventory/shipments/:id/refuel
func (h *InventoryHandler) LogShipmentRefuel(c *gin.Context) {
	userID := c.GetString("user_id")
	shipmentID := c.Param("id")

	var req models.LogShipmentRefuelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некоректні дані: " + err.Error()})
		return
	}

	refuel, err := h.invService.LogShipmentRefuel(c.Request.Context(), shipmentID, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func(uID string, sID string, liters float64) {
		note := fmt.Sprintf("Дозаправка в дорозі: %.1f л", liters)
		_ = h.auditService.LogAction(context.Background(), uID, "CREATE", "SHIPMENT_REFUEL", sID, note)
	}(userID, shipmentID, req.Liters)

	c.JSON(http.StatusCreated, refuel)
}

// GetShipmentRefuels — список усіх дозаправок конкретного рейсу.
// GET /api/inventory/shipments/:id/refuels
func (h *InventoryHandler) GetShipmentRefuels(c *gin.Context) {
	shipmentID := c.Param("id")

	refuels, err := h.invService.GetShipmentRefuels(c.Request.Context(), shipmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if refuels == nil {
		refuels = []*models.ShipmentRefuel{}
	}
	c.JSON(http.StatusOK, refuels)
}

func (h *InventoryHandler) DownloadShipmentPDF(c *gin.Context) {
	shipmentID := c.Param("id")

	pdfBytes, err := h.invService.GenerateShipmentPDF(c.Request.Context(), shipmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := fmt.Sprintf("Waybill_%s.pdf", shipmentID)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "application/pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func (h *InventoryHandler) DownloadResourceQR(c *gin.Context) {
	resourceID := c.Param("id")

	pngBytes, err := h.invService.GenerateResourceQR(resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileName := fmt.Sprintf("QR_%s.png", resourceID)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Header("Content-Type", "image/png")

	c.Data(http.StatusOK, "image/png", pngBytes)
}

func (h *InventoryHandler) UpdateCategory(c *gin.Context) {
	categoryID := c.Param("id")
	userID := c.GetString("user_id")

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильні дані запиту"})
		return
	}

	err := h.invService.UpdateCategory(c.Request.Context(), categoryID, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не вдалося оновити категорію"})
		return
	}

	go func(uID, cID, name string) {
		_ = h.auditService.LogAction(context.Background(), uID, "UPDATE", "CATEGORY", cID, "Оновлено назву/опис категорії: "+name)
	}(userID, categoryID, req.Name)

	c.JSON(http.StatusOK, gin.H{"message": "Категорію успішно оновлено"})
}

func (h *InventoryHandler) DeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")
	userID := c.GetString("user_id")

	err := h.invService.DeleteCategory(c.Request.Context(), categoryID)
	if err != nil {
		if err.Error() == "категорія не порожня" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неможливо видалити: у цій категорії ще є майно"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка видалення категорії"})
		return
	}

	go func(uID, cID string) {
		_ = h.auditService.LogAction(context.Background(), uID, "DELETE", "CATEGORY", cID, "Видалено категорію майна")
	}(userID, categoryID)

	c.JSON(http.StatusOK, gin.H{"message": "Категорію видалено"})
}

func (h *InventoryHandler) SubmitAudit(c *gin.Context) {
	var req models.SubmitAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Невірний формат даних: " + err.Error()})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не знайдено ID користувача"})
		return
	}

	err := h.invService.SubmitInventoryAudit(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка збереження результатів переобліку: " + err.Error()})
		return
	}

	go func(uID, wID string) {
		_ = h.auditService.LogAction(
			context.Background(),
			uID,
			"INVENTORY_AUDIT",
			"WAREHOUSE",
			wID,
			"Проведено переоблік складу. Зафіксовано акт розбіжностей.",
		)
	}(userID, req.WarehouseID)

	c.JSON(http.StatusOK, gin.H{"message": "Акт переобліку успішно сформовано, залишки оновлено!"})
}

// Допоміжна функція для очищення тексту перед порівнянням
func normalizeText(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// 1. Простий і красивий шаблон (без організацій і складів)
func (h *InventoryHandler) DownloadImportTemplate(c *gin.Context) {
	f := excelize.NewFile()
	defer f.Close()
	sheet := "Шаблон"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{
		"Назва майна (Обов'язково)",
		"Назва категорії (Обов'язково)",
		"Кількість",
		"Одиниці (PCS/KG/L/KIT)",
		"Вага 1 од. (кг)",
		"Заводський штрих-код",
		"Мін. залишок",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	f.SetColWidth(sheet, "A", "B", 35)
	f.SetColWidth(sheet, "C", "G", 15)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=OmniLog_Import_Template.xlsx")
	f.Write(c.Writer)
}

// 1. Оновлений Імпорт (з виводом помилок у консоль)
func (h *InventoryHandler) ImportExcel(c *gin.Context) {
	unitIDStr := c.PostForm("unit_id")
	warehouseID := c.PostForm("warehouse_id")
	unitID, _ := strconv.ParseInt(unitIDStr, 10, 64)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не знайдено"})
		return
	}

	// 🛡️ Ліміт розміру Excel: 5 МБ
	const maxExcelSize = 5 << 20 // 5 MB
	if fileHeader.Size > maxExcelSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл занадто великий. Максимальний розмір для імпорту — 5 МБ."})
		return
	}

	// 🛡️ Валідація розширення
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Дозволяються лише файли .xlsx або .xls"})
		return
	}

	file, _ := fileHeader.Open()
	defer file.Close()
	f, _ := excelize.OpenReader(file)
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil || len(rows) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл порожній"})
		return
	}

	ctx := c.Request.Context()
	categories, err := h.invService.GetAllCategories(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Критична помилка завантаження категорій з бази: %v", err),
		})
		return
	}

	successCount := 0
	var importErrors []string

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		getCol := func(idx int) string {
			if idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		name := getCol(0)
		categoryName := getCol(1)
		if name == "" || categoryName == "" {
			continue
		}

		// 1. Шукаємо або створюємо категорію (безпечним методом)
		catID, catErr := h.smartCategoryLookup(ctx, categoryName, &categories)
		if catErr != nil {
			importErrors = append(importErrors, fmt.Sprintf("Рядок %d (%s): %v", i+1, name, catErr))
			continue // Блокуємо і йдемо до наступного
		}

		qty, _ := strconv.Atoi(getCol(2))
		unitType := getCol(3)
		if unitType == "" {
			unitType = "PCS"
		}

		// Захист від ком замість крапок (напр. 0,6 замість 0.6)
		weightStr := strings.ReplaceAll(getCol(4), ",", ".")
		weight, _ := strconv.ParseFloat(weightStr, 64)
		if weight <= 0 {
			weight = 1.0
		}
		minQty, _ := strconv.Atoi(getCol(6))

		// 2. Перевіряємо, чи існує вже такий ресурс
		existingRes, err := h.invService.GetResourceByNameAndWarehouse(ctx, name, warehouseID)

		if err == nil && existingRes != nil {
			// РЕСУРС ІСНУЄ -> Плюсуємо кількість
			newQty := existingRes.Quantity + qty
			updateErr := h.invService.UpdateResourceQuantity(ctx, existingRes.ID, newQty)
			if updateErr != nil {
				importErrors = append(importErrors, fmt.Sprintf("Рядок %d: Помилка оновлення кількості: %v", i+1, updateErr))
			} else {
				successCount++
			}
		} else {
			// РЕСУРС НОВИЙ -> Створюємо
			req := models.CreateResourceRequest{
				UnitID:      unitID,
				WarehouseID: &warehouseID,
				CategoryID:  catID,
				Name:        name,
				Quantity:    qty,
				UnitType:    models.UnitMeasurement(unitType),
				WeightKg:    weight,
				Barcode:     getCol(5),
				MinQuantity: minQty,
				Condition:   models.ConditionNew,
			}

			_, createErr := h.invService.CreateResource(ctx, &req)
			if createErr != nil {
				importErrors = append(importErrors, fmt.Sprintf("Рядок %d (%s): Помилка створення БД - %v", i+1, name, createErr))
			} else {
				successCount++
			}
		}
	}

	// 🔥 ДРУКУЄМО ПОМИЛКИ В КОНСОЛЬ GO ДЛЯ ДЕБАГУ
	if len(importErrors) > 0 {
		fmt.Println("=== ПОМИЛКИ ІМПОРТУ ===")
		for _, e := range importErrors {
			fmt.Println(e)
		}
		fmt.Println("=======================")
	}

	c.JSON(http.StatusOK, gin.H{
		"success_count":   successCount,
		"total_processed": len(rows) - 1,
		"errors":          importErrors,
	})
}

// 2. Виправлений, безпечний розумний пошук
func (h *InventoryHandler) smartCategoryLookup(ctx context.Context, name string, cats *[]models.ResourceCategory) (string, error) {
	search := strings.ToLower(strings.TrimSpace(name))
	if search == "" {
		return "", fmt.Errorf("порожня назва категорії")
	}

	// Етап 1: Шукаємо ТОЧНИЙ збіг (пріоритет)
	for _, c := range *cats {
		if strings.ToLower(c.Name) == search {
			return c.ID, nil
		}
	}

	// Етап 2: Безпечний частковий збіг
	// (Ігноруємо слова менше 4 символів, щоб не ловити "та", "іт", "ка")
	for _, c := range *cats {
		dbName := strings.ToLower(c.Name)
		if len(dbName) > 3 && strings.Contains(search, dbName) {
			return c.ID, nil
		}
		if len(search) > 3 && strings.Contains(dbName, search) {
			return c.ID, nil
		}
	}

	// Етап 3: Якщо нічого не підійшло - створюємо нову
	newCat, err := h.invService.CreateCategory(ctx, &models.CreateCategoryRequest{Name: name})
	if err != nil {
		return "", fmt.Errorf("БД відмовила у створенні категорії: %v", err)
	}
	if newCat != nil && newCat.ID != "" {
		newCat.Name = name
		*cats = append(*cats, *newCat) // Запам'ятовуємо, щоб наступні рядки її використовували
		return newCat.ID, nil
	}

	return "", fmt.Errorf("невідома помилка отримання ID нової категорії")
}

// GetShipmentsByVehicle — повертає список рейсів для конкретного авто
func (h *InventoryHandler) GetShipmentsByVehicle(c *gin.Context) {
	vehicleID := c.Param("id")
	if vehicleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle id is required"})
		return
	}
	list, err := h.invService.ListShipmentsByVehicle(c.Request.Context(), vehicleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка завантаження рейсів"})
		return
	}
	c.JSON(http.StatusOK, list)
}
