package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hellobchain/gotool/gcli"
	"github.com/hellobchain/gotool/gcrypto"
	"github.com/hellobchain/gotool/gexcel"
	"github.com/hellobchain/gotool/gschedule"
	"github.com/hellobchain/gotool/gstr"
	"github.com/hellobchain/gotool/gtree"
	"github.com/hellobchain/gotool/gvalid"
)

func main() {
	fmt.Println(gstr.CamelCase("hello_world")) // helloWorld
	fmt.Println(gcrypto.MD5("123"))            // 202cb962ac59075b964b07152d234b70
	gcli.PrintSuccess("部署完成！")
	// 校验器
	if err := gvalid.New("13800138000").Mobile().Check(); err != nil {
		panic(err)
	}

	// 导出 Excel
	headers := []string{"姓名", "年龄", "邮箱"}
	rows := [][]interface{}{
		{"Alice", 18, "alice@example.com"},
		{"Bob", 20, "bob@example.com"},
	}
	_ = gexcel.WriteCSV("users.csv", headers, rows)   // 直接用 Excel 打开
	_ = gexcel.WriteXLSX("users.xlsx", headers, rows) // 真正 xlsx
	// 打印当前目录（默认）
	_ = gtree.Print()

	// 打印指定目录，仅目录、显示大小、最大 3 层
	opt := gtree.Option{
		DirOnly:  true,
		ShowSize: true,
		MaxDepth: 3,
	}
	_ = gtree.PrintDir("/tmp", opt)

	// 获取树字符串
	st, _ := gtree.String(".", gtree.DefaultOption)
	println(st)

	s := gschedule.New(10) // 最多 10 并发

	// 1. 每 2 秒一次
	s.Add("tick", gschedule.Every(2*time.Second), gschedule.JobFunc(func() {
		log.Println("tick", time.Now().Format("15:04:05"))
	}))

	// 2. cron：每周一到五 09:00
	cron, _ := gschedule.NewCron("0 9 * * 1-5")
	s.Add("job9am", cron, gschedule.JobFunc(func() {
		log.Println("good morning 9am")
	}))

	// 3. 单次延迟
	s.Add("once", gschedule.Delay(3*time.Second), gschedule.JobFunc(func() {
		log.Println("3s later run once")
	}))

	// 阻塞 30 秒后优雅退出
	time.Sleep(30 * time.Second)
	s.Stop()
}
