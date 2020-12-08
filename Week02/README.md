# Error
1. 预定义Error
```golang
    NOTFOUNT := errors.New("bufio: not_found")

```
2. errors.New返回的是指针，为了确保是同一个error，而不止是string相同。
3. 只在关心这个err（比如logger、其他处理）时，使用这个error，使用后不应该继续传递error到上层。
```golang
    package main

    import (
        "database/sql"

        _ "github.com/go-sql-driver/mysql"
        "github.com/pkg/errors"
    )

    type User struct {
        Id   int
        Name string
    }

    // 预定的错误类型
    var NOTFOUND = errors.New("not_found")

    // dao
    func getUser(name string) (*User, error) {
        db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/monitoring")
        if err != nil {
            // 封装预定义的错误类型
            if err == sql.ErrNoRows {
                return nil, errors.WithStack(NOTFOUND)
            } else {
                return nil, errors.WithStack(err)
            }
        }
        defer db.Close()

        user := &User{}

        err = db.QueryRow("select id, name from users where name = ? ", name).Scan(&user.Id, &user.Name)
        if err != nil {
            return user, errors.WithStack(err)
        }
        return user, nil
    }

    //controller
    func getUserByName(name string) (*User, error) {
        return getUser(name)
    }

    //api
    func api() (int, *User) {
        user, err := getUserByName("abc")
        if err != nil {
            // logger.Error("api_fail, err="+err.Error())
            // 比对预定义的错误类型
            if errors.Is(err, NOTFOUND) {
                return 404, user
            }

            return 500, user
        }

        return 200, user
    }

    func main() {
        api()
    }


```
4. 由包自己判断错误类型，而不是预定义错误类型，获得解耦合
```golang
    type tem interface{

    }
    func IsNullErr(err error) bool{
        te, ok:= err(.NullErr)
        return ok && te.Null()
    }
```
5. Wrap Error
   1. only handle error once.Handling an error means inspecting the error value, and making a single decision.
   2. github/pkg/errors
   3. 使用 errors.New / errors.Errorf 返回错误
   4. 使用 errors.Wrap/ errors.Wrapf 保存堆栈信息
   5. 使用errors.Is 判断错误类型
   6. 使用errors.As 改变错误类型为指定预定义错误类型

6. go 1.13 UnWrap方法提供了错误的封装