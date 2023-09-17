package validator

type ShortLinkGenValidator struct {
	Url string `form:"url" binding:"required"`
}

type ShortLinkGetValidator struct {
	Hash string `form:"hash" binding:"required"`
}
