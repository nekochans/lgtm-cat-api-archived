# lgtm-cat-api
LGTMeowで利用するWebAPI

## Getting Started

AWS Lambda + Goで実装しています。

デプロイには AWS CLI を利用しています。 なお、Lambdaの初期構築は Terraform で行なっています。

### AWSクレデンシャルの設定

[名前付きプロファイル](https://docs.aws.amazon.com/ja_jp/cli/latest/userguide/cli-configure-profiles.html) を利用しています。

このプロジェクトで利用しているプロファイル名は `lgtm-cat` です。

### テスト

一部のテストでDBに接続するテストをしています。

DBに接続できないとエラーとなるので、テストの実行前に下記の設定を行ってください。

#### MySQL コンテナの起動
下記のリポジトリの docker-compose を利用してテスト用に MySQL のコンテナを起動。

https://github.com/nekochans/lgtm-cat-migration

#### 環境変数の設定
環境変数を設定してください。

```
export TEST_DB_HOST=localhost:3306
export TEST_DB_USER=local_lgtm_cat_user
export TEST_DB_PASSWORD=password
export TEST_DB_NAME=local_lgtm_cat
```

### デプロイ
ローカルからデプロイする場合、下記のコマンドを実行してください。

```
# STG
make deploy-stg

# PROD
make deploy-prod
```
