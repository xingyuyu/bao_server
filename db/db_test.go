package db

import (
	"testing"
)

func Test_InitDB(t *testing.T) {
	if InitDB() { //try a unit test on function
		t.Log("connect to db success") // 如果不是如预期的那么就报错
	} else {
		t.Error("connect to db success") //记录一些你期望记录的信息
	}
}

func Test_Insert(t *testing.T) {
	sql := "insert into jipiao_exchange(liaotianbao_id, self_city, self_arrive,self_time,expect_city,expect_arrive,expect_time,update_time,status) values('yuxingyu','北京','南京',1548048779,'广州', '上海',1548048779,1548048770,0)"
	ret, err := Insert(&sql)
	if err != nil {
		t.Error("insert fail")
	}
	if ret != 1 {
		t.Error("insert num fail")
	}
	t.Log("insert success")

}
