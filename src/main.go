package main

// 簡単かつ効率的にURLアクセスやHTTPリクエストを実行するためのツール
// exe化して配布することを想定している
// 使い方は、PATHを通して、コマンドライン引数にURLを指定するだけ
// 例: gofetch --url https://example.com
// 例: gofetch -u https://example.com
// 例: gofetch -u https://example.com --output output.txt
// 例: gofetch -u https://example.com -o output.txt
// 例: gofetch -u https://example.com --timeout 10
// 例: gofetch -u https://example.com -t 10
// 例: gofetch -u https://example.com --retry 5
// 例: gofetch -u https://example.com -r 5
// 例: gofetch -u https://example.com --for 10
// 例: gofetch -u https://example.com -f 10
// 例: gofetch --help
// 例: gofetch -h
// 例: gofetch --version
// 例: gofetch -v
// パラメーターは以下の通り
// -u, --url: アクセスするURLを指定する。必須
// -o, --output: 出力先のファイル名を指定する。省略した場合は標準出力に出力される
// -h, --help: ヘルプを表示する。省略した場合は表示されない
// -v, --version: バージョン情報を表示する。省略した場合は表示されない
// -t, --timeout: タイムアウト時間を指定する。省略した場合は30秒
// -f, --for: 回数を指定する。省略した場合は1回
// -r, --retry: リトライ回数を指定する。省略した場合は3回

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	//バージョン情報
	Version = "1.0.0"
)
const (
	// ヘルプメッセージ
	HelpMessage = `
Usage: gofetch [options]
Options:
  -u, --url     URL to fetch (required)
  -o, --output  Output file (default: stdout)
  -t, --timeout Timeout in seconds (default: 30)
  -r, --retry   Retry count (default: 3)
  -f, --for     Number of times to fetch (default: 1)
  -h, --help    Show this help message
  -v, --version Show version information
`
)

// isValidURL checks if the given URL is valid
func isValidURL(inputURL string) bool {
	_, err := url.ParseRequestURI(inputURL)
	return err == nil
}

// main関数
func main() {
	// コマンドライン引数のパース
	// flagパッケージを使用して、コマンドライン引数をパースする
	url := flag.String("u", "", "URL to fetch")
	output := flag.String("o", "", "Output file (default: stdout)")
	timeout := flag.Int("t", 30, "Timeout in seconds")
	retry := flag.Int("r", 3, "Retry count")
	help := flag.Bool("h", false, "Show help message")
	version := flag.Bool("v", false, "Show version information")

	flag.Parse()

	// ヘルプメッセージの表示
	if *help {
		fmt.Print(HelpMessage)
		os.Exit(0)
	}

	// バージョン情報の表示
	if *version {
		fmt.Println("Version:", Version)
		os.Exit(0)
	}

	// URLが指定されていない場合はエラー
	if *url == "" {
		fmt.Println("Error: URL is required")
		fmt.Print(HelpMessage)
		os.Exit(1)
	}

	// URLのバリデーション
	if !isValidURL(*url) {
		fmt.Println("Error: Invalid URL")
		fmt.Print(HelpMessage)
		os.Exit(1)
	}

	// URLのスキームをhttpに変換
	if !strings.HasPrefix(*url, "http://") && !strings.HasPrefix(*url, "https://") {
		*url = "http://" + *url
	}

	// タイムアウト時間の設定
	client := &http.Client{
		Timeout: time.Duration(*timeout) * time.Second,
	}

	var resp *http.Response
	var err error

	for i := 0; i < *retry; i++ {
		resp, err = client.Get(*url)
		if err == nil {
			break
		}
		time.Sleep(time.Second) // リトライまで1秒待つ
	}

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if *output == "" {
		fmt.Println(string(body))
	} else {
		err = ioutil.WriteFile(*output, body, 0644)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}
}
