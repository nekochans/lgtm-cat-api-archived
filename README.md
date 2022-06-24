# lgtm-cat-api
[![ci](https://github.com/nekochans/lgtm-cat-api/actions/workflows/ci.yml/badge.svg)](https://github.com/nekochans/lgtm-cat-api/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/nekochans/lgtm-cat-api/branch/main/graph/badge.svg?token=BCZABFS4P0)](https://codecov.io/gh/nekochans/lgtm-cat-api)

LGTMeowで利用するWebAPI

## Getting Started

AWS Lambda + Goで実装しています。

デプロイには AWS CLI を利用しています。 なお、Lambdaの初期構築は Terraform で行なっています。

### AWSクレデンシャルの設定

[名前付きプロファイル](https://docs.aws.amazon.com/ja_jp/cli/latest/userguide/cli-configure-profiles.html) を利用しています。

このプロジェクトで利用しているプロファイル名は `lgtm-cat` です。

### ローカル上で Docker で動作させる

#### MySQL コンテナの起動
下記のリポジトリの docker-compose を利用して MySQL のコンテナを起動。

https://github.com/nekochans/lgtm-cat-migration

#### 環境変数の設定
環境変数を設定してください。

`lgtm-cat-migration` と同じ値を設定する必要があります。

```
# テストの実行に必要
export TEST_DB_HOST=localhost:3306
export TEST_DB_USER=local_lgtm_cat_user
export TEST_DB_PASSWORD=password
export TEST_DB_NAME=local_lgtm_cat

export DB_HOSTNAME=mysql:3306
export DB_USERNAME=local_lgtm_cat_user
export DB_PASSWORD=password
export DB_NAME=local_lgtm_cat
```

### Docker を起動
下記のコマンドを実行することで、Docker で起動することができます。

`make run-normal`

デバッカーを利用したい場合、は下記を実行してください。

`run-debug`

### テスト

一部のテストでDBに接続するテストをしています。

DBに接続できないとエラーとなるので、テストの実行前に下記の設定を行ってください。

`ローカル上で Docker で動作させる`手順と同様に下記のリポジトリの docker-compose を利用して MySQL のコンテナを起動してください。

https://github.com/nekochans/lgtm-cat-migration

### デプロイ
ローカルからデプロイする場合、下記のコマンドを実行してください。

```
# STG
make deploy-stg

# PROD
make deploy-prod
```
