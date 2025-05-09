# https://akshare.akfamily.xyz/introduction.html
# uv add akshare --upgrade

import akshare as ak
import pandas as pd


def main(symbol :str ,start_date:str ,end_date:str):
    df = ak.stock_zh_a_hist(
        symbol=symbol,
        period="daily", # choice of {'daily', 'weekly', 'monthly'}
        adjust="qfq", # qfq: 返回前复权后的数据; hfq: 返回后复权后的数据
        start_date = start_date,
        end_date = end_date,
    )

    if len(df) == 0:
        return

    # pd常见操作
    print(df.head())  # 前5行
    print(df.tail(3))  # 后3行
    print(df.info())  # 数据概况
    print(df.describe())  # 数值型统计摘要
    print(df.shape)  # 行列数

    print(df["日期"])
    print(df[["日期","开盘"]])
    print(df[df["开盘"] >= 47])


    # df的日期列改为datetime格式
    df["日期"] = pd.to_datetime(df["日期"])

    # 按照日期排序为索引 inplace=True在原始数据修改
    df.set_index("日期", inplace=True)
    df.sort_index(ascending = False,inplace=True)

    # # index=False 为索引是否显示
    df.to_csv( symbol +'.csv')


if __name__ == "__main__":
    main(symbol="002230",start_date="20250401",end_date="20250508")
