package dto

// ── User data extracted from JWT ──
type UserData struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// ── Standard API Response ──
type APIResponse struct {
	Status     bool   `json:"status"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
}

// ── Dashboard HTTP Request DTOs (from frontend) ──

type DateOption struct {
	Date  *string    `json:"date,omitempty"`
	Year  *int       `json:"year,omitempty"`
	Month *int       `json:"month,omitempty"`
	Day   *int       `json:"day,omitempty"`
	Range *DateRange `json:"range,omitempty"`
}

type DateRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type GetUserTransactionsRequest struct {
	WalletID   string     `json:"walletID,omitempty"`
	DateOption DateOption `json:"dateOption"`
}

type GetUserBalanceRequest struct {
	WalletID    string     `json:"walletID,omitempty"`
	Aggregation string     `json:"aggregation" validate:"required,oneof=daily weekly monthly"`
	Range       *DateRange `json:"range,omitempty"`
}

type GetUserFinancialSummaryRequest struct {
	WalletID string     `json:"walletID,omitempty"`
	Range    *DateRange `json:"range,omitempty"`
}
