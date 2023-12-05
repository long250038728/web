package str

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func read() {
	/**
	1.string 保存string.reader类中，参数为s字符串，i为0，prevRune为-1
	2.当读取io.Reader，循环调用string.reader的read方法直到返回长度0(读不到了),把[]byte传入，
	3.read方法取 reader.s[0:([]byte长度)] ,然后执行copy方法把数据拷贝到[]byte中，同时返回长度给外部判断是否退出循环
	4.直到退出循环就能拿到完整的值赋值给[]byte中
	*/
	reader := strings.NewReader("hello")
	b, err := io.ReadAll(reader)
	fmt.Println(string(b))
	fmt.Println(err)

	/**
	字符串拼接的三个效率较高的的处理方式
	*/
	bu := bytes.NewBufferString("hello")
	bu.WriteString(" world")
	fmt.Println(bu.String())

	build := strings.Builder{}
	build.WriteString("hello")
	build.WriteString(" world2")
	fmt.Println(build.String())

	str := "hello world3"
	byt := make([]byte, 0, 100)
	byt = append(byt, str...)
	fmt.Println(string(byt))
}
