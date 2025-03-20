package ctx

import (
	"context"
	"gorm.io/gorm"
	"net/http"
)

const TxKey = "tx"

//func SetTx(db *gorm.DB) context.Context {
//	return context.WithValue(context.Background(), txKey, db)
//}

//func GetTx() *gorm.DB {
//	ctxBackground := context.Background()
//	val := ctxBackground.Value(txKey)
//	return val.(*gorm.DB)
//}

func SetTx(r *http.Request, db *gorm.DB) *http.Request {
	//r.WithContext(context.WithValue(r.Context(), TxKey, db))
	ctx := context.WithValue(r.Context(), TxKey, db)
	return r.WithContext(ctx)
}

func GetTx(r *http.Request) *gorm.DB {
	return r.Context().Value(TxKey).(*gorm.DB)
}
