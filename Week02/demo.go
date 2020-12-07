// business --> service --> dao
package main

import (
	"database/sql"
	"log"

	"github.com/pkg/errors"
)

var errMsg = map[error]string{
	sql.ErrNoRows: "no such data",
}

func dao(id int) (int, error) {
	if id < 0 {
		// return 0, errors.Wrapf(sql.ErrNoRows, "cannot find data with id %d", id)
		// 将原始错误上报，上下层有依赖，好的做法是使用自定义的错误对上层屏蔽技术细节，比如
		// Better one
		return 0, errors.Wrapf(code.NotFound, fmt.Sprintf("sql: %s error: %v", sql, err)

	}
	return id, nil
}

func service(id int) (int, error) {
	return dao(id)
}

func business(id int) string {
	var (
		res int
		err error
	)

	if res, err = service(id); err != nil {
		// 1. 打印详细的堆栈日志
		// 2. 将提示信息返回给最终调用者
		log.Printf("%+v", err)
		return errMsg[errors.Cause(err)]
	}

	log.Printf("result is %d", res)
	return "ok"
}

func main() {
	log.Println(business(2))  // "ok"
	log.Println(business(-1)) // "no such data"
}
