package module

import (
	"fmt"
	"xiaoniu/common"
	"xiaoniu/utils"
)

func AddPriceFormDB(p common.Price) error {
	sqlCmd := "insert into " + common.TablePrice + " (day_price,night_price,uid,`time`) " +
		" values (?,?,?,?) "
	_, err := utils.Execute(sqlCmd, p.DayPrice, p.NightPrice, p.Uid, p.SqlTime)
	if err != nil {
		return err
	}
	return nil
}

func GetPriceFormDB(uid int) (*common.Price, error) {
	db := utils.DB
	p := &common.Price{}
	sqlCmd := fmt.Sprintf("select * from %s where uid = ? ORDER BY time DESC", common.TablePrice)
	proxyRow := db.QueryRowx(sqlCmd, uid)
	err := proxyRow.StructScan(p)
	if err != nil {
		if err.Error() == common.NoRowsError {
			return p, nil
		}
		return nil, err
	}
	return p, nil
}
