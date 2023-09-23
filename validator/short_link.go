package validator

type ShortLinkGenValidator struct {
	Url      string `form:"url" binding:"required"`
	Password string `form:"password"`
}

type ShortLinkGetValidator struct {
	Hash     string `form:"hash" binding:"required"`
	Password string `form:"password"`
}
