import torch
from tqdm import tqdm
import ast
from torch import dropout, nn
import torch.nn.functional as F
import torch.optim as optim
from torch.utils.data import Dataset, DataLoader
from sklearn.model_selection import train_test_split
from sklearn.metrics import roc_auc_score
import pandas as pd
from typing import List
import matplotlib.pyplot as plt

def parse_to_list(s) -> List[int]:
    return ast.literal_eval(s)

def pad_to_maxlen(s: List[int], pad_id: int, max_len: int):
    if len(s) >= max_len:
        return s[:max_len]

    return s + (max_len - len(s)) * [pad_id]

class RecDataset(Dataset):
    def __init__(self, user_df: pd.DataFrame, 
                       site_df: pd.DataFrame, 
                       site_idx_to_id: dict, 
                       id_to_site_idx: dict):
        """
        user_df: 用户特征
            - address: 城市, 单标签稀疏特征(35)
            - tourist_type: 出游类型, 单标签稀疏特征(7)
            - like_type: 喜欢景点类型, 多标签稀疏特征(14) -> 0 - 13 maxlen = 6
            - targets: 出游动机, 多标签稀疏特征(13)  -> 0 - 12 maxlen = 6
            - price_sensitive: 价格是否敏感, 单标签稀疏特征(2)
            - attention: 体验关注细节, 多标签稀疏特征(8) -> 0 -> 7 maxlen = 4
            - label: 标签 1 -> 喜欢 0 -> 不喜欢

        site_df: 景点特真
            - score: 评分, 单值稠密特征
            - hot_degree: 热度, 单值稠密特征
            - address: 城市, 单标签稀疏特征(35)
            - introduce_embed: 介绍, 多维稠密特征(256)
            - price: 价格, 单值稀疏特征(单值稠密特征)
            - positive_comment_rate: 好评率, 单值稠密特征
        """
        self.user_df = user_df
        self.site_df = site_df
        self.site2id = site_idx_to_id
        self.id2site = id_to_site_idx

        self.user_df["like_type"] = (
            self.user_df["like_type"]
            .apply(parse_to_list)
            .apply(lambda s: pad_to_maxlen(s, pad_id=14, max_len=6))
        )

        self.user_df["targets"] = (
            self.user_df["targets"]
            .apply(parse_to_list)
            .apply(lambda s: pad_to_maxlen(s, pad_id=13, max_len=6))
        )

        self.user_df["attention"] = (
            self.user_df["attention"]
            .apply(parse_to_list)
            .apply(lambda s: pad_to_maxlen(s, pad_id=8, max_len=4))
        )
    
    
    def __len__(self):
        return len(self.user_df)
    
    def __getitem__(self, idx):
        row = self.user_df.iloc[idx]
        site_id =  int(self.site2id[row["site_idx"]])       
        site_row = self.site_df.iloc[site_id]

        address = row["address"]
        toursit_type = row["tourist_type"]
        like_type = row["like_type"]
        targets = row["targets"]
        price_sensitive = row["price_sensitive"]
        attention = row["attention"]
        label = row["label"]

        site_score = site_row["score"]
        site_hot_degree = site_row["hot_degree"]
        site_address = site_row["address"]
        site_introduce_embed = site_row["introduce_embed"]
        site_price = site_row["price"]
        site_positive_comment_rate = site_row["positive_comment_rate"]

        sample = {
            "site_id": torch.tensor(site_id, dtype=torch.long),
            "address": torch.tensor(address, dtype=torch.long),
            "tourist_type": torch.tensor(toursit_type, dtype=torch.long),
            "like_type": torch.tensor(like_type, dtype=torch.long),
            "targets": torch.tensor(targets, dtype=torch.long),
            "price_sensitive": torch.tensor(price_sensitive, dtype=torch.long),
            "attention": torch.tensor(attention, dtype=torch.long),

            "site_score": torch.tensor(site_score, dtype=torch.float32),
            "site_hot_degree": torch.tensor(site_hot_degree, dtype=torch.float32),
            "site_introduce_embed": torch.tensor(site_introduce_embed, dtype=torch.float32),
            "site_address": torch.tensor(site_address, dtype=torch.long),
            "site_price": torch.tensor(site_price, dtype=torch.float32),
            "site_positive_comment_rate": torch.tensor(site_positive_comment_rate, dtype=torch.float32),

            "label": torch.tensor(label, dtype=torch.float32)
        }

        return sample

class UserTower(nn.Module):
    def __init__(self,
                 address_dim=35,
                 tourist_type_dim=7,
                 like_type_dim=15,
                 targets_dim=14,
                 price_sensitive_dim=2,
                 attention_dim=9,
                 embed_dim=32,
                 out_dim=64,
                 hidden_dim=128,
                 dropout_rate=0.1):
        super(UserTower, self).__init__()

        self.address_embedding = nn.Embedding(address_dim, embed_dim)
        self.toursit_type_embedding = nn.Embedding(tourist_type_dim, embed_dim)
        self.like_type_embedding = nn.Embedding(like_type_dim, embed_dim)
        self.targets_embedding = nn.Embedding(targets_dim, embed_dim)
        self.price_sensitive_embedding = nn.Embedding(price_sensitive_dim, embed_dim)
        self.attention_embedding = nn.Embedding(attention_dim, embed_dim)

        # 32 * 6 = 192
        input_dim = embed_dim * 6

        self.mlp = nn.Sequential(
            nn.Linear(input_dim, hidden_dim),
            nn.ReLU(),
            nn.Dropout(dropout_rate),
            nn.Linear(hidden_dim, out_dim)
        )
    
    def _pool(self, matrix, embed):
        # matrix: [B, L]
        matrix_embed = embed(matrix)
        # matrix_embed: [B, L, embed_dim]
        s = matrix_embed.mean(dim=1)
        # s: [B, embed_dim]
        return s

    def forward(self, batch):
        # 单值稀疏特征 [B, embed_dim]
        address_embed = self.address_embedding(batch["address"])
        tourist_type_embed = self.toursit_type_embedding(batch["tourist_type"])
        price_sensitive_embed = self.price_sensitive_embedding(batch["price_sensitive"])

        # 多值稀疏特征 [B, L, embed_dim] -> [B, embed_dim]
        like_type_embed = self._pool(batch["like_type"], self.like_type_embedding)
        targets_embed = self._pool(batch["targets"], self.targets_embedding)
        attention_embed = self._pool(batch["attention"], self.attention_embedding)

        # 特征拼接
        features = torch.cat(
            [address_embed, tourist_type_embed, price_sensitive_embed, 
             like_type_embed, targets_embed, attention_embed],
             dim=-1
        )

        user_embed = self.mlp(features)
        # 最后一个维度上做归一化操作
        user_embed = F.normalize(user_embed, p=2, dim=-1)

        return user_embed

class SiteTower(nn.Module):
    def __init__(self, 
                 address_dim=35,
                 introduce_dim=256,
                 embed_dim=32,
                 output_dim=64,
                 hidden_dim=128,
                 introduce_final_dim=64,
                 dropout_rate=0.1):

        super(SiteTower, self).__init__()

        # 地址特征嵌入层
        self.address_embedding = nn.Embedding(address_dim, embed_dim)

        input_dim = embed_dim + introduce_final_dim + 4


        self.introduce_proj = nn.Sequential(
            nn.Linear(introduce_dim, introduce_final_dim),  # 256→64，平衡维度占比
            nn.ReLU(),
            nn.Dropout(dropout_rate)
        )

        # self.register_buffer("address_features", address_features)
        # self.register_buffer("introduce_features", introduce_features)
        # self.register_buffer("other_features", other_features)

        self.dense_norm = nn.BatchNorm1d(4)

        self.mlp = nn.Sequential(
            nn.Linear(input_dim, hidden_dim),
            nn.ReLU(),
            nn.Dropout(dropout_rate),
            nn.Linear(hidden_dim, output_dim)
        )

    def forward(self, batch):
        address = batch["site_address"]
        intro_embed = batch["site_introduce_embed"]
        score = batch["site_score"].unsqueeze(1)
        hot_degree = batch["site_hot_degree"].unsqueeze(1)
        price = batch["site_price"].unsqueeze(1)
        positive_comment_rate = batch["site_positive_comment_rate"].unsqueeze(1)

        other = torch.cat([score, hot_degree, price, positive_comment_rate], dim=-1)
        
        if other.size(0) != 1:
            other = self.dense_norm(other)

        address_embed = self.address_embedding(address)
        intro_embed = self.introduce_proj(intro_embed)

        site_features = torch.cat([
                address_embed, 
                intro_embed, 
                other
                ], dim=-1)

        site_embed = self.mlp(site_features)

        # 归一化操作
        site_embed = F.normalize(site_embed, p=2, dim=-1)
        
        return site_embed

class TwoTower(nn.Module):
    def __init__(self, user_tower: nn.Module,
                       site_tower: nn.Module):
        super(TwoTower, self).__init__()
        self.user_tower = user_tower
        self.site_tower = site_tower
    
    def get_site_embed(self, x):
        return self.site_tower(x)

    def get_user_embed(self, x):
        return self.user_tower(x)

    def forward(self, batch):
        user_embed = self.user_tower(batch)
        site_embed = self.site_tower(batch)

        # [B, embed_dim] -> [B]
        # 注意在 userTower 和 siteTower 中已经归一化操作了
        score = (user_embed * site_embed).sum(dim=-1)

        return score

def train(model: nn.Module,
          train_loader: DataLoader,
          valid_loader: DataLoader,
          epochs=30,
          lr=1e-3,
          weight_decay=1e-5):
    criterion = nn.BCEWithLogitsLoss()
    optimizer = optim.Adam(
        model.parameters(),
        lr=lr,
        weight_decay=weight_decay
    )

    avg_train_loss_list = []
    avg_valid_loss_list = []
    auc_list = []

    for epoch in range(1, epochs+1):
        model.train()
        total_loss = 0.0
        total_sample = 0

        loader = tqdm(train_loader, desc=f"Epoch {epoch}/{epochs}")
        for batch in loader:
            score = model(batch)
            labels = batch["label"]

            loss = criterion(score, labels)

            # 梯度下降 + 更新参数
            optimizer.zero_grad()
            loss.backward()
            optimizer.step()

            total_loss += loss.item() * labels.size(0)
            total_sample += labels.size(0)
        
        avg_train_loss = total_loss / total_sample

        avg_valid_loss, auc = eval(
            model=model,
            valid_loader=valid_loader,
            criterion=criterion
        )

        print('='*20)
        print(f"Epoch {epoch}/{epochs}")
        print(f"train_loss={avg_train_loss:.4f}")
        print(f"valid_loss={avg_valid_loss:.4f}")
        print(f"valid_auc={auc}")
        print('='*20)

        avg_train_loss_list.append(avg_train_loss)
        avg_valid_loss_list.append(avg_valid_loss)
        auc_list.append(auc)

        # 保存模型
        model_save_path = f"../weights/model_{epoch}.pt"
        torch.save(model.state_dict(), model_save_path)

    return avg_train_loss_list, avg_valid_loss_list, auc_list

def eval(model: nn.Module,
         valid_loader: DataLoader,
         criterion: nn.BCEWithLogitsLoss):
         model.eval()
         total_valid_loss = 0.0
         total_valid_sample = 0
         all_pred_score = []
         all_labels = []

         with torch.no_grad():
            for batch in valid_loader:
                scores = model(batch)
                labels = batch["label"]

                loss = criterion(scores, labels)

                total_valid_loss += loss.item() * labels.size(0)
                total_valid_sample += labels.size(0)

                probs = torch.sigmoid(scores)
                all_pred_score.append(probs)
                all_labels.append(labels)

            avg_valid_loss = total_valid_loss / total_valid_sample

            try:
                y_pred = torch.cat(all_pred_score).numpy()
                y_label = torch.cat(all_labels).numpy()
                auc = roc_auc_score(y_label, y_pred)
            except Exception as e:
                print(f"Exception In Auc Caculator: {e}")
                auc = 0.0
            
            return avg_valid_loss, auc

def draw_metrics(avg_train_loss_list: List, avg_valid_loss_list: List, auc_list: List):
    def draw_and_save(result_list: List, 
                      title: str,
                      label: str, 
                      y_label: str, 
                      save_path: str):
        epochs = len(result_list)
        epoch_list = range(1, epochs+1)

        plt.figure(figsize=(8, 5))
        plt.plot(epoch_list, result_list, marker='o', label=label)

        plt.xlabel('epoch')
        plt.ylabel(y_label)
        plt.title(title)
        plt.grid(True, linestyle='--', alpha=0.5)
        plt.legend() # 显示图例

        plt.tight_layout()          # 自动调整边距
        plt.savefig(save_path, dpi=300)  # 保存到文件
        plt.close()                 # 关闭图像，避免在循环中占内存

    draw_and_save(avg_train_loss_list,
                  "train loss",
                  "loss",
                  "loss",
                  "../run/train_loss.png")
    
    draw_and_save(avg_valid_loss_list,
                  "valid loss",
                  "loss",
                  "loss",
                  "../run/valid_loss.png")
    
    draw_and_save(auc_list,
                 "valid auc",
                 "auc",
                 "auc",
                 "../run/valid_auc.png")


if __name__ == '__main__':
    user_features_path = "../../dataset/user_features.csv"
    site_features_path = "../../dataset/site_features.pkl"
    user_df = pd.read_csv(user_features_path)
    site_df = pd.read_pickle(site_features_path)

    # print('DEBUG')
    # print("label NaN count:", user_df["label"].isna().sum())
    # print(user_df["label"].unique())
    # print('DEBUG END')

    train_df, test_df = train_test_split(
        user_df,
        test_size=0.1,
        random_state=42,
        shuffle=True
    )

    site_idx_to_id = { id+1: id for id in range(0, 1000) }
    id_to_site_idx = { id: id+1 for id in range(0, 1000) }

    train_dataset = RecDataset(
        user_df=train_df,
        site_df=site_df,
        site_idx_to_id=site_idx_to_id,
        id_to_site_idx=id_to_site_idx
    )

    test_dataset = RecDataset(
        user_df=test_df,
        site_df=site_df,
        site_idx_to_id=site_idx_to_id,
        id_to_site_idx=id_to_site_idx
    )

    train_loader = DataLoader(
        train_dataset,
        batch_size=256,
        shuffle=True
    )

    valid_loader = DataLoader(
        test_dataset,
        batch_size=256,
        shuffle=False
    )

    userTower = UserTower()
    siteTower = SiteTower()

    model = TwoTower(userTower, siteTower)

    print("模型训练开始")
    avg_train_list, avg_valid_list, auc_list = train(model, train_loader, valid_loader)
    print('模型训练结束')

    draw_metrics(avg_train_list, avg_valid_list, auc_list)
    print('成功保存训练结果')
