/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-20 11:50:21
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-27 12:12:09
 */

package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	alfheimdbwal "github.com/dj456119/AlfheimDB-WAL"
	"github.com/sirupsen/logrus"
)

var index int64
var wal *alfheimdbwal.AlfheimDBWAL

//curl "http://localhost:12345/single?data=hahaha"
func SingeWrite(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	v := req.Form.Get("data")

	logrus.Info("Request data is ", v)
	t1 := time.Now().UnixNano() / 1e6
	buff := make([]byte, 8+8+len(v))
	lItem := alfheimdbwal.NewLogItemBuff(index, []byte(v), buff, true)
	wal.WriteLog(lItem, buff)
	t2 := time.Now().UnixNano() / 1e6
	index++
	w.Write([]byte(fmt.Sprintf("%d\n", lItem.Index)))
	logrus.Info(fmt.Sprintf("%d", t2-t1))
}

//curl "http://localhost:12345/batch?data=hahaha&count=100"
func BatchWrite(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	v := req.Form.Get("data")
	logrus.Info("Request data is ", v)
	batchCount, _ := strconv.ParseInt(req.Form.Get("count"), 10, 32)
	buff := make([]byte, len(v)*int(batchCount)+16*int(batchCount))
	lists := []*alfheimdbwal.LogItem{}
	pos := 0
	for i := 0; i < int(batchCount); i++ {
		lists = append(lists, alfheimdbwal.NewLogItemBuff(index, []byte(v), buff[pos:], true))
		index++
		pos = pos + 16 + len(v)
	}
	t1 := time.Now().UnixNano() / 1e6
	wal.BatchWriteLog(lists, buff)
	t2 := time.Now().UnixNano() / 1e6
	w.Write([]byte(fmt.Sprintf("%d\n", t2-t1)))
	logrus.Info(fmt.Sprintf("%d", t2-t1))
}

//curl "http://localhost:12345/get?index=3"
func GetLog(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	v := req.Form.Get("index")
	index, _ := strconv.ParseInt(v, 10, 32)
	t1 := time.Now().UnixNano() / 1e6
	result := wal.GetLog(index)
	t2 := time.Now().UnixNano() / 1e6

	logrus.Info(string(result), fmt.Sprintf("%d", t2-t1))
	w.Write(result)
}

//curl "http://localhost:12345/benchmarks?perLength=84&batchCount=100&loop=1000"
func Benchmarks(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	perLength, _ := strconv.ParseInt(req.Form.Get("perLength"), 10, 32)
	batchCount, _ := strconv.ParseInt(req.Form.Get("batchCount"), 10, 32)
	loop, _ := strconv.ParseInt(req.Form.Get("loop"), 10, 32)
	logrus.Info("perlength: ", perLength, " batchCount: ", batchCount, " loop: ", loop)
	data := make([]byte, perLength)
	for i := range data {
		data[i] = 'a'
	}
	t1 := time.Now().UnixNano() / 1e6
	buff := make([]byte, perLength*batchCount+batchCount*16)
	for x := 0; x < int(loop); x++ {

		lists := make([]*alfheimdbwal.LogItem, batchCount)
		pos := 0
		for i := 0; i < int(batchCount); i++ {
			lists = append(lists, alfheimdbwal.NewLogItemBuff(index, data, buff[pos:], true))
			index++
			pos = pos + 16 + int(perLength)
		}
		t3 := time.Now().UnixNano() / 1e6
		wal.BatchWriteLog(lists, buff)
		t4 := time.Now().UnixNano() / 1e6
		fmt.Println(" cost :", (t4 - t3))
	}

	t2 := time.Now().UnixNano() / 1e6

	logrus.Info(fmt.Sprintf("%d\n", t2-t1))
	w.Write([]byte(fmt.Sprintf("%d\n", t2-t1)))
}

//curl "http://localhost:12345/delete?startIndex=10&endIndex=20"
func Delete(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	startIndex, _ := strconv.ParseInt(req.Form.Get("startIndex"), 10, 32)
	endIndex, _ := strconv.ParseInt(req.Form.Get("endIndex"), 10, 32)
	wal.TruncateLog(startIndex, endIndex)
	w.Write([]byte("delete ok\n"))
}

func Init(aw *alfheimdbwal.AlfheimDBWAL) {
	wal = aw
	index = wal.MaxIndex + 1
	logrus.Info("Http Test Server is start")
	http.HandleFunc("/single", SingeWrite)
	http.HandleFunc("/batch", BatchWrite)
	http.HandleFunc("/get", GetLog)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/benchmarks", Benchmarks)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	wal := alfheimdbwal.AlfheimDBWAL{}
	wal.Dirname = "data/"
	wal.MaxItems = 1000
	wal.IsBigEndian = true
	wal.Mutex = new(sync.Mutex)
	wal.BuildDirIndex()
	Init(&wal)

}
