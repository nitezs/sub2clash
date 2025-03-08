package validator

type ShortLinkGenValidator struct {
	Url      string `form:"url" binding:"required"`
	Password string `form:"password"`
	CustomID string `form:"customId"`
}

type GetUrlValidator struct {
	Hash     string `form:"hash" binding:"required"`
	Password string `form:"password"`
}

type ShortLinkUpdateValidator struct {
	Hash     string `form:"hash" binding:"required"`
	Url      string `form:"url" binding:"required"`
	Password string `form:"password" binding:"required"`
}
