package file

import (
	"net/http"

	"github.com/yeisme/notevault/internal/logic/file"
	"github.com/yeisme/notevault/internal/svc"
	"github.com/yeisme/notevault/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// (高级) 获取文件两个版本之间的差异信息 (主要针对文本文件)。
func GetFileVersionDiffHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileVersionDiffRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := file.NewGetFileVersionDiffLogic(r.Context(), svcCtx)
		resp, err := l.GetFileVersionDiff(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
