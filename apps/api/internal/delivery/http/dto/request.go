package dto

type RegisterRequest struct {
	Username string `json:"username" binding:"required"          example:"john_doe"`
	Password string `json:"password" binding:"required,min=6"    example:"secret123"`
	Role     string `json:"role"                                 example:"user"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"john_doe"`
	Password string `json:"password" binding:"required" example:"secret123"`
}

type AnnouncementRequest struct {
	Title   string `json:"title"   binding:"required" example:"Libur Nasional"`
	Content string `json:"content" binding:"required" example:"Diberitahukan bahwa tanggal 17 Agustus adalah hari libur nasional."`
}
