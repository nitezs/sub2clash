package validator

type ShortLinkGenValidator struct {
	Url      string `form:"url" binding:"required"`
	Password string `form:"password"`
}

type ShortLinkGetValidator struct {
	Hash     string `form:"hash" binding:"required"` // Hash: 短链接
	Password string `form:"password"`
}

type ShortLinkUpdateValidator struct {
	Hash     string `form:"hash" binding:"required"`
	Url      string `form:"url" binding:"required"`
	Password string `form:"password" binding:"required"`
}
