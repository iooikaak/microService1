package user

import (
	"context"

	model "github.com/iooikaak/microService1/database/mysql/user"
)

func (d *Dao) GetUserInfo(ctx context.Context, userID int32) (*model.UserInfo, error) {
	result := &model.UserInfo{}
	err := d.Db.Context(ctx).Model(&model.UserInfo{}).Where("id = ?", userID).First(result).Error
	return result, err
}
