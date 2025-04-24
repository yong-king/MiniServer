package sequence

import (
	"database/sql"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const sqlReplaceIntoSub = `REPLACE INTO sequence (stub) VALUES ('a')`

type MySQL struct{
	conn sqlx.SqlConn
}

func NewMysql(dns string) Sequence {
	return &MySQL{
		conn: sqlx.NewMysql(dns),
	}
}

func (m *MySQL) Next() (seq uint64, err error){
	// prepare 预编译
	var stmt sqlx.StmtSession
	stmt, err = m.conn.Prepare(sqlReplaceIntoSub)
	if err != nil {
		logx.Errorw("conn.Prepare failed", logx.LogField{Key: "err", Value: err.Error()})
		return 0, err
	}
	defer stmt.Close()

	// 执行
	var rest sql.Result
	rest, err = stmt.Exec()
	if err != nil {
		logx.Errorw("stmt.Exec failed", logx.LogField{Key: "err", Value: err.Error()})
		return 0, err
	}

	// 获取刚刚插入的主键id
	var lid int64
	lid, err = rest.LastInsertId()
	if err != nil {
		logx.Errorw("rest.LastInsertId failed", logx.LogField{Key: "err", Value: err.Error()})
		return 0, err
	}
	seq = uint64(lid)
	return
}