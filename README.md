# kanban API

カンバン方式のタスク管理アプリ **kanban** のAPIサーバーのリポジトリです。

- アプリURL: https://k4nban.com
- アプリ概要: https://github.com/jo-tbhac/kanban-readme
- フロントエンドのリポジトリ: https://github.com/jo-tbhac/kanban

## 必要条件

以下導入手順はMacOS専用です。

### Go

```
# Goのインストール
brew install go

# GOPATHの設定
echo "export GOPATH=$(go env GOPATH)" >> ~/.bash_profile
echo "export PATH=$PATH:$(go env GOPATH)/bin" >> ~/.bash_profile
source ~/.bash_profile
```

### MySQL

```
# MySQLのインストール
brew install mysql

# MySQLの起動
mysql.server start

# MySQLへログイン
mysql -u root

# データベースの作成
CREATE DATABASE <任意のデータベース名> DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
```

### AWS

ファイルアップロード機能を使用するためにAWSアカウントが必要です。

またS3バケットの作成と、S3へアクセスするためのAccessKeyIDとSecretAccessKeyを取得する必要があります。

#### AWSアカウントの作成

[こちら](https://aws.amazon.com/jp/premiumsupport/knowledge-center/create-and-activate-aws-account/)を参照してください。

#### AccessKeyIDとSecretAccessKeyの取得

セキュリティのために`AmazonS3FullAccess`のみアクセスが許可されたIAMユーザーを作成することを推奨します。（詳細については[こちら](https://docs.aws.amazon.com/ja_jp/IAM/latest/UserGuide/id_users_create.html)を参照してください）

取得したAccessKeyIDとSecretAccessKeyは、環境変数もしくは共有資格情報ファイルに設定してください。（詳細については[こちら](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)を参照してください）

#### S3バケットの作成

[こちら](https://docs.aws.amazon.com/ja_jp/AmazonS3/latest/user-guide/create-bucket.html)を参照してください。

## プロジェクトのセットアップ

以下のコマンドでリポジトリをクローンします。

```
git clone https://github.com/jo-tbhac/kanban-api.git
```

プロジェクトのルートディレクトリに移動し、設定ファイルを作成します。

```
cp config.yml.sample config.yml
```

作成した`config.yml`を開き、自身の環境に応じた値を入力してください。

```
aws:
    bucket: 作成したS3バケットの名前
    region: バケットを作成したリージョン  # 例: ap-northeast-1
database:
    user: データベースのユーザー名
    name: データベース名
    host: データベースのホスト
    password: データベースのパスワード
    driver: データベースのドライバー
    log_mode: データベースアクセスのログを表示するかどうか true or false
web:
    port: アプリケーションを起動するポート番号
    origin: httpアクセスを許可するクライアントサイドのOrigin  # 例: http://localhost:3000
```

ここまで完了したら以下のコマンドでアプリケーションを実行できます。

```
go run main.go
```
