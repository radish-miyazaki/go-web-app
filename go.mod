module github.com/radish-miyazaki/go-web-app

go 1.19

// errgroup: ゴルーチンで処理を行った際にエラーを返せるようにするためのパッケージ
require golang.org/x/sync v0.1.0

// 環境変数を取り扱うためのパッケージ
require github.com/caarlos0/env/v7 v7.0.0 // indirect

require (
	github.com/go-chi/chi/v5 v5.0.8 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.11.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
)
