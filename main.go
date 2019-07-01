package main

import (
	"fmt"
	"github.com/pruknil/ydbwrapper/orm"
	"github.com/pruknil/ydbwrapper/ydb"
)

func main() {
	o := orm.NewOrm()
	profile := ydb.Profile{}
	profile.Age = 30

	user := ydb.User{}
	user.Profile = &profile
	user.Name = "slene"

	fmt.Println(o.Insert(profile))
	fmt.Println(o.Insert(user))
}
