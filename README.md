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

#### 前提条件
BuildKit を有効化してください。

最新バージョン(20.10.16)ではデフォルトで有効になっていると思います。

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

- `make run-normal-build` build が必要な場合のみこちらを利用
- `make run-normal`

デバッカーを利用したい場合、は下記を実行してください。

- `make run-debug-build` # build が必要な場合のみこちらを利用
- `make run-debug`

### テスト

#### BD 接続
一部のテストでDBに接続するテストをしています。

DBに接続できないとエラーとなるので、テストの実行前に下記の設定を行ってください。

- `ローカル上で Docker で動作させる`手順と同様に下記のリポジトリの docker-compose を利用して MySQL のコンテナを起動

https://github.com/nekochans/lgtm-cat-migration

- 下記のコマンドで、Docker 上でテストを実行
  - `make test-build` # build が必要な場合のみこちらを利用
  - `make test`

#### mock の自動生成
mock の自動生成ツール [matryer/moq](https://github.com/matryer/moq) を利用しています。

lgtm-cat-api-dev コンテナ上で下記のようなコマンドを実行することで、mockが自動生成されます。コマンドの詳細は moq のドキュメントを参照。

```
moq -out domain/lgtm_image_repository_test.moq.go ./domain LgtmImageRepository
```

### デプロイ
ローカルからデプロイする場合、下記のコマンドを実行してください。

```
# STG
make deploy-stg

# PROD
make deploy-prod
```
