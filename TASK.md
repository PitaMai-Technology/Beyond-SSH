# 開発タスクリスト

## 1. プロジェクトの初期設定
- [x] Vite + Vue 3 プロジェクトのセットアップ
- [x] UIフレームワーク（Vuetify等）の導入（Twitter風UI・TreeView用）
- [x] Go言語によるWASMプロジェクトの初期化（`main.go`の作成）
- [x] Go WASMをローカルのnpmパッケージとしてビルド・インストールするスクリプトの作成

## 2. フロントエンド開発（Vite + Vue）
- [x] Twitter風の全体レイアウト作成（左サイドバー、メインエリア、右サイドバー）
- [x] xterm.jsを使用しない、独自のSSHコマンド入力・結果表示コンポーネントの実装
- [x] Material UIのTreeView風のファイルツリー表示コンポーネント（SFTP用）の実装
- [x] セッション管理・ユーザー認証UIの実装

## 3. Go WASMバックエンドモジュール
- [x] JavaScriptから呼び出せるGo関数（`connectSSH`, `executeCommand`等）のWASMエクスポート
- [x] `golang.org/x/crypto/ssh` を用いたSSH接続・コマンド実行ロジックの実装
- [ ] SSH キー認証対応
- [x] SFTPプロトコルを用いたファイル操作ロジックの実装
- [x] ネットワーク通信制限を回避するためのプロキシ通信ロジック（またはNode.js実行環境の整備）の実装

## 4. 結合テストとUI調整
- [x] フロントエンドとWASMモジュールの結合
- [ ] 実際のサーバーに対するSSH/SFTP接続テスト
- [x] デザイン・マイクロインタラクション（Twitter風の見た目）のブラッシュアップ
