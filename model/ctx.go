package model

type GetterCtx struct {
	Table   string `json:"table"`
	Ext     string `json:"ext"`
	Ids     []uint `json:"ids"`
	PaperId uint   `json:"paper_id"`
}

func (ctx *GetterCtx) GetIds() []uint {
	return ctx.Ids
}

func (ctx *GetterCtx) GetPaperId() uint {
	return ctx.PaperId
}
