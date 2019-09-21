package main

import (
	"net/http"
	"os"
	localSpike2 "spikeSystem/localSpike"
	remoteSpike2 "spikeSystem/remoteSpike"
	"spikeSystem/util"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
)

var (
	localSpike localSpike2.LocalSpike
	remoteSpike remoteSpike2.RemoteSpikeKeys
	redisPool *redis.Pool
	done chan int
)

//初始化要使用的结构体和redis连接池
func init() {
	localSpike = localSpike2.LocalSpike{
		LocalInStock:     150,
		LocalSalesVolume: 0,
	}
	remoteSpike = remoteSpike2.RemoteSpikeKeys{
		SpikeOrderHashKey:  "ticket_hash_key",
		TotalInventoryKey:  "ticket_total_nums",
		QuantityOfOrderKey: "ticket_sold_nums",
	}
	redisPool = remoteSpike2.NewPool()
	done = make(chan int, 1)
	done <- 1
}

func main() {
	http.HandleFunc("/buy/ticket", handleReq)
	http.ListenAndServe(":3005", nil)
}

//处理请求函数,根据请求将响应结果信息写入日志
func handleReq(w http.ResponseWriter, r *http.Request) {
	redisConn := redisPool.Get()
	LogMsg := ""
	<-done
	//全局读写锁
	if localSpike.LocalDeductionStock() && remoteSpike.RemoteDeductionStock(redisConn) {
		util.RespJson(w, 1,  "抢票成功", nil)
		LogMsg = LogMsg + "result:1,localSales:" + strconv.FormatInt(localSpike.LocalSalesVolume, 10)
	} else {
		util.RespJson(w, -1, "已售罄", nil)
		LogMsg = LogMsg + "result:0,localSales:" + strconv.FormatInt(localSpike.LocalSalesVolume, 10)
	}
	//将抢票状态写入到log中
	done <- 1
	writeLog(LogMsg, "./stat.log")
}

func writeLog(msg string, logPath string) {
	fd, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer fd.Close()
	content := strings.Join([]string{msg, "\r\n"}, "")
	buf := []byte(content)
	fd.Write(buf)
}
