###加权平均负载均衡统计

```cassandraql
package main

import (
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/buy/ticket", handleReq)
	http.ListenAndServe(":3004", nil)
}

//处理请求函数,根据请求将响应结果信息写入日志
func handleReq(w http.ResponseWriter, r *http.Request) {
	failedMsg :=  "handle in port:"
	writeLog(failedMsg, "./stat.log")
}

//写入日志
func writeLog(msg string, logPath string) {
	fd, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer fd.Close()
	content := strings.Join([]string{msg, "\r\n"}, "3004")
	buf := []byte(content)
	fd.Write(buf)
}
```

####封装响应结构体，ab并发测试代码
```cassandraql
package main

import (
	"net/http"
	"spikeSystem/util"
)

func main() {
	http.HandleFunc("/buy/ticket", handleReq)
	http.ListenAndServe(":3005", nil)
}

//处理请求函数,根据请求将响应结果信息写入日志
func handleReq(w http.ResponseWriter, r *http.Request) {
	util.RespJson(w, -1, "已售罄", nil)
}
```

### go语言整数形原子性加减
```cassandraql
package main

import (
	"sync"
	"sync/atomic"
)

func main()  {
	var wg sync.WaitGroup
	var counter int64
	counter = 0
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}

	wg.Wait()
	println(counter)
}
```