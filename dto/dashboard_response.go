package dto

type DashboardResponse struct {
	TotalUsers     int64 `json:"total_users"`
	TotalPoli      int64 `json:"total_poli"`
	TotalWaiting   int64 `json:"total_waiting"`
	TotalCalled    int64 `json:"total_called"`
	TotalCompleted int64 `json:"total_completed"`
}
