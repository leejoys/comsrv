package pgdb

import (
	"comsrv/pkg/storage"
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

var ErrorDuplicatePost error = errors.New("SQLSTATE 23505")

// Хранилище данных.
type Store struct {
	db *pgxpool.Pool
}

//New - Конструктор объекта хранилища.
func New(connstr string) (*Store, error) {

	db, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	// проверка связи с БД
	err = db.Ping(context.Background())
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Store{db: db}, nil
}

//Close - освобождение ресурса
func (s *Store) Close() {
	s.db.Close()
}

//Comments - получение всех публикаций
func (s *Store) Comments() ([]storage.Comment, error) {
	rows, err := s.db.Query(context.Background(),
		`SELECT 
		comments.id, 
		comments.autor, 
		comments.content, 
		comments.pubtime, 
		comments.parentpost,
		comments.parentcomment
	FROM comments;`)

	if err != nil {
		return nil, err
	}

	var posts []storage.Comment
	for rows.Next() {
		var p storage.Comment
		err = rows.Scan(
			&p.ID,
			&p.Author,
			&p.Content,
			&p.PubTime,
			&p.ParentPost,
			&p.ParentComment,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, rows.Err()
}

//CommentsN получить комментарии к новости n
func (s *Store) CommentsN(n int) ([]storage.Comment, error) {
	rows, err := s.db.Query(context.Background(),
		`SELECT 
		comments.id, 
		comments.autor, 
		comments.content, 
		comments.pubtime, 
		comments.parentpost,
		comments.parentcomment
	FROM comments
	WHERE parentpost=$1;`, n)

	if err != nil {
		return nil, err
	}

	var comments []storage.Comment
	for rows.Next() {
		var c storage.Comment
		err = rows.Scan(
			&c.ID,
			&c.Author,
			&c.Content,
			&c.PubTime,
			&c.ParentPost,
			&c.ParentComment,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, rows.Err()
}

//AddComment - создание новой публикации
func (s *Store) AddComment(c storage.Comment) error {
	_, err := s.db.Exec(context.Background(), `
	INSERT INTO comments (
		author, 
		content, 
		pubtime, 
		parentpost,
		parentcomment) 
	VALUES ($1,$2,$3,$4,$5);`,
		c.Author,
		c.Content,
		c.PubTime,
		c.ParentPost,
		c.ParentComment)
	return err
}

//UpdateComment - обновление по id значения author,  pubtime, parentpost, и parentcomment
func (s *Store) UpdateComment(c storage.Comment) error {
	_, err := s.db.Exec(context.Background(), `
	UPDATE comments 
	SET author=$2,
	content=$3,
	pubtime=$4,
	parentpost=$5,
	parentcomment=$6
	WHERE id=$1;`,
		c.ID,
		c.Author,
		c.Content,
		c.PubTime,
		c.ParentPost,
		c.ParentComment)
	return err
}

//DeleteComment - удаляет пост по id
func (s *Store) DeleteComment(c storage.Comment) error {
	_, err := s.db.Exec(context.Background(), `
	DELETE FROM comments 
	WHERE id=$1;`, c.ID)
	return err
}
