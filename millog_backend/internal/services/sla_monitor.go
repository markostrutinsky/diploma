package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"Omnilog_backend/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SLAMonitor struct {
	db        *pgxpool.Pool
	reqRepo   *repositories.SupplyRequestRepository
	auditRepo *repositories.AuditLogRepository
	emailSvc  EmailService
}

func NewSLAMonitor(db *pgxpool.Pool, reqRepo *repositories.SupplyRequestRepository, auditRepo *repositories.AuditLogRepository, emailSvc EmailService) *SLAMonitor {
	return &SLAMonitor{
		db:        db,
		reqRepo:   reqRepo,
		auditRepo: auditRepo,
		emailSvc:  emailSvc,
	}
}

// Start запускає фоновий процес (Cron)
func (s *SLAMonitor) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute) // ⏳ ДЛЯ ТЕСТУ: 1 хвилина. В продакшені зміниш на time.Hour
	log.Println("🛡️ SLA Monitor успішно запущено у фоновому режимі...")

	go func() {
		for {
			select {
			case <-ticker.C:
				s.CheckPendingRequests(ctx)
			case <-ctx.Done():
				ticker.Stop()
				log.Println("🛑 SLA Monitor зупинено.")
				return
			}
		}
	}()
}

// CheckPendingRequests зроблено з великої літери, щоб хендлер теж міг його викликати
func (s *SLAMonitor) CheckPendingRequests(ctx context.Context) (int, error) {
	overdueLimitHours := 24 // ⏳ ДЛЯ ТЕСТУ: постав 0

	requests, err := s.reqRepo.GetOverdueRequests(ctx, s.db, "PENDING", overdueLimitHours)
	if err != nil {
		log.Println("🚨 SLA Monitor: Помилка БД -", err)
		return 0, err
	}

	if len(requests) == 0 {
		return 0, nil
	}

	log.Printf("⚠️ SLA ALERT! Ескалація %d прострочених заявок!\n", len(requests))

	for _, req := range requests {
		// Тепер імейл динамічний!
		managerEmail := req.ManagerEmail
		details := fmt.Sprintf("Заявка очікує підтвердження понад %d годин. Автоматична ескалація.", overdueLimitHours)

		reqID := req.ID
		authorID := req.CreatedBy

		go func(rID, det, aID, targetEmail string) {
			bgCtx := context.Background()

			// 1. Зміна статусу на ESCALATED — тільки якщо заявка досі PENDING
			// Це захищає від ситуації, коли заявка вже скасована між вибіркою і оновленням
			escalated, err := s.reqRepo.EscalateStatus(bgCtx, s.db, rID, "Автоматична ескалація SLA")
			if err != nil {
				log.Println("❌ SLA Monitor: Помилка зміни статусу:", err)
				return
			}
			if !escalated {
				// Заявка вже не PENDING (скасована, відхилена, відправлена тощо) — пропускаємо
				log.Printf("⏭️ SLA Monitor: Заявка %s вже не в стані PENDING, ескалацію пропущено.\n", rID)
				return
			}

			// 2. Аудит
			err = s.auditRepo.LogAction(bgCtx, s.db, aID, "SLA_VIOLATION", "REQUEST", rID, det)
			if err != nil {
				log.Println("❌ Помилка запису в аудит SLA:", err)
			}

			// 3. Відправка Email на динамічну адресу
			err = s.emailSvc.SendSLAAlert(targetEmail, rID, overdueLimitHours)
			if err != nil {
				log.Printf("❌ Помилка відправки SLA Email на %s: %v\n", targetEmail, err)
			} else {
				log.Printf("✅ SLA Ескалація успішна! Сповіщення надіслано на: %s\n", targetEmail)
			}
		}(reqID, details, authorID, managerEmail)
	}

	return len(requests), nil
}

// GetEscalatedCount повертає поточну кількість заявок зі статусом ESCALATED
// (тих, що вже були ескальовані раніше і досі не опрацьовані).
func (s *SLAMonitor) GetEscalatedCount(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM supply_requests WHERE status = 'ESCALATED'`).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetPendingOverdueSoonCount повертає кількість заявок у статусі PENDING,
// які ще не досягли порогу ескалації, але вже мають вік > 50% від ліміту.
func (s *SLAMonitor) GetPendingStats(ctx context.Context) (pending int, overdueSoon int, err error) {
	err = s.db.QueryRow(ctx,
		`SELECT
		   COUNT(*) FILTER (WHERE status = 'PENDING'),
		   COUNT(*) FILTER (WHERE status = 'PENDING' AND created_at < NOW() - INTERVAL '12 hours')
		 FROM supply_requests`).Scan(&pending, &overdueSoon)
	return
}
