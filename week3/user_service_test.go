package week3

import (
    "context"
    "database/sql"
    "log"
    "testing"

    _ "github.com/go-sql-driver/mysql"
)

func TestGetUser(t *testing.T){
    db, err := sql.Open("mysql",
        "root:123456@tcp(127.0.0.1:3306)/hello")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    ctx:=context.Background()
    dao:=NewUserDao(db)
    service := NewUserService(ctx, dao)
    user, err := service.GetUser(ctx, 1)
    if err!=nil{
        t.Fatal(err)
    }
    t.Log(user)
}