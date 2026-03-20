# simpledb

書籍 [**「Database Design and Implementation」**](https://link.springer.com/book/10.1007/978-3-030-33836-7) をもとに、RDBをGoでスクラッチ実装するプロジェクトです。

## 概要

Javaで書かれたサンプル実装をGoで再実装しながら、RDBの内部構造（ストレージエンジン、ログ管理、バッファ管理、クエリ処理等）を学ぶことを目的としています。

## 実装状況

| 章 | 内容 | 状態 |
|---|---|---|
| 3章 | ファイル管理（FileManager / Page / BlockID） | 完了 |
| 4章 | ログ管理（LogManager） | 実装中 |

## ディレクトリ構成

```
.
├── internal/
│   ├── file/       # ファイル管理（3章）
│   └── log/        # ログ管理（4章）
└── main.go
```

## 開発

### 必要なツール

```sh
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
```

### コマンド

```sh
make test   # テスト実行
make lint   # 静的解析
make fmt    # フォーマット
```
