package models

import "time"

// ContractorMembershipStatus — стан співпраці підрядника з конкретною організацією (tenant).
type ContractorMembershipStatus string

const (
	// MembershipPending — підрядник подав заявку (або спробував узяти завдання),
	// організація ще не прийняла рішення.
	MembershipPending ContractorMembershipStatus = "PENDING"
	// MembershipApproved — організація підтвердила підрядника; він може брати її завдання.
	MembershipApproved ContractorMembershipStatus = "APPROVED"
	// MembershipRejected — організація відхилила співпрацю.
	MembershipRejected ContractorMembershipStatus = "REJECTED"
)

// ContractorMembership — зв'язок «підрядник ↔ організація» з власним статусом схвалення.
// Дозволяє одному глобальному підряднику співпрацювати з кількома організаціями,
// при цьому кожна організація самостійно вирішує, чи допускати його до своїх завдань.
type ContractorMembership struct {
	ID           string                     `json:"id"`
	ContractorID string                     `json:"contractor_id"`
	TenantID     string                     `json:"tenant_id"`
	Status       ContractorMembershipStatus `json:"status"`
	Note         *string                    `json:"note,omitempty"`
	RequestedAt  time.Time                  `json:"requested_at"`
	DecidedAt    *time.Time                 `json:"decided_at,omitempty"`
	DecidedBy    *string                    `json:"decided_by,omitempty"`

	// Обчислювані поля (JOIN на users) — для адмін-панелі організації.
	ContractorName  string  `json:"contractor_name,omitempty"`
	ContractorEmail string  `json:"contractor_email,omitempty"`
	ContractorPhone *string `json:"contractor_phone,omitempty"`

	// Обчислюване поле (JOIN на tenants) — для self-view підрядника.
	TenantName string `json:"tenant_name,omitempty"`
}
