package file

import (
	"net/http"

	"github.com/yeisme/notevault/internal/logic/file"
	"github.com/yeisme/notevault/internal/svc"
	"github.com/yeisme/notevault/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 将文件恢复到特定版本。
func RevertFileVersionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RevertFileVersionRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := file.NewRevertFileVersionLogic(r.Context(), svcCtx)
		resp, err := l.RevertFileVersion(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
