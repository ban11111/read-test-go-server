package model

type GetterCtx struct {
	Ids []uint `json:"ids"`
}

func (ctx *GetterCtx) GetIds() []uint {
	return ctx.Ids
}