package file_demo

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	f, err := os.Open("testdata/my_file.txt")
	require.NoError(t, err)
	data := make([]byte, 64)
	n, err := f.Read(data)
	fmt.Println(n)
	require.NoError(t, err)
	f.Close()

	// 打开一个 append only 模式的文件
	f, err = os.OpenFile("testdata/my_file.txt", os.O_APPEND, fs.ModeAppend)
	require.NoError(t, err)
	n, err = f.WriteString("hello world")
	fmt.Println(n)
	require.NoError(t, err)
	f.Close()

	f, err = os.Create("testdata/my_file2.txt")
	require.NoError(t, err)
	f.WriteString("hello world")
	f.Close()
}
