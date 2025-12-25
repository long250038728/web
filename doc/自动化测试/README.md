## playwright 与 Selenium对比
* Selenium 生态庞大但是功能陈旧
  * 需要等待元素出现后才能使用
  * 支持素有主流浏览器
  * 需要下载对应的浏览器驱动等
* playwright 速度快开箱即用
  * 主要支持chromium，firefox，webKit等web现代应用
  * 可通过pip一键安装
  * 可以代码生成不用写代码
  * 支持mcp

---
## Selenium
安装
```shell
pip install selenium

# == mac ==
sudo vi ~/.zshrc
alias google-chrome="/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome"

# == 驱动 ==
 https://googlechromelabs.github.io/chrome-for-testing/ 下载对应驱动解压安装
```

代码编写

[selenium代码](p_selenium.py)

代码运行
```shell
python p_selenium.py
```


---

## playwright
安装
```shell
# 安装 playwright
pip install playwright
# 安装 playwright驱动
playwright install
```

代码自动生成
```shell
# -o 文件输出 -b 使用xxx浏览器  
playwright codegen -o p_playwright.py -b cr "https://www.baidu.com"
```

代码运行
```shell
python p_playwright.py
```

mcp
```json
{
  "mcpServers": {
    "playwright": {
      "command": "npx",
      "args": [
        "@playwright/mcp@latest"
      ]
    }
  }
}
```