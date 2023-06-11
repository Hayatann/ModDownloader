package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"golang.org/x/sync/semaphore"
)
func main(){
    urls := []string{
        "https://taikosanjiro-humenroom.net/wp-content/uploads/shinryu-ura.zip",
        "https://taikosanjiro-humenroom.net/wp-content/uploads/vixtory-ura.zip",
    }
    downloadParallel(urls)
}
func downloadParallel(urls []string) {
    var wg sync.WaitGroup
    var s = semaphore.NewWeighted(5) // 同時実行するgoroutineの数を指定
    for _, u := range urls {
        wg.Add(1)
        go downloadFromURL(u, &wg, s)
    }
    wg.Wait() // goroutineに投げた全ての処理が終了するまで待機
}
func downloadFromURL(_url string, wg *sync.WaitGroup, s *semaphore.Weighted) {
    defer wg.Done() // この関数の実行が終了したことをsync.WaitGroupに伝える
    if err := s.Acquire(context.Background(), 1); err != nil { // セマフォを1つロック
        return
    }
    defer s.Release(1) // この関数の実行完了時にセマフォを1つ解放
    u, err := url.Parse(_url)
    if err != nil {
        log.Fatal(err)
    }
    path := u.Path
    segments := strings.Split(path, "/")
    fmt.Println(segments) //
    fileName := segments[len(segments)-1] // URLからファイル名を作成
    fileName = "dist/" + fileName
    if f, err := os.Stat("dist"); os.IsNotExist(err) || !f.IsDir() {
        os.Mkdir("dist", 0753) // ディレクトリがなければ作成
    }
    file, err := os.Create(fileName) // ファイル名からファイルを作成;
    if err != nil {
        log.Fatal(err)
    }
    resp, err := http.Get(_url) // URLからレスポンスを取得
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    size, err := io.Copy(file, resp.Body) // 取得したレスポンスをファイルにコピー
    fmt.Printf("%s %dKB\n", fileName, size/1024)
    defer file.Close()
}
