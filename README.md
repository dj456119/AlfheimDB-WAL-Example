<!--
 * @Descripttion: 
 * @version: 
 * @Author: cm.d
 * @Date: 2021-11-26 21:06:03
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-26 21:22:34
-->

# AlfheimDB-WAL-Example

A AlfheimDB-WAL Example.

# Default Url

+ Write: `curl "http://localhost:12345/single?data=hahaha"`
+ BatchWrite: `curl "http://localhost:12345/batch?data=hahaha&count=100"`
+ Get: `curl "http://localhost:12345/get?index=3"`
+ Benchmarks: `curl "http://localhost:12345/benchmarks?perLength=84&batchCount=100&loop=1000"`
+ Delete: `curl "http://localhost:12345/delete?startIndex=10&endIndex=20"`
