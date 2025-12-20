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

SITE_NEXT_JS_ID = '__NEXT_DATA__'

FETCH_COMMENT_URL = 'https://m.ctrip.com/restapi/soa2/13444/json/getCommentCollapseList'

class PositiveExtraManager:
    conn: Connection = None
    comment_infos: List[Tuple[int, str,str, str]] = []
    def _init_conn(self):
        self.conn = pymysql.connect(
            port=MYSQL_PORT,
            user=MYSQL_USER,
            host=MYSQL_ADDRESS,
            passwd=MYSQL_PASSWORD,
            database=MYSQL_DATABASE,
            cursorclass=pymysql.cursors.DictCursor
        )
    def insert_to_command_info(self,*record):
        self.comment_infos.extend(record)
    def get_number_of_command_info(self) -> int:
        return len(self.comment_infos)
    
    def store_into_mysql(self):
        if self.conn is None or not self.conn.open:
            self._init_conn()
        
        sql = """
        insert into tb_positive_comment_extra2
        (site_idx,content, tourist_type, ip) values 
        (%s,%s,%s,%s)
        """

        cursor = self.conn.cursor()
        try:
            cursor.executemany(sql, self.comment_infos)
            self.conn.commit()
            print(f"successfully insert {len(self.comment_infos)} into mysql")
        except Exception as e:
            self.conn.rollback()
            print(f"Exception insert into comment extra: {e}")
        finally:
            cursor.close()
            self.conn.close()
        
        self.comment_infos = []

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

class PoiIDParser:
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

def parse_cur_comments(cur_poi_id) -> List[Tuple[str,str, str]]:
        cur_site_comments = []
        if not cur_poi_id == -1:
            for idx in range(2):
                payload = {
                    "arg": {
                        "channelType": 2,
                        "collapseType": 0,
                        "commentTagId": 0,
                        "pageIndex": idx,
                        "pageSize": 50,
                        "poiId": cur_poi_id,
                        "sourceType": 1,
                        "sortType": 3,
                        "starType": 0
                    },
                    "head": {
                        "cid": "",
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

                resp = requests.post(FETCH_COMMENT_URL, json=payload, headers=headers)

                items = resp.json()['result']['items']

                if items:
                    for comment in items:
                        content = comment.get('content', '') or ''
                        tourist_type = comment.get('touristTypeDisplay', '') or ''
                        ip_location = comment.get('ipLocatedName', '') or ''
                        cur_site_comments.append((content,tourist_type, ip_location))

        return cur_site_comments

def fetch_positive_comment_extra(start_index=1):
    db = PositiveExtraManager()
    urls = db.fetch_site_urls()
    for cur_idx, url in enumerate(urls):
        site_idx = cur_idx + 1
        if site_idx < start_index:
            continue
        resp = requests.get(url)
        text = resp.text
        parser = PoiIDParser(text)
        poi_id = parser.parser_poi_id()

        # 发送请求
        cur_site_comments = parse_cur_comments(poi_id)
        cur_site_comments_info = []

        for comment_info in cur_site_comments:
            cur_site_comments_info.append((site_idx, comment_info[0], comment_info[1], comment_info[2]))

        print(f'='*10 + f"{cur_idx + 1}" + '='*10)
        if len(cur_site_comments_info) != 0:
            print(cur_site_comments_info[0][1])
            print(cur_site_comments_info[0][2])
            print(cur_site_comments_info[0][3])
        else:
            print('empty')
            print('='*20)

        db.insert_to_command_info(*cur_site_comments_info)

        if db.get_number_of_command_info() >= 1000 or cur_idx == 999:
            db.store_into_mysql()

if __name__ == '__main__':
    fetch_positive_comment_extra(start_index=132)
