# pip install selenium

# == mac ==
# sudo vi ~/.zshrc
# alias google-chrome="/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome"

# == 驱动 ==
# https://googlechromelabs.github.io/chrome-for-testing/ 下载对应驱动解压安装


# 切换虚拟环境
# uv venv
# source .venv/bin/activate
# uv add selenium
# python3 main.py

import json
import random

from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service

from selenium.webdriver.support.wait import WebDriverWait
from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC


def init_driver() -> webdriver.Chrome:
    chromedriver_path = "/Users/linlong/Downloads/chromedriver-mac-x64/chromedriver"
    # '--verbose',  log_output=sys.stdout,
    service = Service(executable_path=chromedriver_path,
                      service_args=['--headless=new', '--no-sandbox',
                                    '--disable-dev-shm-usage',
                                    '--disable-gpu',
                                    '--ignore-certificate-errors',
                                    '--ignore-ssl-errors',
                                    ])

    options = Options()
    options.add_argument('--disable-gpu')  # 禁用GPU渲染
    options.add_argument('--incognito')  # 无痕模式
    options.add_argument('--ignore-certificate-errors-spki-list')  # 忽略与证书相关的错误
    options.add_argument('--disable-notifications')  # 禁用浏览器通知和推送API
    options.add_argument('--disable-extensions')  # 禁用浏览器扩展
    options.add_argument('--start-maximized')
    # options.add_argument(f'user-agent={get_UA()}')   # 修改用户代理信息
    # options.add_argument('--window-name=huya_test')  # 设置初始窗口用户标题
    # options.add_argument('--window-workspace=1')  # 指定初始窗口工作区  #
    # options.add_argument('--force-dark-mode')  # 使用暗模式
    # options.add_argument('--start-fullscreen')  # 指定浏览器是否以全屏模式启，与进入浏览器后按F11效果相同
    return webdriver.Chrome(options=options, service=service)


def get_UA():
    UA_list = [
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.54 Safari/537.36',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4651.0 Safari/537.36',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.0.0 Safari/537.36'
    ]
    UA = UA_list[random.randint(0, len(UA_list) - 1)]
    return UA


def waitElement(d: webdriver.Chrome, xpath: str):
    return WebDriverWait(d, 10).until(EC.element_to_be_clickable(By.XPATH, xpath))


def waitElements(d: webdriver.Chrome, xpath: str):
    return WebDriverWait(d, 10).until(EC.presence_of_all_elements_located(By.XPATH, xpath))


job_tpl = """
内容: {}
状态: {}
类别: {}
方式: {}
时间: {}
操作人: {}
     """


def main():
    driver = init_driver()
    try:
        driver.get("https://test.zhubaoe.cn/erp.php/login/index.html?code=sz208356")

        username = waitElement(driver, f'//*[@id="username"]')
        password = waitElement(driver, f'//*[@id="app"]/div[1]/div/div[1]/div/div/div[2]/div/form/div[2]/div/input')
        summit = waitElement(driver, f'//*[@id="app"]/div[1]/div/div[1]/div/div/div[2]/div/form/button')

        username.send_keys("大戈")
        password.send_keys("654321")
        summit.click()

        crm = waitElement(driver, f'//*[@id="site-navbar-collapse"]/ul[1]/li[3]')
        crm.click()

        messageMenu = waitElement(driver, f'/html/body/div[1]/div/div[1]/div/ul/li[4]/a/span[1]')
        messageMenu.click()

        messageListMenu = waitElement(driver, f'/html/body/div[1]/div/div[1]/div/ul/li[4]/ul/li[1]/a/span')
        messageListMenu.click()

        arr = []

        waitElement(driver, f'//*[@id="record_body"]/tr[1]')

        for i in waitElements(driver, f'//*[@id="record_body"]/tr'):
            arr.append({
                "context": waitElement(i, f'td[1]').text,
                "status": waitElement(i, f'td[2]').text,
                "type": waitElement(i, f'td[3]').text,
                "sys": waitElement(i, f'td[4]').text,
                "time": waitElement(i, f'td[5]').text,
                "admin": waitElement(i, f'td[6]').text,
                "txt": job_tpl.format(
                    waitElement(i, f'td[1]').text,
                    waitElement(i, f'td[2]').text,
                    waitElement(i, f'td[3]').text,
                    waitElement(i, f'td[4]').text,
                    waitElement(i, f'td[5]').text,
                    waitElement(i, f'td[6]').text,
                )
            })
        driver.save_screenshot("page_screenshot.png")

        with open("data.json", "w") as w:
            w.write(json.dumps(arr, ensure_ascii=False))

    except Exception as e:
        print(e)

    finally:
        driver.close()


if __name__ == '__main__':
    main()
