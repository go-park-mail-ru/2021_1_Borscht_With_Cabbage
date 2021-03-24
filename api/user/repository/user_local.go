package repository

import (
	"backend/api/domain"
	errors "backend/utils"
	"database/sql"
	"fmt"
	"net/http"
)

type userRepo struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepo {
	return &userRepo{DB: db}
}

func NewUserRepo() domain.UserRepo {
	return &userRepo{}
}

func (u *userRepo) Create(newUser domain.User) error {
	// todo exists checking
	err := u.DB.QueryRow("insert into users (name, number, email, password) values ($1, $2, $3, $4)",
		newUser.Name, newUser.Phone, newUser.Email, newUser.Password)
	if err != nil {
		// TODO
	}

	return nil
}

// для авторизации
func (u *userRepo) CheckUserExists(userToCheck domain.UserAuth) (domain.User, error) {
	DBuser, err := u.DB.Query("select uid, name, avatar from users where (phone=&1 or email=$1) and password=$2",
		userToCheck.Login, userToCheck.Password)
	if err != nil {
		// todo
	}
	user := new(domain.User)
	for DBuser.Next() {
		err = DBuser.Scan(
			&user.Uid,
			&user.Name,
			&user.Avatar,
		)
		if err != nil {
			// TODO
		}
	}

	return *user, nil
}

func (u *userRepo) GetByUid(uid int32) (domain.User, error) {
	DBuser, err := u.DB.Query("select name, phone, email, avatar, address, mainAddress from users where uid=$1", uid),
	if err != nil {
		return domain.User{}, errors.Authorization("user not found")
	}
	user := new(domain.User)
	for DBuser.Next() {
		err = DBuser.Scan(
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.MainAddress,
			&user.Avatar,
		)
		if err != nil {
			// TODO
		}
	}
	return *user, nil
}

func (u *userRepo) Update(newUser domain.UserData) error {
	// todo

	//for i, user := range *ctx.Users {
	//	if user.Email == newUser.Email && user.Phone != ctx.User.Phone { // если у кого-то другого уже есть такой email
	//		return errors.Create(http.StatusBadRequest, "user with this email already exists")
	//	}
	//	if user.Phone == newUser.Phone && user.Phone != ctx.User.Phone { // если у кого-то другого уже есть такой телефон
	//		return errors.Create(http.StatusBadRequest, "user with this number already exists")
	//	}
	//
	//	if user.Phone == ctx.User.Phone {
	//		if newUser.Password != "" {
	//			if newUser.PasswordOld != user.Password {
	//				fmt.Println(newUser.PasswordOld, " ", user.Password)
	//				return errors.Create(http.StatusBadRequest, "invalid old password")
	//			}
	//			(*ctx.Users)[i].Password = newUser.Password
	//		}
	//
	//		(*ctx.Users)[i].Phone = newUser.Phone
	//		(*ctx.Users)[i].Email = newUser.Email
	//		(*ctx.Users)[i].Name = newUser.Name
	//
	//		fmt.Println(*ctx.Users)
	//		//return cc.SendOK(profileEdits)
	//		return nil
	//	}
	//}

	return errors.Authorization("user not found")
}
