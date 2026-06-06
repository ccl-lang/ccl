package cclValues

func NewCCLCodeContext() *CCLCodeContext {
	ctx := &CCLCodeContext{}
	ctx.initialize()
	return ctx
}
