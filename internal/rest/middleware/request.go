package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ReqIdKey string

const reqIdKey = ReqIdKey("reqid")

func UseRequestId() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqId, _ := uuid.NewRandom()
			ctx := r.Context()
			ctx = context.WithValue(ctx, reqIdKey, reqId.String())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetRequestId(r *http.Request) string {
	return GetReqIdCtx(r.Context())
}

func GetReqIdCtx(ctx context.Context) string {
	reqId := ctx.Value(reqIdKey).(string)
	return reqId
}
