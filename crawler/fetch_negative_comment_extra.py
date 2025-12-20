import requests
import pymysql
import pymysql.cursors
from pymysql.connections import Connection
from typing import List, Tuple
from bs4 import BeautifulSoup
import json

MYSQL_ADDRESS = "localhost"
MYSQL_PORT = 3306
MYSQL_USER = 'root'
MYSQL_PASSWORD = '123456'
MYSQL_DATABASE = "software_design"

MAX_PAGE_SIZE = 50
FETCH_COMMENT_URL = 'https://m.ctrip.com/restapi/soa2/13444/json/getCommentCollapseList'
SITE_NEXT_JS_ID = '__NEXT_DATA__'


EACH_INSERT_SIZE = 250



class DbManager:
    conn: Connection = None
    comment_infos: List[Tuple[int, str, str, str]] = []

    def _init_conn(self):
        self.conn = pymysql.connect(
            port=MYSQL_PORT,
            user=MYSQL_USER,
            host=MYSQL_ADDRESS,
            passwd=MYSQL_PASSWORD,
            database=MYSQL_DATABASE,
            cursorclass=pymysql.cursors.DictCursor
        )

    def fetch_site_urls(self) -> List[str]:
        if self.conn is None or not self.conn.open:
            self._init_conn()

        sql = 'SELECT url FROM tb_urls WHERE url_type = 0'
        cursors = self.conn.cursor()
        urls = []
        try:
            cursors.execute(sql)
            urls_obj = cursors.fetchall()
            urls = [url_obj['url'] for url_obj in urls_obj]
        except Exception as e:
            print(f"fetch url failed: {e}")
        finally:
            cursors.close()
            self.conn.close()

        return urls   
    
    def insert_comments(self, *new_comments_info):
        self.comment_infos.extend(new_comments_info)
    
    def get_number_of_comments(self):
        return len(self.comment_infos)
    
    def store_comments(self):
        if self.conn is None or not self.conn.open:
            self._init_conn()

        sql =\
            """
            INSERT INTO `tb_negative_comment_extra` 
            (site_idx, content, tourist_type, ip)  
            VALUES
            (%s, %s, %s, %s)
            """

        cursor = self.conn.cursor()

        try:
            cursor.executemany(sql, self.comment_infos)
            self.conn.commit()
            print(f"Successfully Insert {len(self.comment_infos)} rows!")
        except Exception as e:
            self.conn.rollback()
            print(f"Exception: {e}")
        finally:
            cursor.close()
            self.conn.close()

        self.comment_infos = []

class SimpleHtmlParser:
    def __init__(self, html_text: str):
        self.parser = BeautifulSoup(html_text, 'html.parser')
    def parser_poi_id(self) -> int:
        # 解析得到  poiId
        script_tag = self.parser.find('script', id=SITE_NEXT_JS_ID)
        poiId = -1
        if script_tag:
            json_data = json.loads(script_tag.text)
            poiId = json_data.get('props').get('pageProps').get('initialState').get('poiDetail').get('poiId')

        return poiId

def fetch_negative_comments(comment_start_idx: int = 0):
    db = DbManager()

    # 1. get the urls
    urls = db.fetch_site_urls()

    print(f"Successfully get {len(urls)} rows from tb_urls")

    for cur_idx, url in enumerate(urls):
        # cur_site_idx = cur_idx + 1 = 90 已经插入
        # cur_idx = 89
        # cur_idx <= 89 已经插入
        # cur_idx < 90
        if cur_idx < comment_start_idx:
            continue

        cur_site_idx = cur_idx + 1
        cur_site_resp = requests.get(url)
        parser = SimpleHtmlParser(cur_site_resp.text)
        cur_poi_id = parser.parser_poi_id()
        cur_site_comments = []
        if not cur_poi_id == -1:
            payload_negative = {
                "arg": {
                    "channelType": 2,
                    "collapseType": 0,
                    "commentTagId": -12,
                    "pageIndex": 0,
                    "pageSize": MAX_PAGE_SIZE,
                    "poiId": cur_poi_id,
                    "sourceType": 1,
                    "sortType": 3,
                    "starType": 0
                },
                "head": {
                    "cid": "09031179414218738290",
                    "ctok": "",
                    "cver": "1.0",
                    "lang": "01",
                    "sid": "8888",
                    "syscode": "09",
                    "auth": "",
                    "xsid": "",
                    "extension": []
                }
            }

            # 模拟浏览器发送请求
            headers = {
                    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
                    "Content-Type": "application/json"
            }

            resp = requests.post(FETCH_COMMENT_URL, json=payload_negative, headers=headers)

            items = resp.json()['result']['items']

            if items:
                for comment in items:
                    content = comment.get('content', '') or ''
                    tourist_type = comment.get('touristTypeDisplay', '') or ''
                    ip_location = comment.get('ipLocatedName', '') or ''
                    cur_site_comments.append((content,tourist_type, ip_location))

        cur_site_negative_comments = [(cur_site_idx, comment_info[0], comment_info[1], comment_info[2]) for comment_info in cur_site_comments]


        print('-'*10 + f'{cur_site_idx}' + '-'*10)
        if len(cur_site_negative_comments) >= 1:
            print(cur_site_negative_comments[0][1])
            print(cur_site_negative_comments[0][2])
            print(cur_site_negative_comments[0][3])
        else:
            print("empty")
        print('-'*21)

        db.insert_comments(*cur_site_negative_comments)

        if db.get_number_of_comments() >= EACH_INSERT_SIZE or cur_site_idx >= len(urls):
            db.store_comments()


if __name__ == '__main__':
    fetch_negative_comments()

