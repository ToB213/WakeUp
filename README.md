# どきどき！！ファイル破壊アラーム！！！

このプロジェクトは，指定した時間にアラームを設定し，その時間が過ぎたら指定したディレクトリ内のファイルを自動的に削除する，Go言語で作成されたシンプルなWebアプリケーションです．また，ファイル削除を防ぐための「起きた確認」機能もあります．

## 特徴

- アラームの時間を指定して設定
- アラーム時間が過ぎたら、指定したディレクトリ内のファイルを自動削除
- アラーム時間前に「起きた確認」ボタンを押すことで、ファイル削除を回避

## プロジェクト構成

```bash
.
├── go.mod                # Goモジュールファイル
├── go.sum                # Go依存関係ファイル
├── main.go               # メインのGoアプリケーション
├── static
│   └── index.html        # アプリケーションのHTMLファイル
├── templates
│   ├── alarm.html        # アラーム設定のHTMLテンプレート
│   └── style.css         # スタイリング用CSSファイル
├── test                  # ファイル削除のテスト用ディレクトリ
└── wake_up.db            # SQLiteデータベースファイル
```

## 前提条件
- Go（バージョン 1.16 以上）がインストールされていること
- データベースにSQLite3を使用
- Webブラウザ

## 使い方

### アラームを設定する:

フォームにアラーム時間を `YYYY-MM-DD HH:MM` の形式で入力します
削除対象のディレクトリパスを指定します

### 起きた確認:

アラーム時間前に「起きた確認」ボタンをクリックすると，アラームが確認済みとして登録され，ファイル削除が回避されます
