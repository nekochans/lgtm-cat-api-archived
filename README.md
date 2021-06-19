# lgtm-cat-api
LGTMeowで利用するWebAPI

## Getting Started

AWS Lambda + Goで実装しています。

デプロイには AWS CLI を利用しています。 なお、Lambdaの初期構築は Terraform で行なっています。

### AWSクレデンシャルの設定

[名前付きプロファイル](https://docs.aws.amazon.com/ja_jp/cli/latest/userguide/cli-configure-profiles.html) を利用しています。

このプロジェクトで利用しているプロファイル名は `lgtm-cat` です。

### デプロイ
ローカルからデプロイする場合、下記のコマンドを実行してください。

```
# STG
make deploy-stg

# PROD
make deploy-prod
```
