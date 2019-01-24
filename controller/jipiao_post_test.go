package controller

import (
	"bao_server/db"
	"log"
	"testing"
)

func Test_InitDB(t *testing.T) {
	if db.InitDB() { //try a unit test on function
		t.Log("connect to db success") // 如果不是如预期的那么就报错
	} else {
		t.Error("connect to db success") //记录一些你期望记录的信息
	}
}
func Test_Insert111(t *testing.T) {
	t1, _ := ParseTime("02-02")
	t2, _ := ParseTime("02-05")
	var one JipiaoInfo
	one.weixinId = "weixintest1"
	one.selfCity = "北京"
	one.selfArrive = "上海"
	one.selfTime = uint64(t1)
	one.expectCity = "深圳"
	one.expectArrive = "成都"
	one.expectTime = uint64(t2)
	one.liaotianbaoID = "llll1"
	sql1 := constructInserSql(&one)
	log.Println(sql1)
	db.Insert(&sql1)

	t1, _ = ParseTime("02-05")
	t2, _ = ParseTime("02-02")
	var two JipiaoInfo
	two.weixinId = "weixintest2"
	two.selfCity = "深圳"
	two.selfArrive = "成都"
	two.selfTime = uint64(t1)
	two.expectCity = "北京"
	two.expectArrive = "上海"
	two.expectTime = uint64(t2)
	two.liaotianbaoID = "llll2"
	sql1 = constructInserSql(&two)
	log.Println(sql1)
	db.Insert(&sql1)

	t1, _ = ParseTime("02-14")
	t2, _ = ParseTime("02-20")
	one.weixinId = "weixintest3"
	one.selfCity = "广州"
	one.selfArrive = "杭州"
	one.selfTime = uint64(t1)
	one.expectCity = "深圳"
	one.expectArrive = "哈尔滨"
	one.expectTime = uint64(t2)
	one.liaotianbaoID = "llll3"
	sql1 = constructInserSql(&one)
	log.Println(sql1)
	db.Insert(&sql1)
}
