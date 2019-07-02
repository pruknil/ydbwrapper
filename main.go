package main

import (
	"fmt"
	"github.com/pruknil/ydbwrapper/orm"
	"github.com/pruknil/ydbwrapper/ydb"
)

func main() {
	o := orm.NewOrm()
	profile := ydb.Profile{
		Id:   0,
		Age:  30,
		User: nil,
	}

	user := ydb.User{
		Id:      0,
		Name:    "slene",
		Profile: &profile,
	}

	fmt.Println(o.Insert(profile))
	fmt.Println(o.Insert(user))
}
