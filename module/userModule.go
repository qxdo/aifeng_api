package module

import (
	"errors"
	"fmt"
	"xiaoniu/common"
	"xiaoniu/utils"
)

func Login(user common.User) (*common.User, error) {
	sqlCmd := fmt.Sprintf("select * from %s where username=? and password=? and disabled=0", common.TableUser)
	row := utils.DB.QueryRowx(sqlCmd, user.Username, user.Password)
	if row == nil {
		return nil, errors.New("login row is empty")
	}
	u := &common.User{}
	err := row.StructScan(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}
