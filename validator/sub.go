package validator

type SubQuery struct {
	Sub     string `form:"sub" json:"name" binding:"required"`
	Mix     bool   `form:"mix,default=false" json:"email" binding:""`
	Refresh bool   `form:"refresh,default=false" json:"age" binding:""`
}
