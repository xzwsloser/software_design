import pymysql
import pandas as pd
import csv

MYSQL_HOST      = 'localhost'
MYSQL_PORT      = 3306
MYSQL_USER      = 'root'
MYSQL_PASSWD    = '123456'
MYSQL_DATABASE  = 'software_design'

def transfrom_mysql_to_csv():
    conn = pymysql.connect(
        host=MYSQL_HOST,
        port=MYSQL_PORT,
        user=MYSQL_USER,
        passwd=MYSQL_PASSWD,
        database=MYSQL_DATABASE
    )

    tables_to_read = ['tb_sites', 
                      'tb_site_extra', 
                      'tb_positive_comment_extra2', 
                      'tb_negative_comment_extra']
    
    for table_name in tables_to_read:
        sql = f"SELECT * FROM {table_name}"
        df = pd.read_sql(sql, conn)
        df.to_csv(
            f'../dataset/{table_name}.csv',
            index=False,
            encoding='utf-8-sig',
            quoting=csv.QUOTE_NONNUMERIC,
            escapechar='\\' # 转义字符
        )
        print(f"successfully save {table_name}.csv")

    conn.close()

if __name__ == "__main__":
    transfrom_mysql_to_csv()
