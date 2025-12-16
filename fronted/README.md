# 可视化旅游景点推荐系统(前端)
## 项目启动
1. 首先需要对于项目进行构建:
```shell
npm run build
```
2. 生成静态资源文件夹 `dist`, 可以创建一个专门的目录用于运行:
```shell
mkdir run
mv ./dist ./run
```
3. 配置 `nginx`:
```nginx
# 必须包含 events 块
events {
    worker_connections 1024;
}

http {
    # 包含 MIME 类型，否则网页可能没有样式
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    server {
        listen 8888;
        server_name localhost;

        root /home/xzw/projects/software_design/fronted/run/dist;
        index index.html;

        location / {
            try_files $uri $uri/ /index.html;
        }

        location /api/ {
            proxy_pass http://localhost:9999/api/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }
    }
}
```

4. 启动 `nginx` 服务器(同时注意目标目录的权限问题), 同时注意使用绝对路径
```shell
chmod +x <your-home-path>
sudo chmod -R 755 /<your-project-path>/dist 

sudo nginx -c /<your-project-path/nginx.conf
```
