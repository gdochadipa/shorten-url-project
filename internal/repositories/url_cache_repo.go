package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ochadipa/url-shorterner-project/internal/db"
	"github.com/ochadipa/url-shorterner-project/internal/model"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (u *UrlRepository) StoreCacheUrl(ctx context.Context, urlModel *model.Url) error {

	jsonUrl, err := json.Marshal(urlModel)
	if err != nil {
		zap.L().Error("failed.parse.to.json", zap.Error(err))
		return err
	}

	key := fmt.Sprintf("url:%s", urlModel.ID)
	ttl := 24 * time.Hour

	err = db.Rds.Set(ctx, key, jsonUrl, ttl).Err()

	if err != nil {
		zap.L().Error("failed.set.to.redis", zap.Error(err))
		return err
	}

	zap.L().Info("successfully.set.to.redis")
	return nil
}

func (u *UrlRepository) GetCacheUrl(ctx context.Context, id string) (*model.Url, error) {
	zap.L().Info("get.from.redis")
	key := fmt.Sprintf("url:%s", id)
	value, err := db.Rds.Get(ctx, key).Result()

	if err == redis.Nil {
		zap.L().Info("not.found.redis")
		return nil, nil
	} else if err != nil {
		zap.L().Error("found.error", zap.Error(err))
		return nil, err
	}

	var url *model.Url
	err = json.Unmarshal([]byte(value), &url)

	if err != nil {
		zap.L().Error("failed.set.to.redis", zap.Error(err))
		return nil, err
	}

	return url, nil

}
