# SuperCrawler

## 使用方法

- 编译
```go
go build main.go -o supercrawler
```

- 运行
```go
./supercrawler [options] [value]
```

两种模式：
1. update: 只爬取最近的一页（更新）
2. all：   全量爬取（历史）

```text
Usage of supercrawler:
  -m string
        Update / All (default "update")
  -t string
        target to crawl, list of register func, use ',' to split, like 'aa,bb,cc'
  -v    show version and exit.
```

- plugin: 插件调用总函数
- core：插件调用细节