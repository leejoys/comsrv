package storage

// Comment - комментарий.
type Comment struct {
	ID            int    `json:"ID"`            // номер записи
	Author        string `json:"Author"`        // автор комментария
	Content       string `json:"Content"`       // содержание комментария
	PubTime       int64  `json:"PubTime"`       //время комментария для БД и фронта
	ParentPost    int    `json:"ParentPost"`    // ID родительской новости
	ParentComment int    `json:"ParentComment"` // ID родительского комментария
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	CommentsN(int) ([]Comment, error) // получение n последних комментариев
	AddComment(Comment) error         // создание нового комментария
	UpdateComment(Comment) error      // обновление комментария
	DeleteComment(Comment) error      // удаление комментария по ID
	Close()                           // освобождение ресурса
}
