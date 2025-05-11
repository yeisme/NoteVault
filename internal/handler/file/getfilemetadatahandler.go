package file

import (
	"net/http"

	"github.com/yeisme/notevault/internal/logic/file"
	"github.com/yeisme/notevault/internal/svc"
	"github.com/yeisme/notevault/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取特定文件的元数据。可选获取特定版本的元数据。
func GetFileMetadataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetFileMetadataRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := file.NewGetFileMetadataLogic(r.Context(), svcCtx)
		resp, err := l.GetFileMetadata(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
