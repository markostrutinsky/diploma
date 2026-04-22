package repositories

import "context"

// tenantCtxKey має бути ідентичний тому, що використовує middleware.
// Щоб уникнути циклу імпортів, ми визначаємо власний тип-"приватний" ключ,
// але порівняння йде по значенню string, покладеному під ключем middleware.
// Тому приймаємо значення через helper-функцію, яку handler/service встановить явно.
type tenantCtxKeyT struct{}

var TenantCtxKey = tenantCtxKeyT{}

// WithTenant кладе tenant_id у context (використовується в тестах / сервісах).
func WithTenant(ctx context.Context, tenantID string) context.Context {
	if tenantID == "" {
		return ctx
	}
	return context.WithValue(ctx, TenantCtxKey, tenantID)
}

// TenantFromCtx читає tenant_id із context; якщо його немає — повертає "".
// "" означає «SYSTEM_ADMIN або unauth» — repository такий запит не фільтрує.
func TenantFromCtx(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v := ctx.Value(TenantCtxKey)
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}
