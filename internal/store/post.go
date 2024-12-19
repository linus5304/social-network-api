package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) GetUserFeed(ctx context.Context, userId int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	fSearch := fmt.Sprintf("%%%s%%", fq.Search)
	fmt.Printf("tags: %v", fq)

	qeury := `
		select 
			p.id, p.user_id, p.title, p."content" , p.created_at ,p.updated_at , p."version", 
			count(c.id) as comments_count, p.tags, u.username 
		from posts p 
		left join "comments" c on c.post_id  = p.id
		left join users u on p.user_id = u.id 
		join followers f on f.follower_id = p.user_id or p.user_id = $1
		where 
			f.user_id = $1 and 
			(p.tags @> $5 or array_length($5, 1) = 0) and
			(p.title ilike $4 or p.content ilike $4)
		group by p.id, u.username 
		order by p.created_at ` + fq.Sort + `
		limit $2 offset $3
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, qeury, userId, fq.Limit, fq.OffSet, fSearch, pq.Array(fq.Tags))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []PostWithMetadata
	for rows.Next() {
		var p PostWithMetadata
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.Version,
			&p.CommentsCount,
			pq.Array(&p.Tags),
			&p.User.Username,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, p)
	}

	return feed, nil
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
	insert into posts (content, title, user_id, tags)
	values ($1, $2, $3, $4) returning id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostStore) GetById(ctx context.Context, id int64) (*Post, error) {
	query :=
		`select id, user_id, title, content, created_at, updated_at, tags, version from posts where id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var post Post

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostStore) Delete(ctx context.Context, id int64) error {
	query := `delete from posts where id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return nil
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
	update posts 
	set title = $1, content = $2, version = version + 1
	where id = $3 and version = $4
	returning version`

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
