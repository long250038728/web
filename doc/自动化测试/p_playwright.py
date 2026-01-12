import re
from playwright.sync_api import Playwright, sync_playwright, expect


def run(playwright: Playwright) -> None:
    browser = playwright.chromium.launch(headless=False)
    context = browser.new_context()
    page = context.new_page()
    page.goto("https://moss.zhubaoe.cn/erp.php/login/index.html?code=sz568302")
    page.get_by_role("textbox", name="登录账号").click()
    page.get_by_role("textbox", name="登录账号").fill("林龙")
    page.get_by_role("textbox", name="密码").click()
    page.get_by_role("textbox", name="密码").fill("123456")
    page.get_by_role("textbox", name="密码").press("Enter")
    page.get_by_role("link", name="SCRM").click()
    page.get_by_role("button", name="查询", exact=True).click()
    with page.expect_popup() as page1_info:
        page.locator(".el-table__fixed-right > .el-table__fixed-body-wrapper > .el-table__body > tbody > .el-table__row.hover-row > .el-table_1_column_18 > .cell > div > span").first.click()
    page1 = page1_info.value
    page1.get_by_role("tab", name="订单记录").click()
    page1.get_by_text("新品销退单 历史销售单 维修单 打金单据 所属门店业务类型订单编号成交日期商品条码商品名称数量净金重标签价销售金价回收价格成交金价工费金额应售金额活动总折扣率活").click()
    page1.locator("label").filter(has_text="历史销售单").click()
    page1.locator("label").filter(has_text="维修单").click()
    page1.get_by_role("tab", name="跟进记录").click()
    page1.get_by_role("tab", name="优惠券记录").click()
    page1.get_by_role("tab", name="积分记录").click()
    page1.get_by_role("tab", name="优惠券记录").click()
    page1.get_by_role("link", name="积分活动").click()
    page1.get_by_role("button", name="查询").click()
    page1.close()
    page.close()

    # ---------------------
    context.close()
    browser.close()


with sync_playwright() as playwright:
    run(playwright)
