package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SequenceModel = (*customSequenceModel)(nil)

type (
	// SequenceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSequenceModel.
	SequenceModel interface {
		sequenceModel
	}

	customSequenceModel struct {
		*defaultSequenceModel
	}
)

// NewSequenceModel returns a model for the database table.
func NewSequenceModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SequenceModel {
	return &customSequenceModel{
		defaultSequenceModel: newSequenceModel(conn, c, opts...),
	}
}
