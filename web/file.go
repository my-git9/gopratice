package web

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type FileUploader struct {
	FileField string
	DstPathFunc func( *multipart.FileHeader ) string
}

func (u FileUploader) Handle() HandleFunc  {
	return func(ctx *Context) {
		// 上传文件的逻辑处理

		// 第一步：读到文件内容
		// 第二步：计算目标路径
		// 第三步：保存文件
		// 第四步：返回结果
		file, header, err := ctx.Req.FormFile("file")
		if err != nil {
			ctx.RespJSON(http.StatusBadRequest, err.Error())
			return
		}
		defer file.Close()
		// 创建目标文件
		dst, err := os.OpenFile(u.DstPathFunc(header),  os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			ctx.RespJSON(http.StatusInternalServerError, err.Error())
			return
		}
		defer dst.Close()

		// 复制数据
		_, err = io.CopyBuffer(dst, file, nil)
		if err != nil {
			ctx.RespJSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.RespJSON(http.StatusOK, "上传成功")
	}
}

// 文件下载
type FileDownloader struct {
	Dir string

}

func (d FileDownloader) Handle() HandleFunc  {
	return func(ctx *Context) {
		// 用的是 xxx?file=xxx.txt
		req, err := ctx.QueryValue("file")
		if err != nil {
			ctx.RespJSON(http.StatusBadRequest, err.Error())
			return
		}
		req = filepath.Clean(req)
		dst := filepath.Join(d.Dir, req)
		// 做一些校验，防止攻击者用相对路径下载其他文件
		//dst, err = filepath.Abs(dst)
		//if string.Contains(dst, d.Dir) {
		//	ctx.RespJSON(http.StatusBadRequest, "文件不存在")
		//	return
		//}
		fn :=  filepath.Base(dst)
		header := ctx.Resp.Header()

		// 核心的请求头
		header.Set("Content-Disposition", "attachment; filename="+fn)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")

		http.ServeFile(ctx.Resp, ctx.Req, dst)
	}
}
