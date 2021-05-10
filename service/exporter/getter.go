package exporter

import "read-test-server/dao"

type UserGetter struct{}

func (u *UserGetter) Table() string {
	return "user"
}

func (u *UserGetter) Getter(ctx GetterCtx) (data interface{}, err error) {
	usersByIds, err := dao.QueryUsersByIds(ctx.GetIds())
	return usersByIds, err
}

type AnswerGetter struct{}

func (u *AnswerGetter) Table() string {
	return "answer"
}

func (u *AnswerGetter) Getter(ctx GetterCtx) (data interface{}, err error) {
	usersByIds, err := dao.QueryAnswersByUidsAndPaperId(ctx.GetIds(), ctx.GetPaperId())
	return usersByIds, err
}