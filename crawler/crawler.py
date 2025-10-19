import requests
from bs4 import BeautifulSoup
from enum import Enum
from typing import List, Tuple
import pymysql.cursors
from pymysql.connections import  Connection
import re
import json


# @Description: 特殊 Url, 类名定义
INIT_URL = 'https://you.ctrip.com/sight/china110000.html'
FETCH_COMMENT_URL = 'https://m.ctrip.com/restapi/soa2/13444/json/getCommentCollapseList'

MAX_URL_NUMBER = 1000
MAX_SITES_IN_MEMORY = 50
MAX_PAGE_SIZE = 50
MAX_SITES_COMMENT_IN_MEMORY = 1000

SITE_URL_CLASS_NAME = 'titleModule_name__Li4Tv'
SITE_INFO_CLASS_NAME = 'baseInfoModule'
SITE_IMAGE_INFO_CLASS_NAME = 'swiper'
SITE_TEXT_INFO_CLASS_NAME = 'baseInfoMain'
SITE_TITLE_VIEW_CLASS_NAME = 'titleView'
SITE_TITLE_CLASS_NAME = 'title'
SITE_COMMENT_CLASS_NAME = 'comment'
SITE_COMMENT_SCORE_CLASS_NAME = 'commentScore'
SITE_COMMENT_NUM_CLASS_NAME = 'commentScoreNum'
SITE_BASE_INFO_CLASS_NAME = 'baseInfoContent'
SITE_BASE_INFO_CONTEXT_CLASS_NAME = 'baseInfoItem'
SITE_BASE_INFO_CONTEXT_TEXT_CLASS_NAME = 'baseInfoText'
SITE_OPEN_TIME_CLASS_NAME = 'baseInfoText cursor openTimeText'
SITE_PHONE_SPAN_CLASS_NAME = 'phoneHeaderItem'
SITE_HOT_DEGREE_CLASS_NAME = 'heatScoreText'
SITE_IMAGE_CLASS_NAME = 'swiper-slide'
SITE_INTRO_CLASS_NAME = 'LimitHeightText'

SITE_NEXT_JS_ID = '__NEXT_DATA__'

MYSQL_ADDRESS = "localhost"
MYSQL_PORT = 3306
MYSQL_USER = 'root'
MYSQL_PASSWORD = '123456'
MYSQL_DATABASE = "software_design"

class Site:
     name: str
     score: float
     address: str
     hot_degree: float
     introduce: str
     open_time: str
     phone: str
     images: List[str]

     def __init__(self,
                  name: str,
                  score: float,
                  address: str,
                  hot_degree: float,
                  introduce: str,
                  open_time: str,
                  phone: str,
                  images: List[str]
                  ):
         self.name = name
         self.score = score
         self.address = address
         self.hot_degree = hot_degree
         self.introduce = introduce
         self.open_time = open_time
         self.phone = phone
         self.images = images

     def __str__(self):
         info =\
         f'''
         名称: {self.name}
         分数: {self.score}
         地址: {self.address}
         热度: {self.hot_degree}
         开放时间: {self.open_time}
         官方电话: {self.phone}
         '''

         image_info = '\n图片列表:\n'
         for image_url in self.images:
             image_info += image_url
             image_info += '\n'

         intro_info = f'\n介绍:\n{self.introduce}'

         return info + image_info + intro_info

# @Description: Url 存储以及持久化(只持久化存储 Url)
class UrlManager:
    database: Connection = None
    def __init__(self):
        self.site_urls = []
    def insert_into_urls(self, *urls):
        self.site_urls.extend(urls)
    def _connect_to_mysql(self):
        try:
            self.database = pymysql.connect(
                host=MYSQL_ADDRESS,
                port=MYSQL_PORT,
                user= MYSQL_USER,
                passwd=MYSQL_PASSWORD,
                database=MYSQL_DATABASE,
                cursorclass=pymysql.cursors.DictCursor
            )
        except Exception as e:
            print(f'Exception: {e}')

    def store_into_database(self):
        if not self.database or not self.database.open:
            self._connect_to_mysql()
        sql = 'INSERT INTO tb_urls (url, url_type) VALUES (%s, %s)'
        data_lists = [(url, 0) for url in self.site_urls]
        cursor = self.database.cursor()
        try:
            cursor.executemany(sql, data_lists)
            self.database.commit()
            print('Site Url Store Successfully!')
        except Exception as e:
            print(f'Site Url Store Exception {e}')
            self.database.rollback()
        finally:
            cursor.close()
            self.database.close()

    def fetch_from_database(self) -> List[str]:
        if not self.database or not self.database.open:
            self._connect_to_mysql()

        sql = 'SELECT url from tb_urls WHERE url_type = 0'
        cursor = self.database.cursor()
        urls = []
        try:
            cursor.execute(sql)
            url_tuples = cursor.fetchall()
            urls = [cur['url'] for cur in url_tuples]
        except Exception as e:
            print(f'Fetch Data Exception: {e}')
        finally:
            cursor.close()
            self.database.close()

        return urls

class SiteInfoManager:
    connection: Connection = None
    site_infos: List[Tuple[Site, int]]
    comment_infos: List[Tuple[int, str]]
    def __init__(self):
        # 存储景区信息
        # 数据存储格式 (Site, idx)
        self.site_infos = []

        # 存储评论信息
        # 数据存储格式 (idx, content)
        self.comment_infos = []

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

    def insert_to_memory(self, site_info):
        self.site_infos.append(site_info)

    def get_number_of_data_memory(self):
        return len(self.site_infos)

    def insert_comment_to_memory(self, *comment_info):
        self.comment_infos.extend(comment_info)

    def get_number_of_comment_memory(self):
        return len(self.comment_infos)

    def store_into_database(self):
        if not self.connection or not self.connection.open:
            self._connect_to_mysql()

        sql =\
        """
        INSERT INTO tb_sites (name, score, address, hot_degree, introduce, open_time, phone, images, site_idx)
        VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)
        """

        data_lists = []
        for site, idx in self.site_infos:
            images_str = ','.join(site.images)
            data_lists.append((site.name, site.score, site.address,
                               site.hot_degree, site.introduce, site.open_time,
                               site.phone, images_str, idx))

        cursor = self.connection.cursor()

        try:
            cursor.executemany(sql, data_lists)
            self.connection.commit()
            print(f'Insert {len(self.site_infos)} data into database Successfully!')
        except Exception as e:
            print(f'Exception while insert site info to database: {e}')
            self.connection.rollback()
        finally:
            cursor.close()
            self.connection.close()

        # 清空内存中存储的信息
        self.site_infos = []

    def store_comment_to_database(self):
        if not self.connection or not self.connection.open:
            self._connect_to_mysql()

        sql =\
            """
            INSERT INTO `tb_comments` (site_idx, content)  
            VALUES
            (%s, %s)
            """

        data_lists = []
        for idx, comment_info in self.comment_infos:
            data_lists.append((idx, comment_info))

        cursor = self.connection.cursor()

        try:
            cursor.executemany(sql, data_lists)
            self.connection.commit()
            print(f'Insert {len(self.comment_infos)} data into database Successfully!')
        except Exception as e:
            print(f'Exception while insert site info to database: {e}')
            self.connection.rollback()
        finally:
            cursor.close()
            self.connection.close()

        self.comment_infos = []

# @Description: HTML 页面解析器, 根据 HTML 页面解析不同的元素(Url, 页码, 景点信息等)
class HtmlParser:
    def __init__(self, html_text: str):
        self.parser = BeautifulSoup(html_text, 'html.parser')
    def parse_site_url(self) -> List[str]:
        div_tags = self.parser.find_all('div', class_=SITE_URL_CLASS_NAME)
        span_tags = [tag.find('span') for tag in div_tags]
        a_tags = [tag.find('a') for tag in span_tags]
        urls = [a_tag.get('href') for a_tag in a_tags]
        return urls
    def parse_site_info(self) -> Site:
        # 解析单个景点信息
        site_total_info_div = self.parser.find('div', class_=SITE_INFO_CLASS_NAME)
        site_image_info_div = site_total_info_div.find('div', class_=SITE_IMAGE_INFO_CLASS_NAME)
        site_text_info_div = site_total_info_div.find('div', class_=SITE_TEXT_INFO_CLASS_NAME)

        # 解析文本信息
        name = ''
        score = 0.0
        address = ''
        hot_degree = 0.0
        introduce = ''
        open_time = ''
        phone = ''
        images = []

        # 1. 名称
        title_view_div = site_text_info_div.find('div', class_=SITE_TITLE_VIEW_CLASS_NAME)
        if title_view_div:
            title_div = title_view_div.find('div', class_=SITE_TITLE_CLASS_NAME)
            if title_div:
                name_tag = title_div.find('h1')
                if name_tag:
                    name = name_tag.text

        # 2. 评分
        comment_div = site_text_info_div.find('div', class_=SITE_COMMENT_CLASS_NAME)
        if comment_div:
            comment_score_div = comment_div.find('div', class_=SITE_COMMENT_SCORE_CLASS_NAME)
            if comment_score_div:
                score_tag = comment_score_div.find('p', class_=SITE_COMMENT_NUM_CLASS_NAME)
                if score_tag:
                    score = float(score_tag.text)

        # 基础信息解析
        base_info_div = site_text_info_div.find('div', class_=SITE_BASE_INFO_CLASS_NAME)
        base_info_items = base_info_div.find_all('div', class_=SITE_BASE_INFO_CONTEXT_CLASS_NAME)

        # 3. 地址
        if len(base_info_items) >= 1:
            address_tag = base_info_items[0].find('p', class_=SITE_BASE_INFO_CONTEXT_TEXT_CLASS_NAME)
            if address_tag:
                address = address_tag.text

        # 4. 开放时间
        if len(base_info_items) >= 2:
            open_time_p_tag = base_info_items[1].find('p', class_=SITE_OPEN_TIME_CLASS_NAME)
            if open_time_p_tag:
                open_time = open_time_p_tag.get_text(strip=False)

        # 5. 官方电话
        if len(base_info_items) >= 3:
            phone_span_tag = base_info_items[2].find('span', class_=SITE_PHONE_SPAN_CLASS_NAME)
            if phone_span_tag:
                phone = phone_span_tag.get_text(strip=False)

        # 6. 热度
        hot_degree_div = site_text_info_div.find('div', class_=SITE_HOT_DEGREE_CLASS_NAME)
        if hot_degree_div:
            hot_degree = float(hot_degree_div.text)

        # 7. 图片信息
        images_divs = site_image_info_div.find_all('div', class_=SITE_IMAGE_CLASS_NAME)
        # pattern = re.compile(r'background-image: url\("(.*?)"\)')
        if images_divs:
            for image_div in images_divs:
                style = image_div.get('style', '')
                pattern = re.compile(r'url\(([^)]+)\)')
                match = pattern.search(style)
                image_url = style
                if match:
                    image_url = match.group(1)
                images.append(image_url)

        # 8. 介绍信息
        intro_p_tags = self.parser.find('div', class_=SITE_INTRO_CLASS_NAME).find('div').find_all('p')
        if intro_p_tags:
            for p_tag in intro_p_tags:
                introduce += p_tag.text
                if not p_tag.find('img'):
                    introduce += '\n'

        # 返回 site
        cur_site = Site(name=name,
                        score=score,
                        hot_degree=hot_degree,
                        address=address,
                        images=images,
                        introduce=introduce,
                        open_time=open_time,
                        phone=phone)

        return cur_site

    def parser_poi_id(self) -> int:
        # 解析得到  poiId
        script_tag = self.parser.find('script', id=SITE_NEXT_JS_ID)
        poiId = -1
        if script_tag:
            json_data = json.loads(script_tag.text)
            poiId = json_data.get('props').get('pageProps').get('initialState').get('poiDetail').get('poiId')

        return poiId

# @Description: 枚举类型, 描述爬虫状态
class ClawlerState(Enum):
    UNSTART = 1,
    PARSEMAINPAGE = 2,
    PARSESITE = 3,
    END = 4

# @Description: 使用 State Machine 控制爬虫状态(可添加 ProxyPool 等逻辑)
class ClawlerController:
    def __init__(self):
        self.init_url = INIT_URL
        self.state = ClawlerState.UNSTART
        self.url_manager = UrlManager()
        self.site_manager = SiteInfoManager()
    def start(self, start_idx=0, comment_start_idx=0, if_parse_url=True, if_parse_site=True, if_parse_comment=True):
        self.state = ClawlerState.PARSEMAINPAGE
        if if_parse_url:
            # 1. 构造页码对应 Url
            page_urls = []
            cur_page_url = ''
            for i in range(1, 301):
                cur_page_url = f'https://you.ctrip.com/sight/china110000/s0-p{i}.html'
                page_urls.append(cur_page_url)
            # 2. 解析 MainPage 得到 Url
            total_num_of_urls = 0
            for idx, cur_url in enumerate(page_urls):
                cur_resp = requests.get(cur_url)
                cur_html = cur_resp.text
                cur_parser = HtmlParser(cur_html)
                cur_urls = cur_parser.parse_site_url()
                self.url_manager.insert_into_urls(*cur_urls)
                total_num_of_urls += len(cur_urls)
                print('-'*10 + f'{idx + 1}' + '-'*10)
                for url in cur_urls:
                    print(url)
                if total_num_of_urls >= MAX_URL_NUMBER:
                    print('-'*10 + 'End' + '-'*10)
                    print(f'total number of urls: {total_num_of_urls}')
                    break
            self.url_manager.store_into_database()

        if if_parse_site:
            site_urls = self.url_manager.fetch_from_database()
            for idx, site_url in enumerate(site_urls):
                if idx < start_idx:
                    continue
                # 处理单个景点信息
                cur_site = self.parse_cur_site(site_url)
                print('-'*10 + str(idx + 1) + '-'*10)
                print(cur_site)
                self.site_manager.insert_to_memory((cur_site, idx+1))
                if self.site_manager.get_number_of_data_memory() >= MAX_SITES_IN_MEMORY:
                    self.site_manager.store_into_database()

        if if_parse_comment:
            site_urls = self.url_manager.fetch_from_database()
            for idx, site_url in enumerate(site_urls):
                # 处理单个景点评论信息
                if idx < comment_start_idx:
                    continue
                cur_site_comments = self.parse_cur_comments(site_url)

                if len(cur_site_comments) > 0:
                    print('-'*10 + f'{idx + 1}' + '-'*10)
                    print(f'Successfully Get {len(cur_site_comments)} comments')
                    print(cur_site_comments[0])

                    # insert into database
                    comment_info = [(idx+1, comment) for comment in cur_site_comments]
                    self.site_manager.insert_comment_to_memory(*comment_info)
                else:
                    print(f'Successfully Get {len(cur_site_comments)} comments')

                if (self.site_manager.get_number_of_comment_memory() >= MAX_SITES_COMMENT_IN_MEMORY
                        or idx == len(site_urls) - 1):
                    self.site_manager.store_comment_to_database()

    def parse_cur_site(self, site_url) -> Site:
        cur_site_resp = requests.get(site_url)
        parser = HtmlParser(cur_site_resp.text)
        cur_site = parser.parse_site_info()
        return cur_site

    def parse_cur_comments(self, site_url) -> List[str]:
        cur_site_resp = requests.get(site_url)
        parser = HtmlParser(cur_site_resp.text)
        cur_poi_id = parser.parser_poi_id()
        cur_site_comments = []
        if not cur_poi_id == -1:
            for idx in range(2):
                payload = {
                    "arg": {
                        "channelType": 2,
                        "collapseType": 0,
                        "commentTagId": 0,
                        "pageIndex": idx,
                        "pageSize": MAX_PAGE_SIZE,
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
                        cur_site_comments.append(comment['content'])

        return cur_site_comments




if __name__ == '__main__':
    clawer = ClawlerController()
    clawer.start(comment_start_idx= 161, if_parse_url=False, if_parse_site=False)