package file

import (
	"net/http"

	"github.com/yeisme/notevault/internal/logic/file"
	"github.com/yeisme/notevault/internal/svc"
	"github.com/yeisme/notevault/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DownloadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileDownloadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		l := file.NewDownloadFileLogic(r.Context(), svcCtx, r, w)
		err := l.DownloadFile(&req)
		if err != nil {
			return
		}
	}
}
