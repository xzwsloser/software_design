import sys
import requests
from bs4 import BeautifulSoup
from enum import Enum
from typing import List, Tuple
import pymysql.cursors
from pymysql.connections import  Connection
import re
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
import time


SITE_PRICE_LABEL = 'priceView_box__8FAr4'
SITE_REAL_PRICE_LABEL = 'priceView_real-price-text__xmmuA'

HOT_TAG_CURRENT_LABEL = 'hotTags'
PER_HOT_TAG_LABEL = 'hotTag'

MYSQL_ADDRESS = "localhost"
MYSQL_PORT = 3306
MYSQL_USER = 'root'
MYSQL_PASSWORD = '123456'
MYSQL_DATABASE = "software_design"

class SiteExtraHtmlParser:
    def __init__(self, html_text: str):
        self.parser = BeautifulSoup(html_text, 'html.parser')
    def parse_price(self) -> List[int]:
        # 解析票价
        price_div_tags = self.parser.find_all('div', class_=SITE_PRICE_LABEL)
        result = []
        for price_div in price_div_tags:
            real_price_span = price_div.find('span', class_=SITE_REAL_PRICE_LABEL)
            if real_price_span:
                price_str = real_price_span.text
                try:
                    price = int(price_str[1:])
                except Exception:
                    price = 0
                result.append(price)
            else:
                result.append(0)
        return result

    def parse_comment_count(self) -> Tuple[int,int]:
        positive_count = 0
        negative_count = 0

        hot_tags_span = self.parser.find_all('span', class_=PER_HOT_TAG_LABEL)
        if len(hot_tags_span) == 4:
            positive_comment_str = hot_tags_span[1].text
            negative_comment_str = hot_tags_span[3].text
            print(positive_comment_str)
            print(negative_comment_str)
            m = re.search(r'\d+', positive_comment_str)
            if m:
                positive_count = int(m.group(0))

            m = re.search(r'\d+', negative_comment_str)
            if m:
                negative_count = int(m.group(0))

        return (positive_count, negative_count)

class SiteExtraManager:
    connection: Connection = None

    # positive_count, negative_count, site_idx
    idx_to_price_count: List[Tuple[int,int,int]] = []

    def insert_into_comment_count(self, price_count: Tuple[int,int,int]):
        self.idx_to_price_count.append(price_count)

    def get_price_comment_len(self) -> int:
        return len(self.idx_to_price_count)

    # 存储价格信息
    def _connect_to_mysql(self):
        try:
            self.connection = pymysql.connect(
                host=MYSQL_ADDRESS,
                port=MYSQL_PORT,
                user= MYSQL_USER,
                passwd=MYSQL_PASSWORD,
                database=MYSQL_DATABASE,
                cursorclass=pymysql.cursors.DictCursor
            )
        except Exception as e:
            print(f'Connect to Mysql Failed: {e}')
        
    def store_price_info(self, idx_to_price: List[Tuple[int, int]]):
        if self.connection is None or not self.connection.open:
            self._connect_to_mysql()

        sql = """
        INSERT INTO tb_site_extra
        (site_idx, price) VALUES (%s,%s)
        """
        cursor = self.connection.cursor()

        try:
            pass
            cursor.executemany(sql, idx_to_price)
            self.connection.commit()
            print(f"Successfully Insert {len(idx_to_price)} to mysql")
        except Exception as e:
            print(f'Exception In insert pipeline {e}')
            self.connection.rollback()
        finally:
            cursor.close()
            self.connection.close()
    
    def update_comment_count(self):
        if self.connection is None or not self.connection.open:
            self._connect_to_mysql()
        
        sql = """
        UPDATE `tb_site_extra` 
        SET positive_comment_count = %s,
            negative_comment_count = %s
        WHERE site_idx = %s
        """

        cursor = self.connection.cursor()
        try:
            cursor.executemany(sql, self.idx_to_price_count)
            self.connection.commit()
            print(f"successfully update {len(self.idx_to_price_count)} record")
        except Exception as e:
            print(f"Exception in update process {e}")
            self.connection.rollback()
        finally:
            cursor.close()
            self.connection.close()

        self.idx_to_price_count = []

    def fetch_from_database(self) -> List[str]:
        if not self.connection or not self.connection.open:
            self._connect_to_mysql()

        sql = 'SELECT url from tb_urls WHERE url_type = 0'
        cursor = self.connection.cursor()
        urls = []
        try:
            cursor.execute(sql)
            url_tuples = cursor.fetchall()
            urls = [cur['url'] for cur in url_tuples]
        except Exception as e:
            print(f'Fetch Data Exception: {e}')
        finally:
            cursor.close()
            self.connection.close()

        return urls

def fetch_prices_of_sites():
    site_extra_mananger = SiteExtraManager()
    page_urls = []
    cur_page_url = ''
    for i in range(1, 301):
        cur_page_url = f'https://you.ctrip.com/sight/china110000/s0-p{i}.html'
        page_urls.append(cur_page_url)
    
    site_idx = 1
    for cur_idx, page_url in enumerate(page_urls):
        cur_resp = requests.get(page_url)
        cur_text = cur_resp.text
        parser = SiteExtraHtmlParser(cur_text)
        prices = parser.parse_price()
        idx_to_price = []
        for price in prices:
            idx_to_price.append((site_idx, price))
            site_idx += 1
            if site_idx > 1000:
                break

        print('='*10 + f"{cur_idx + 1}" + "="*10)
        print(idx_to_price)
        site_extra_mananger.store_price_info(idx_to_price)

        if site_idx > 1000:
            break

def fetch_comment_counts(chrome_options, start_idx=1, end_idx=250):
    site_extra_mananger = SiteExtraManager()
    urls = site_extra_mananger.fetch_from_database()

    for cur_idx, url in enumerate(urls):
        site_idx = cur_idx + 1 
        if site_idx < start_idx or site_idx > end_idx:
            break

        # cur_resp = requests.get(url) 
        # cur_text = cur_resp.text

        driver = webdriver.Chrome(options=chrome_options)
        driver.get(url)
        time.sleep(1)
        cur_text = driver.page_source
        driver.quit()

        parser = SiteExtraHtmlParser(cur_text)
        positive_comment, negative_comment = parser.parse_comment_count()
        site_extra_mananger.insert_into_comment_count((positive_comment, negative_comment, site_idx))
        print(f"="*10 + f"{site_idx}" + "="*10)
        print(url)
        print(f"p: {positive_comment} n: {negative_comment}")

        if site_extra_mananger.get_price_comment_len() >= 25 or site_idx == end_idx:
            site_extra_mananger.update_comment_count()

if __name__ == '__main__':
    # 模拟浏览器后台运行
    chrome_options = Options()
    chrome_options.add_argument("--headless")
    chrome_options.add_argument("--disable-gpu")

    start_idx = 1
    end_idx = 250
    if len(sys.argv) >= 3:
        start_idx = int(sys.argv[1])
        end_idx = int(sys.argv[2])

    fetch_comment_counts(chrome_options=chrome_options, start_idx=start_idx, end_idx=end_idx)
