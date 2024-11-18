アプリケーションの起動：
docker-compose up --build

ブラウザで以下のURLにアクセス：
http://localhost:8080

## API エンドポイント

- `GET /`: メインページ
- `POST /api/v1/users`: ユーザー作成
- `POST /api/v1/todos`: Todo作成
- `GET /api/v1/todos`: Todo一覧取得
- `PUT /api/v1/todos/:id`: Todo更新
- `DELETE /api/v1/todos/:id`: Todo削除
