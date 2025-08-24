package web

import (
	"github.com/stretchr/testify/require"
	"html/template"
	"log"
	"mime/multipart"
	"path/filepath"
	"testing"
)

func TestUpload(t *testing.T)  {
	tpl, err := template.ParseGlob("testdata/tpls/*.gohtml")
	require.NoError(t, err)
	engine := &GoTemplateEngine{
		T: tpl,
	}

	h := NewHTTPServer(ServerWithTemplateEngine(engine))
	h.Get("/upload", func(ctx *Context) {
		err := ctx.Render("upload.gohtml", nil)
		if err != nil {
			log.Println(err)
		}
	})

	fu := FileUploader{
		FileField: "myfile",
		DstPathFunc: func(header *multipart.FileHeader) string {
			return filepath.Join("testdata", header.Filename)
		},
	}

	h.Post("/upload", fu.Handle())

	h.Start(":8080")
}

func TestDownload(t *testing.T) {
	h := NewHTTPServer()

	fu := FileDownloader{
		Dir: "testdata/test.txt",
	}

	h.Get("/download", fu.Handle())
	h.Start(":8080")
}
