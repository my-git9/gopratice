package template_demo

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"html/template"
	"strings"
	"testing"
)

func TestHelloWorld(t *testing.T)  {
	type User struct {
		Name string
	}
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
切片长度: {{len .Slice}}
say hello: {{.Hello "Tom" "Jerry"}}
打印数字: {{printf "%.2f" 1.234}}
`)

	require.NoError(t, err)

	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, FuncCall{})

	require.NoError(t, err)
	assert.Equal(t, `
切片长度: 2
say hello: Tom.Jerry
打印数字: 1.23
`, buffer.String())
}

func TestForLoop(t *testing.T)  {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
{{- range $idx, $ele := .}}
{{- $idx}}
{{- end}}
`)

	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, make([]int, 100))
	require.NoError(t, err)
	assert.Equal(t, strings.Repeat("0-", 100), buffer.String())

}

type FuncCall struct {
	Slice []string
}

func (f FuncCall) Hello(first string, last string) string {
	return fmt.Sprintf("%s.%s", first, last)
}
