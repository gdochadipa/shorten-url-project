package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ochadipa/url-shorterner-project/internal/model"
	"github.com/ochadipa/url-shorterner-project/internal/repositories"
	"github.com/ochadipa/url-shorterner-project/pkg"
	"go.uber.org/zap"
)

type UrlService struct {
	UrlRepo *repositories.UrlRepository
}

func (s *UrlService) StoreUrl(ctx context.Context, url string) (*model.Url, error) {
	/**
	 * - validation
		*  - store on db
		* - store on redis
	*/

	if err := pkg.ValidateURL(url); err != nil {
		zap.L().Error("invalid.url", zap.Error(err))
		return nil, err
	}

	shortUrl, err := s.UrlRepo.StoreUrl(ctx, url)

	if err != nil {
		zap.L().Error("failed.insert", zap.Error(err))
		return nil, err
	}

	urlModel := &model.Url{
		ID:        *shortUrl,
		URL:       url,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = s.UrlRepo.StoreCacheUrl(ctx, urlModel)
	if err != nil {
		zap.L().Error("failed.insert.to.redis", zap.Error(err))
		return nil, err
	}

	return urlModel, nil
}

func (s *UrlService) GetUrl(ctx context.Context, id string) (*model.Url, error) {
	/**
	 * - check on redis
	 * - check on db
	 * - if exists on db but not on redis, store on redis
	 * - if not exists, then return not found
	 */

	if urlModel, err := s.UrlRepo.GetCacheUrl(ctx, id); urlModel != nil {
		return urlModel, nil
	} else if err != nil {
		zap.L().Error("failed.get.from.redis", zap.Error(err))
	}

	urlModel, err := s.UrlRepo.GetUrl(ctx, id)
	if err != nil {
		zap.L().Error("failed.get.from.databases", zap.Error(err))
	}
	if urlModel != nil {
		zap.L().Info("what is this?", zap.Any("urlmodel",urlModel))
		err = s.UrlRepo.StoreCacheUrl(ctx, urlModel)
		if err != nil {
			zap.L().Error("failed.insert.to.redis", zap.Error(err))
			return nil, err
		}
		return urlModel, nil
	}

	return nil, fmt.Errorf("not.found")
}


func NewService(repo *repositories.UrlRepository)*UrlService{
	return &UrlService{
		UrlRepo: repo,
	}
}
