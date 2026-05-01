# my-todo
## 環境
- Go 1.25.4
- MySQL 8.0
- Atlas 1.2.1

## 起動方法
1. Docker起動
    ```bash
    docker-compose up -d
    ```
2. DBマイグレーション
    ```bash
    make atlas-apply
    ```
3. サーバ起動
    ```bash
    make run
    ```
