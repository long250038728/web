### text/template

可以自定义生成相关的文件或文档

#### 基本操作
{{ .Name}}          传入的变量.Name
{{ $var="hello" }}  定义变量$var
{{- .Name}}         空格去除

#### 遍历 (支持数组，切片，map，channel)
{{range .List}} {{.}}  {{end}}                          遍历传入的变量.List   {{.}}代表item
{{range .List}} {{.Name}}  {{end}}                      遍历传入的变量.List   {{.Name}}代表item.Name
{{range $index,$item := .List}} {{$index}}  {{end}}     遍历传入的变量.List   把index及item用变量$index,$item接收，后面方便获取


#### 函数
{{len .List}}  
{{myFunc .List}}  
```
系统：
len：返回一个字符串、数组、切片、映射或通道的长度。
index：返回一个字符串、数组或切片中指定位置的元素。
printf：根据格式字符串输出格式化的字符串。
range：遍历一个数组、切片、映射或通道，并输出其中的每个元素。
with：设置当前上下文中的变量。

自定义
funcMap :=template.FuncMap{
    "myFunc": myFunc,
}
template.New("name").Funcs(funcMap)
```


#### 判断
{{if eq .Name "sam"}}Hello, {{.Name}} {{end}}                如果.Name == "sam" 显示中间的内容
{{if eq .Name "sam"}} sam {{else}} other {{end}}           如果.Name == "sam" 显示sam 否知显示other

```
eq：当等式 arg1 == arg2 成立时，返回 true，否则返回 false
ne：当不等式 arg1 != arg2 成立时，返回 true，否则返回 false
lt：当不等式 arg1 < arg2 成立时，返回 true，否则返回 false
le：当不等式 arg1 <= arg2 成立时，返回 true，否则返回 false
gt：当不等式 arg1 > arg2 成立时，返回 true，否则返回 false
ge：当不等式 arg1 >= arg2 成立时，返回 true，否则返回 false
```


