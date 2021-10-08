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

//Posts - получение всех публикаций
func (s *Store) Comments() ([]storage.Comment, error) {
	rows, err := s.db.Query(context.Background(),
		`SELECT 
	posts.id, 
	posts.autor, 
	posts.content, 
	posts.pubtime, 
	posts.parentpost,
	posts.parentcomment
	FROM posts;`)

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

//PostsN - получение N публикаций
func (s *Store) CommentsN(n int) ([]storage.Comment, error) {
	rows, err := s.db.Query(context.Background(),
		`SELECT 
	posts.id, 
	posts.autor, 
	posts.content, 
	posts.pubtime, 
	posts.parentpost,
	posts.parentcomment
	FROM posts;`)

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

//AddPost - создание новой публикации
func (s *Store) AddComment(p storage.Comment) error {
	_, err := s.db.Exec(context.Background(), `
	INSERT INTO posts (
		author, 
		content, 
		pubtime, 
		parentpost,
		parentcomment) 
	VALUES ($1,$2,$3,$4,$5);`,
		p.Author,
		p.Content,
		p.PubTime,
		p.ParentPost,
		p.ParentComment)
	return err
}

//UpdatePost - обновление по id значения title,  pubdate, pubtime, и link
func (s *Store) UpdateComment(p storage.Comment) error {
	_, err := s.db.Exec(context.Background(), `
	UPDATE posts 
	SET author=$2,
	content=$3,
	pubtime=$4,
	parentpost=$5,
	parentcomment=$6
	WHERE id=$1;`,
		p.ID,
		p.Author,
		p.Content,
		p.PubTime,
		p.ParentPost,
		p.ParentComment)
	return err
}

//DeletePost - удаляет пост по id
func (s *Store) DeleteComment(p storage.Comment) error {
	_, err := s.db.Exec(context.Background(), `
	DELETE FROM posts 
	WHERE id=$1;`, p.ID)
	return err
}
