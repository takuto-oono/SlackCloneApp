# 環境構築
1. コードをクローン
2. 開発メンバーの誰かからconfig.iniを共有してもらい、backend/config配下に保存する。
3. docker-compose upを実行

# テスト用データ生成コマンド(フロントエンド用)
## 注意
保存されているデータはすべて削除される。(いちいちデータ削除しなくても大丈夫)

## 利用手順
1. docker-compose up
2. docker ps コマンドでbackendのどちらか(1 or 2)のどちらかのコンテナIDを調べる
3. docker exec -it {ステップ2で取得したもの} bash コマンドでコンテナ内に入る
4. go/src 配下で go run main.go testdata コマンドを実行
5. しばらく待機して、フロントエンドにアクセス

## 保存されているデータ(改良していく予定なので、実際のコード(set_data.go)を見たほうが正確かも)
ユーザー
ユーザー名: user1, user2, user3
パスワード ３ユーザーともabc123

workspaceはtest_workspaceは必ずある。
