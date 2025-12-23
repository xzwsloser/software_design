from datetime import time
import pandas as pd
from openai import OpenAI, RateLimitError
from typing import List
import time
import json

"""
用户特征:

ip 地址

景点喜好类型:
    0 亲子同乐
    1 观光游览
    2 夜游观景
    3 自然风光
    4 名胜古迹
    5 户外活动
    6 展馆展览
    7 动植物园
    8 冬季滑雪
    9 主题乐园
    10 体闲娱乐
    11 温泉泡汤
    12 水上活动
    13 空中体验

出游类型(单标签):
    0 其他出游     1308 + NaN
    1 单独旅行     4772
    2 商务出差       18
    3 家庭亲子    16182
    4 情侣夫妻     7140
    5 朋友出游    10630
    6 陪同父母     2254

出游动机(多标签):
    0 其他
    1 历史文化溯源
    2 自然景观观赏
    3 亲子遛娃互动
    4 主题乐园狂欢
    5 城市地标打卡
    6 休闲度假放松
    7 网红地标打卡
    8 文化艺术体验
    9 户外探险猎奇
    10 家庭团聚出游
    11 治愈系散心
    12 节庆主题体验

价格敏感与否(二分类):
    0 价格敏感型
    1 价格不敏感

体验关注细节(多分类): 
    0 排队效率敏感
    1 设备完善度敏感
    2 服务质量敏感
    3 行程规划偏好
    4 舒适度敏感（
    5 导览体验敏感
    6 消费透明敏感
    7 无障碍设施敏感
"""
SILICONFLOW_API_KEY = "<YOUR_API_KEY>"
CHAT_MODEL_URL = "https://api.siliconflow.cn/v1"

GLM_MODEL_106B = 'zai-org/GLM-4.5-Air'
GLM_MODEL_4_6_335B = 'zai-org/GLM-4.6'

def make_prompt(comments: List[str]) -> str:
    """
    将评论以结构化 JSON 的形式嵌入 prompt，显式给出 num_comments，减少长度不一致问题。
    """
    num_comments = len(comments)

    # 先构造一个结构化的输入 JSON 字符串，避免手写 {{ }}
    input_json = {
        "num_comments": num_comments,
        "comments": comments,
    }
    input_json_str = json.dumps(input_json, ensure_ascii=False, indent=2)

    PROMPT = f"""
你是一个旅游评论标签助手。请根据给定的多条评论，从多个维度为每条评论分别选择合适的标签。

================
景点喜好类型 (likeType，选择 1 - 6 个标签，对应编号为整数列表):
    0 亲子同乐
    1 观光游览
    2 夜游观景
    3 自然风光
    4 名胜古迹
    5 户外活动
    6 展馆展览
    7 动植物园
    8 冬季滑雪
    9 主题乐园
    10 体闲娱乐
    11 温泉泡汤
    12 水上活动
    13 空中体验

出游动机 (target，选择 1 - 6 个标签，对应编号为整数列表):
    0 其他
    1 历史文化溯源
    2 自然景观观赏
    3 亲子遛娃互动
    4 主题乐园狂欢
    5 城市地标打卡
    6 休闲度假放松
    7 网红地标打卡
    8 文化艺术体验
    9 户外探险猎奇
    10 家庭团聚出游
    11 治愈系散心
    12 节庆主题体验

价格敏感与否 (priceSensitive，选择 0 或 1，对应单个整数):
    0 价格敏感型
    1 价格不敏感

体验关注细节 (attention，多标签):
    0 排队效率敏感
    1 设备完善度敏感
    2 服务质量敏感
    3 行程规划偏好
    4 舒适度敏感
    5 导览体验敏感
    6 消费透明敏感
    7 拍照出片敏感
================
要求：
    - likeType 必须是一个包含 1 至 6 个整数的列表，每个整数只能是 0~13 之间的编号, 没有符合要求的标签返回 [10]。
    - target 必须是一个包含 1 至 6 个整数的列表，每个整数只能是 0~12 之间的编号, 没有符合要求的标签返回 [0]。
    - priceSensitive 必须是一个整数，只能是 0 或 1。
    - attention 必须是一个包含 1 至 4 个整数的列表，每个整数只能是 0~7 之间的编号, 没有符合要求的标签的返回 [4]。
    - 返回的数据只使用每一项标签前面的“编号”（整数），不要返回中文。
    - 需要对每一条评论分别打标签，按照评论在输入中的顺序依次输出结果。
    - 最终回答中只输出一个合法的 JSON 数组，不要包含任何多余文字、注释或解释。
    - 不要使用代码块标记 ```，不要加前缀或后缀文字。

输出格式说明：
    - 输入中共有 num_comments 条评论（从 0 到 num_comments-1 编号）。
    - 返回一个长度严格等于 num_comments 的 JSON 数组。
    - 数组中的第 i 个元素对应输入 JSON 中 comments[i] 的标签结果。
    - 每个元素都是一个对象，包含字段 "likeType"、"target"、"priceSensitive"、"attention"。

请先在你的思考中数出 num_comments = N（不要在最终答案中输出这个数字或任何解释），
然后生成一个长度必须等于 N 的 JSON 数组，数组第 i 个元素对应第 i 条评论。

下面是需要你分析的评论输入 JSON：
========
{input_json_str}
========

注意：
    - 输出的 JSON 数组长度一定要和 num_comments 完全一致。
    - 例如 num_comments = {num_comments}，则输出的 JSON 数组中必须包含 {num_comments} 个对象。
    - 严格保证：数组长度 = num_comments，且数组第 i 个元素对应 comments[i]。

现在开始生成结果：只输出 JSON 数组本身，不要输出任何多余文字。
"""
    return PROMPT.strip()

class LLM:
    def __init__(self):
        self.client = OpenAI(api_key=SILICONFLOW_API_KEY,
                             base_url=CHAT_MODEL_URL)
        # self.model = CHAT_MODEL
        self.model = GLM_MODEL_4_6_335B
    def chat(self, content: List[str]):
        prompt = make_prompt(content)
        resp = self.client.chat.completions.create(  
            model=self.model,  
            messages=[    
                {"role": "user", "content": prompt}  
            ],  
            temperature=0.7,  
            max_tokens=1024*40,
        )
        return resp.choices[0].message.content
    def chat_with_try(self, content: List[str]):
        prompt = make_prompt(content)

        max_retries = 7          # 最大重试次数
        base_backoff = 1.0       # 初始退避时间（秒）

        for attempt in range(max_retries):
            try:
                resp = self.client.chat.completions.create(
                    model=self.model,
                    messages=[
                        {"role": "user", "content": prompt}
                    ],
                    temperature=0.7,
                    max_tokens=1024 * 40,
                )
                return resp.choices[0].message.content

            except RateLimitError as e:
                # 触发 TPM / RPM 等限流时进行重试
                if attempt == max_retries - 1:
                    # 最后一次仍失败，直接抛出
                    raise

                # 指数退避 + 随机抖动，避免所有请求同时重试
                sleep_time = base_backoff * (2 ** attempt) + random.uniform(0, 0.5)
                print(f"TPM Error, sleep time: {sleep_time}")
                time.sleep(sleep_time)

def parse_comments(start_idx = 1, 
                   end_idx=1000, 
                   default_positive_start_idx = 0,
                   default_negative_start_idx = 0,
                   default_positive_size=20,
                   default_negative_size=10,
                   retry_list: List[int] = []):
    positive_df = pd.read_csv("../dataset/tb_positive_comment_extra2.csv")
    negative_df = pd.read_csv("../dataset/tb_negative_comment_extra.csv")
    llm = LLM()

    process_list = retry_list if len(retry_list) != 0 else range(start_idx, end_idx+1)
    for site_idx in process_list:
        start = time.time()
        positive_comment_lists = (
            positive_df.loc[positive_df["site_idx"] == site_idx, "content"]
            .fillna("")
            .tolist()
        )

        negative_comment_list = (
            negative_df.loc[negative_df["site_idx"] == site_idx, "content"]
            .fillna("")
            .tolist()
        )

        positive_ip_list = (
            positive_df.loc[positive_df["site_idx"] == site_idx, "ip"]
            .fillna("")
            .tolist()
        )

        negative_ip_list = (
            negative_df.loc[negative_df["site_idx"] == site_idx, "ip"]
            .fillna("")
            .tolist()
        )

        positive_tourist_list = (
            positive_df.loc[positive_df['site_idx'] == site_idx, "tourist_type"]
            .fillna("")
            .tolist()
        )

        negative_tourist_list = (
            negative_df.loc[negative_df['site_idx'] == site_idx, "tourist_type"]
            .fillna("")
            .tolist()
        )

        if len(positive_comment_lists) == 0 or len(negative_comment_list) == 0:
            continue

        pos_slice = slice(default_positive_start_idx, default_positive_start_idx+default_positive_size)
        positive_comment_to_process = positive_comment_lists[pos_slice]
        positive_ip_to_save = positive_ip_list[pos_slice]
        positive_tourist_to_save = positive_tourist_list[pos_slice]

        neg_slice = slice(default_negative_start_idx, default_negative_start_idx+default_negative_size)
        negative_comment_to_process = negative_comment_list[neg_slice]
        negative_ip_to_save = negative_ip_list[neg_slice]
        negative_tourist_to_save = negative_tourist_list[neg_slice]


        comments = positive_comment_to_process + negative_comment_to_process
        L = len(comments)
        print(f"L = {L}")

        ip_to_save = positive_ip_to_save + negative_ip_to_save
        tourist_to_save = positive_tourist_to_save + negative_tourist_to_save
        label_to_save = [1]*len(positive_comment_to_process) + [0]*len(negative_comment_to_process)
        site_idx_to_save = [site_idx]*L

        content = llm.chat_with_try(comments)

        try:
            datas = json.loads(content)
            print(datas)
            likeTypes = []
            targets = []
            priceSensitive = []
            attention = []

            print(len(datas))
            for obj in datas:
                likeTypes.append('|'.join(map(str,obj['likeType'])))
                targets.append('|'.join(map(str, obj['target'])))
                priceSensitive.append(obj['priceSensitive'])
                attention.append('|'.join(map(str,obj['attention'])))

            # print(len(site_idx_to_save))
            # print(len(ip_to_save))
            # print(len(tourist_to_save))
            # print(len(likeTypes))
            # print(len(targets))
            # print(len(priceSensitive))
            # print(len(attention))
            # print(len(label_to_save))

            df = pd.DataFrame({
                "site_idx": site_idx_to_save,
                "ip": ip_to_save,
                "tourist_type": tourist_to_save,
                "likeType": likeTypes,
                "targets": targets,
                "priceSensitive": priceSensitive,
                "attention": attention,
                "label": label_to_save
            })

            df.to_csv(f"../dataset/comment1/comment_databaset_{site_idx}.csv")
        except Exception as e:
            print(f"cur site_idx: {site_idx} Error: {e}")
            write_to_log(site_idx, e)
        
        end = time.time()

        print(f"[{site_idx}] 耗时 {end - start:.6f} s")
        time.sleep(5)

def write_to_log(site_idx, error_msg):
    with open('./error1.log', 'a+') as f:
        f.write(f"[{site_idx} error msg: {error_msg}]\n")
        
if __name__ == '__main__':
    parse_comments(
                   start_idx=0,
                   default_positive_size=10,
                   default_negative_size=5,
                   )

