package repositories

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/ochadipa/url-shorterner-project/internal/db"
	"github.com/ochadipa/url-shorterner-project/internal/model"
	"github.com/ochadipa/url-shorterner-project/pkg"
	"go.uber.org/zap"
)

type UrlRepository struct {
	url *model.Url
}

func (u *UrlRepository) StoreUrl(ctx context.Context,uri  string) (*string, error) {
	tableName := u.url.GetTableName()
	now := time.Now()
	randomBytes := make([]byte, 6)

	const maxRetries = 5 // Maximum attempts to generate a unique short URL
	var shortURL string

		query := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	for i := range maxRetries {
		fmt.Printf("iteraition %d",i)
		n, err := rand.Read(randomBytes)

		if err != nil {
			return nil, fmt.Errorf(" failed to get random: %w", err)
		}

		if n != len(randomBytes) {
			fmt.Errorf("expected to read %d bytes, but got %d", len(randomBytes), n)
			continue
		}

		shortURL = pkg.RandomString(5)
		sql, args, err := query.Insert(tableName).
			Columns("id", "url", "created_at", "updated_at").
			Values(shortURL, uri, now, now).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			return  nil, fmt.Errorf("failed to build SQL query for short URL generation: %w", err)
		}

		tx, err := db.Db.Begin()
		if err != nil {
			return nil,  fmt.Errorf("failed to begin database transaction: %w", err)
		}


		_, err = tx.ExecContext(ctx, sql, args...)

		if err != nil {
			tx.Rollback()

			if strings.Contains(err.Error(), "duplicate key") ||
				strings.Contains(err.Error(), "unique constraint") ||
				strings.Contains(err.Error(), "UNIQUE constraint failed") {
				continue
			}
			return nil, fmt.Errorf("failed to execute DB insert for short URL: %w", err)
		}

		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to commit DB transaction for short URL: %w", err)
		}
		return &shortURL, nil
	}

	return nil, fmt.Errorf("failed to generate unique short URL after %d retries. All attempts resulted in a collision.", maxRetries)
}

func (u *UrlRepository) GetUrl(ctx context.Context, id string) (*model.Url, error) {
	tableName := u.url.GetTableName()

	query := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.Select("*").From(tableName).Where(squirrel.Eq{"id": id}).Limit(1).ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build SQL select url: %w", err)
	}

	row, err := db.Db.QueryContext(ctx, sql, args...)

	if err != nil {
		return nil, fmt.Errorf("failed get the url: %w", err)
	}

	defer row.Close()
	urls, err := u.scanRows(row)

	if err != nil || len(urls) > 1 {
		return nil, fmt.Errorf("failed patch the url: %w", err)
	}

	if len(urls) == 0 {
		zap.L().Info("arguments", zap.String("sql", sql), zap.Any("args", args))
		zap.L().Error("not.found.url.sql",zap.Error(err))
		return nil, fmt.Errorf("url is not found: %w", err)
	}

	return urls[0], nil
}

func (u *UrlRepository) DeleteUrl(ctx context.Context, id string) error {
	tableName := u.url.GetTableName()

	sql, args, err := sq.Delete(tableName).Where(sq.Eq{"id": id}).ToSql()

	if err != nil {
		return fmt.Errorf(" failed to delete url: %w", err)
	}

	tx, err := db.Db.Begin()
	if err != nil {
		return fmt.Errorf(" failed to  begin delete url: %w", err)
	}

	_, err = tx.ExecContext(ctx, sql, args...)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf(" executed delte url: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf(" failed to commit: %w", err)
	}
	return nil

}

func (u *UrlRepository) scanRows(r *sql.Rows) ([]*model.Url, error) {
	var urls []*model.Url
	for r.Next() {
		var id, shortUrl string
		var createdAt, updatedAt time.Time

		err := r.Scan(&id, &shortUrl, &createdAt, &updatedAt)

		if err != nil {
			return nil, err
		}

		url := &model.Url{
			ID:        id,
			URL:       shortUrl,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		urls = append(urls, url)
	}

	return urls, nil
}


func NewUrlRepo (urlModel *model.Url) *UrlRepository{
	return &UrlRepository{
		url: urlModel,
	}
}
