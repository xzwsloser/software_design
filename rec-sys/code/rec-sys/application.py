from pymilvus import MilvusClient
from rec_model import SiteTower, TwoTower, UserTower
import pandas as pd
import torch
from typing import List
import pymysql 
import pymysql.cursors
from rec_sys import RecSys

MILVUS_ADDRESS = "http://localhost:19530"

MYSQL_HOST = "localhost"
MYSQL_USER = "root"
MYSQL_PORT = 3306
MYSQL_PASSWD = "123456"
MYSQL_DATABASE = "software_design"

class VectorStore:
    def __init__(self):
        self.client = MilvusClient(MILVUS_ADDRESS)
    def insert_site_embed(self, data):
        self.client.insert(collection_name="tb_site_embed",
                           data=data)
    def search_site_embed(self, user_embed, query_size=100):
        res = self.client.search(
            collection_name="tb_site_embed",
            anns_field="site_embed",
            data=[user_embed],
            limit=query_size,
            output_fields=["site_idx"],
            search_params={"metric_type": "IP"}
        )

        site_idxs = [ record['entity']['site_idx'] for record in res[0] ]
        scores = [ record['distance'] for record in res[0] ]

        return site_idxs, scores
    def delte_useless_sites(self):
        del_str = 'site_idx in [780,668,926,288,673,553,300,688,944,954,826,447,196,199,711,968,978,469,221,989,607,992,480,491,366,751,623,880,628,757,636]'
        resp = self.client.delete(
            collection_name="tb_site_embed",
            filter=del_str
        )

        print(resp)

def store_site_embed():
    site_df = pd.read_pickle("../../dataset/site_features.pkl")
    user_tower = UserTower()
    site_tower = SiteTower()
    model = TwoTower(user_tower, site_tower)

    model_path = "../best_model.pt"
    state_dict = torch.load(model_path, map_location="cpu")
    model.load_state_dict(state_dict)

    vector_store = VectorStore()


    for site_id in range(0, 1000):
        site_idx = site_id + 1
        site_row = site_df.iloc[site_id]

        site_score = site_row["score"]
        site_hot_degree = site_row["hot_degree"]
        site_price = site_row["price"]
        site_positive_comment_rate = site_row["positive_comment_rate"]
        site_address = site_row["address"]
        site_introduce_embed = site_row["introduce_embed"]

        batch = {
            "site_score": torch.tensor(site_score, dtype=torch.float32).unsqueeze(0),
            "site_hot_degree": torch.tensor(site_hot_degree, dtype=torch.float32).unsqueeze(0),
            "site_introduce_embed": torch.tensor(site_introduce_embed, dtype=torch.float32).unsqueeze(0),
            "site_address": torch.tensor(site_address, dtype=torch.long).unsqueeze(0),
            "site_price": torch.tensor(site_price, dtype=torch.float32).unsqueeze(0),
            "site_positive_comment_rate": torch.tensor(site_positive_comment_rate, dtype=torch.float32).unsqueeze(0),
        }


        print(f"site_score shape: {batch['site_score'].shape}")

        site_embed = model.get_site_embed(batch)
        site_embed = site_embed.detach().cpu().numpy().reshape(-1)

        data = { "site_idx": site_idx, "site_embed": site_embed }
        vector_store.insert_site_embed(data)
        print(f"successfully insert [{site_idx}] embed to milvus!, dim: {site_embed.shape}")

    vector_store.delte_useless_sites() 


def recommand_for_current_user(model: TwoTower):
    address_id = int(input('请输入你的省份编号'))
    tourist_type = int(input('请输入你的出游类型'))
    like_type = input("请输入你喜欢的景点类型, 使用,分隔")
    targets = input("请输入你的旅游目标, 使用,分隔")
    price_sensitive = int(input("请输入价格是否敏感"))
    attention = input("请输入旅游关注点")

    like_type = [ int(l) for l in like_type.split(",") ]
    targets = [ int(t) for t in targets.split(",") ]
    attention = [ int(a) for a in attention.split(",") ]

    def pad_to_maxlen(l: List[int], pad_id: int, max_len: int):
        if len(l) >= max_len:
            return l[:max_len]
        return l + (max_len - len(l)) * [pad_id]
    
    like_type = pad_to_maxlen(like_type, 14, 6)
    targets = pad_to_maxlen(targets, 13, 6)
    attention = pad_to_maxlen(attention, 8, 4)

    print("="*10 + "你的输入" + "="*10)
    print(f"address_id = {address_id}")
    print(f"tourist_type = {tourist_type}")
    print(f"like_type = {like_type}")
    print(f"targets = {targets}")
    print(f"price_sensitive = {price_sensitive}")
    print(f"attention = {attention}")
    print("="*24)

    user_features = {
        "address": torch.tensor(address_id, dtype=torch.long).unsqueeze(0),
        "tourist_type": torch.tensor(tourist_type, dtype=torch.long).unsqueeze(0),
        "like_type": torch.tensor(like_type, dtype=torch.long).unsqueeze(0),
        "targets": torch.tensor(targets, dtype=torch.long).unsqueeze(0),
        "price_sensitive": torch.tensor(price_sensitive, dtype=torch.long).unsqueeze(0),
        "attention": torch.tensor(attention, dtype=torch.long).unsqueeze(0),
    }

    user_embed = model.get_user_embed(user_features)

    client = VectorStore()

    with torch.no_grad():
        user_embed = user_embed.detach().cpu().numpy().reshape(-1)

    related_site_idxs, scores = client.search_site_embed(user_embed, query_size=50)
    
    print('='*10 + 'vector store search result' + '='*10)
    for idx, (site_idx, score) in enumerate(zip(related_site_idxs, scores)):
        print(f"Rank [{site_idx}] Score: {score}")
    print('='*30)

    connection = pymysql.connect(
        host=MYSQL_HOST,
        port=MYSQL_PORT,
        user=MYSQL_USER,
        passwd=MYSQL_PASSWD,
        database=MYSQL_DATABASE,
        cursorclass=pymysql.cursors.DictCursor
    )

    query_range_str = ",".join([ str(idx) for idx in related_site_idxs ])

    sql = f"""
        SELECT * FROM tb_sites 
            WHERE site_idx IN ({query_range_str})
            ORDER BY FIELD(site_idx, {query_range_str})
    """

    extra_sql = f"""
        SELECT * FROM tb_site_extra 
            WHERE site_idx IN ({query_range_str})
            ORDER BY FIELD(site_idx, {query_range_str})
    """

    print(f"sql = {sql}")

    cursor = connection.cursor()

    try:
        cursor.execute(sql)
        site_infos =  cursor.fetchall()

        cursor.execute(extra_sql)
        site_extra_infos = cursor.fetchall()

        for idx, (site_info, site_extra_info)  in enumerate(zip(site_infos, site_extra_infos)):
            print("="*10 + f"recomment {idx + 1}" + "="*10)
            print(f"景点序号: {site_info['site_idx']}")
            print(f"景点名称: {site_info['name']}")
            print(f"景点评分: {site_info['score']}")
            print(f"景点地址: {site_info['address']}")
            print(f"景点热度: {site_info['hot_degree']}")
            print(f"景点介绍： \n{site_info['introduce']}\n")
            print(f"景点开放时间: {site_info['open_time']}")
            print(f"景点票价: {site_extra_info['price']}")
            print(f"景点好评数: {site_extra_info['positive_comment_count']}")
            print(f"景点差评数: {site_extra_info['negative_comment_count']}")
            print("="*30)
    except Exception as e:
        print(f"Exception: {e}")
    finally:
        cursor.close()
        connection.close()

def recommand_by_rec_sys():
    rec_sys = RecSys()

    address_id = int(input('请输入你的省份编号'))
    tourist_type = int(input('请输入你的出游类型'))
    like_type = input("请输入你喜欢的景点类型, 使用,分隔")
    targets = input("请输入你的旅游目标, 使用,分隔")
    price_sensitive = int(input("请输入价格是否敏感"))
    attention = input("请输入旅游关注点")

    like_type = [ int(l) for l in like_type.split(",") ]
    targets = [ int(t) for t in targets.split(",") ]
    attention = [ int(a) for a in attention.split(",") ]

    print("="*10 + "你的输入" + "="*10)
    print(f"address_id = {address_id}")
    print(f"tourist_type = {tourist_type}")
    print(f"like_type = {like_type}")
    print(f"targets = {targets}")
    print(f"price_sensitive = {price_sensitive}")
    print(f"attention = {attention}")
    print("="*24)

    related_site_idxs, scores = rec_sys.recommand_for_current_user(user_id=0,
                                                           address_id=address_id,
                                                           tourist_type_id=tourist_type,
                                                           like_type=like_type,
                                                           targets_type=targets,
                                                           attention_type=attention,
                                                           price_sensitive=price_sensitive)


    connection = pymysql.connect(
        host=MYSQL_HOST,
        port=MYSQL_PORT,
        user=MYSQL_USER,
        passwd=MYSQL_PASSWD,
        database=MYSQL_DATABASE,
        cursorclass=pymysql.cursors.DictCursor
    )

    query_range_str = ",".join([ str(idx) for idx in related_site_idxs ])

    sql = f"""
        SELECT * FROM tb_sites 
            WHERE site_idx IN ({query_range_str})
            ORDER BY FIELD(site_idx, {query_range_str})
    """

    extra_sql = f"""
        SELECT * FROM tb_site_extra 
            WHERE site_idx IN ({query_range_str})
            ORDER BY FIELD(site_idx, {query_range_str})
    """

    print(f"sql = {sql}")

    cursor = connection.cursor()

    try:
        cursor.execute(sql)
        site_infos =  cursor.fetchall()

        cursor.execute(extra_sql)
        site_extra_infos = cursor.fetchall()

        for idx, (site_info, site_extra_info)  in enumerate(zip(site_infos, site_extra_infos)):
            print("="*10 + f"recomment {idx + 1}" + "="*10)
            print(f"景点序号: {site_info['site_idx']}")
            print(f"景点名称: {site_info['name']}")
            print(f"景点评分: {site_info['score']}")
            # print(f"景点地址: {site_info['address']}")
            print(f"景点热度: {site_info['hot_degree']}")
            # print(f"景点介绍： \n{site_info['introduce']}\n")
            # print(f"景点开放时间: {site_info['open_time']}")
            print(f"景点票价: {site_extra_info['price']}")
            print(f"景点好评数: {site_extra_info['positive_comment_count']}")
            print(f"景点差评数: {site_extra_info['negative_comment_count']}")
            print("="*30)
    except Exception as e:
        print(f"Exception: {e}")
    finally:
        cursor.close()
        connection.close()

    print('='*10 + 'vector store search result' + '='*10)
    for idx, (site_idx, score) in enumerate(zip(related_site_idxs, scores)):
        print(f"Rank [{site_idx}] Score: {score}")
    print('='*30)



if __name__ == "__main__":
    # store_site_embed()
    # site_df = pd.read_pickle("../../dataset/site_features.pkl")
    # user_tower = UserTower()
    # site_tower = SiteTower()
    # model = TwoTower(user_tower, site_tower)

    # model_path = "../best_model.pt"
    # state_dict = torch.load(model_path, map_location="cpu")
    # model.load_state_dict(state_dict)

    # recommand_for_current_user(model)
    recommand_by_rec_sys()
