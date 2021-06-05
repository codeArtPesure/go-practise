package error_handle

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
)

type User struct {
    Id int64
    Name string
}

type UserDao interface {
    OneById(ctx context.Context,id int64) (*User,error)
}

type UserDaoImpl struct {
    db *sql.DB
}

func (u *UserDaoImpl)OneById(ctx context.Context,pk int64)(*User,error){
    var (
      id   int64
      name string
    )
    rows, err := u.db.Query(`SELECT id,name From user where id = ?`,pk)
    if err != nil {
      return nil,fmt.Errorf("UserDao OneById failed %w",err)
    }
    defer rows.Close()

    for rows.Next() {
      err=rows.Scan(&id, &name)
      if err != nil {
          return nil,fmt.Errorf("UserDao OneById failed %w",err)
      }
    }
    return &User{Id:id,Name:name},nil
}

func NewUserDao( db *sql.DB) *UserDaoImpl{
    return &UserDaoImpl{
        db:db,
    }
}

type UserServiceImpl struct {
    context context.Context
    dao UserDao
}

func NewUserService( context context.Context,dao UserDao) *UserServiceImpl{
    return &UserServiceImpl{
        context:context,
        dao:dao,
    }
}

func (s *UserServiceImpl) GetUser(ctx context.Context,userId int64)(*User,error){
    user, err := s.dao.OneById(ctx, userId)
    if err!=nil&&errors.Is(err,sql.ErrNoRows){   // 数据库无该用户
        return nil,&CustomError{1,"用户不存在"}
    }
    if err!=nil{
        return nil,fmt.Errorf("userservice.GetUser err : %w",err)
    }
    return user,nil
}


