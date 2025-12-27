from pymilvus import MilvusClient
from rec_model import (
    UserTower,
    SiteTower,
    TwoTower
)
import torch
from pymysql import Connection
import pymysql.cursors
from typing import List, Tuple

MILVUS_ADDRESS = "http://localhost:19530"

MYSQL_HOST = "localhost"
MYSQL_USER = "root"
MYSQL_PORT = 3306
MYSQL_PASSWD = "123456"
MYSQL_DATABASE = "software_design"

# 最大评论数量(用于计算热度)
MAX_COMMENT_COUNT = 175822 

# 最大好评数量
MAX_POSITIVE_COUNT = 170427

# 热度满分值
MAX_HOT_DEGREE = 10.0

# 评分满分值
MAX_SCORE = 5.0

class MilvusStore:
    def __init__(self):
        self.client = MilvusClient(MILVUS_ADDRESS)
        self.site_tb_name = "tb_site_embed"
        self.user_tb_name = "tb_user_embed"

    def insert_user_embed(self, user_id, user_embed):
        data = { 
            "user_id": user_id,
            "user_embed": user_embed
        }

        resp = self.client.insert(
            collection_name=self.user_tb_name,
            data=data
        )

        return resp['insert_count'] == 1
    
    def query_for_user(self, user_id):
        filter = f'user_id == {user_id}'

        resp = self.client.query(
            collection_name=self.user_tb_name,
            filter=filter,
            output_fields=["user_id", "user_embed"]
        )

        if len(resp) == 0:
            return -1, None

        return resp[0]['user_id'], resp[0]['user_embed']
    
    def delete_user(self, user_id):
        filter = f'user_id == {user_id}'

        resp = self.client.delete(
            collection_name=self.user_tb_name,
            filter=filter
        )

        return resp['delete_count'] == 1

    def update_user(self, user_id, user_embed):
        is_exists,_ = self.query_for_user(user_id)

        if is_exists != -1:
            self.delete_user(user_id)

        self.insert_user_embed(user_id, user_embed)

    def retrieval_site_embed(self, user_embed, query_size=200):
        resp = self.client.search(
            collection_name=self.site_tb_name,
            anns_field="site_embed",
            data=[user_embed],
            limit=query_size,
            output_fields=["site_idx"],
            search_params={
                "metric_type": "IP"
            }
        )

        site_idx_list = [ record['entity']['site_idx'] for record in resp[0] ]
        distance_list = [ record['distance'] for record in resp[0] ]

        return site_idx_list, distance_list

class RecSys:
    model: TwoTower
    mysql_client: Connection
    vector_store: MilvusStore
    def __init__(self, rec_model_path: str = "../best_model.pt"):
        user_tower = UserTower()
        site_tower = SiteTower()
        self.model = TwoTower(user_tower, site_tower)
        state_dict = torch.load(rec_model_path, map_location="cpu")
        self.model.load_state_dict(state_dict)

        self.vector_store = MilvusStore()

        self.mysql_client = None

        self.like_type_pad_id, self.like_type_maxlen = 14, 6
        self.target_type_pad_id, self.target_type_maxlen = 13, 6
        self.attention_type_pad_id, self.attention_type_maxlen = 8, 4
    
    def _init_mysql_client(self):
        self.mysql_client = pymysql.connect(
            host=MYSQL_HOST,
            passwd=MYSQL_PASSWD,
            user=MYSQL_USER,
            port=MYSQL_PORT,
            database=MYSQL_DATABASE,
            cursorclass=pymysql.cursors.DictCursor
        )

    def _pad_to_maxlen(self, s: List[int], pad_id: int, max_len: int) -> List[int]:
        s = list(set(s))
        if len(s) >= max_len:
            return s[:max_len]
        
        return s + [pad_id] * (max_len - len(s))
    
    def _query_site_info(self, site_idxs: List[int]):
        if self.mysql_client is None or not self.mysql_client.open:
            self._init_mysql_client()

        query_site_str = ",".join([ str(site_idx) for site_idx in site_idxs ])

        sql = f"""
            SELECT score, hot_degree FROM tb_sites 
                WHERE site_idx IN ({query_site_str})
                ORDER BY FIELD(site_idx, {query_site_str})
        """

        cursor = self.mysql_client.cursor()

        site_scores = []
        site_hot_degrees = []

        try:
            cursor.execute(sql)
            site_infos = cursor.fetchall()
            site_scores = [ site_info['score'] for site_info in site_infos ]
            site_hot_degrees = [ site_info['hot_degree'] for site_info in site_infos ]
        except Exception as e:
            print(f"Exception occurred in query state: {e}")
        finally:
            cursor.close()
            self.mysql_client.close()
        
        return site_scores, site_hot_degrees
    
    def _query_site_extra_info(self, site_idxs: List[int]):
        if self.mysql_client is not None or not self.mysql_client.open:
            self._init_mysql_client()

        query_site_str = ",".join([ str(site_idx) for site_idx in site_idxs ])

        sql = f"""
        SELECT  positive_comment_count, negative_comment_count FROM tb_site_extra 
            WHERE site_idx IN ({query_site_str})
            ORDER BY FIELD(site_idx, {query_site_str})
        """

        cursor = self.mysql_client.cursor() 

        positive_comment_counts = []
        negative_comment_counts = []

        try:
            cursor.execute(sql)
            site_extra_infos = cursor.fetchall()

            positive_comment_counts = [ site_extra_info['positive_comment_count'] for site_extra_info in site_extra_infos ]
            negative_comment_counts = [ site_extra_info['negative_comment_count'] for site_extra_info in site_extra_infos ]

        except Exception as e:
            print(f"Exception occurred in query site extra infos: {e}")
        finally:
            cursor.close()
            self.mysql_client.close()

        return positive_comment_counts, negative_comment_counts

    def rank(self, site_score_pairs: List[Tuple[int, float]],
                   distance_score_rate=1.0,
                   site_score_rate=0,
                   hot_degree_rate=0,
                   comment_count_rate=0) -> List[Tuple[int, float]]:
        """
        Description: 根据景点自身评分、热度、评论数等信息结合推荐系统给出的评分重新计算分数, 并且重排
        """
        site_idxs = [ site_idx for site_idx, _ in site_score_pairs ]
        distance_scores = [ score for _, score in site_score_pairs ]

        # 获取景点本身信息
        site_scores, site_hot_degrees = self._query_site_info(site_idxs)
        positive_comment_counts, negative_comment_counts = self._query_site_extra_info(site_idxs)
        comment_counts = [ p + n for p, n in zip(positive_comment_counts, negative_comment_counts) ]

        site_scores = [ score / MAX_SCORE for score in site_scores ]
        site_hot_degrees = [ hot_degree / MAX_HOT_DEGREE for hot_degree in site_hot_degrees ]
        comment_counts = [ count / MAX_COMMENT_COUNT for count in comment_counts ]

        scores = [ s_1 * distance_score_rate + 
                   s_2 * site_score_rate + 
                   s_3 * hot_degree_rate + 
                   s_4 * comment_count_rate 
                   for s_1, s_2, s_3, s_4 in zip(distance_scores, site_scores, 
                   site_hot_degrees, comment_counts)]
        
        new_site_score_pairs = [(site_idx, score) for site_idx, score in zip(site_idxs, scores)]

        # 降序排序
        new_site_score_pairs = sorted(new_site_score_pairs, key=lambda x: x[1], reverse=True)

        return new_site_score_pairs

    def recommand_for_current_user(self, user_id: int,
                                         address_id: int,
                                         tourist_type_id: int,
                                         price_sensitive: int,
                                         like_type: List[int],
                                         targets_type: List[int],
                                         attention_type: List[int],
                                         update=False,
                                         limit=200):
        user_embed = None

        is_updated = False
        if not update:
            queryed_user_id, queryed_user_embed = self.vector_store.query_for_user(user_id)
            if queryed_user_id != -1:
                is_updated = True
                user_embed = queryed_user_embed
        

        if not is_updated:
            # 更新 or 插入 用户向量
            like_type = self._pad_to_maxlen(like_type, 
                                            self.like_type_pad_id, 
                                            self.like_type_maxlen)

            targets_type = self._pad_to_maxlen(targets_type,
                                            self.target_type_pad_id,
                                            self.target_type_maxlen) 
            
            attention_type = self._pad_to_maxlen(attention_type, 
                                                self.attention_type_pad_id,
                                                self.attention_type_maxlen)

            user_features = {
                "address": torch.tensor(address_id, dtype=torch.long).unsqueeze(0),
                "tourist_type": torch.tensor(tourist_type_id, dtype=torch.long).unsqueeze(0),
                "like_type": torch.tensor(like_type, dtype=torch.long).unsqueeze(0),
                "targets": torch.tensor(targets_type, dtype=torch.long).unsqueeze(0),
                "price_sensitive": torch.tensor(price_sensitive, dtype=torch.long).unsqueeze(0),
                "attention": torch.tensor(attention_type, dtype=torch.long).unsqueeze(0),
            }

            user_embed = self.model.get_user_embed(user_features) 

            with torch.no_grad():
                user_embed = user_embed.detach().cpu().numpy().reshape(-1)

            # 更新用户向量记录
            self.vector_store.update_user(user_id=user_id, user_embed=user_embed)

        # 检索景点信息
        relevant_sites, scores = self.vector_store.retrieval_site_embed(user_embed=user_embed, query_size=limit)

        site_score_pairs = list(zip(relevant_sites, scores))

        ranked_pairs = self.rank(site_score_pairs=site_score_pairs) 
        ranked_site_idxs = [ site_idx for site_idx, _ in ranked_pairs ]
        ranked_site_scores = [ score for _, score in ranked_pairs ]

        return ranked_site_idxs, ranked_site_scores
