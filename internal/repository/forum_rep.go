package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DrusGalkin/forum-client/internal/entity"
	"go.uber.org/zap"
	"time"
)

type ForumRepository interface {
	GetAllThreads() ([]entity.Thread, error)
	GetThreadByID(id int) (entity.Thread, error)
	CreateThread(thread entity.Thread) (entity.Thread, error)
	DeleteThreadByID(id int) error
	GetThreadsByUserID(userId int) ([]entity.Thread, error)
	CreatePost(post entity.Post) (entity.Post, error)
	GetPostsByThreadID(threadID int) ([]entity.Post, error)
	DeletePostByID(id int) error
	GetPostsByUserID(id int) ([]entity.Post, error)
	GetChatPosts(threadID int) ([]entity.Post, error)
	LinkPostToChat(chat entity.Chat) error
	CheckUserByID(user entity.User, id int) (bool, error)
	GetPostByID(id int) (entity.Post, error)
	EditThread(thread entity.Thread, userID int) error
}

type forumRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewForumRepository(db *sql.DB, logger *zap.Logger) ForumRepository {
	return &forumRepository{
		db:     db,
		logger: logger,
	}
}

func (f *forumRepository) GetAllThreads() ([]entity.Thread, error) {
	f.logger.Info("Получение всех тредов по ID")
	query := `SELECT id, title, content, create_at, user_id 
			  FROM threads 
			  ORDER BY create_at DESC`
	rows, err := f.db.Query(query)
	if err != nil {
		f.logger.Error("Ошибка получения тредов", zap.Error(err))
		return nil, fmt.Errorf("Ошибка получения тредов: %w", err)
	}
	defer rows.Close()

	var threads []entity.Thread
	for rows.Next() {
		var thread entity.Thread
		if err := rows.Scan(&thread.ID, &thread.Title, &thread.Content, &thread.CreateAt, &thread.UserID); err != nil {
			f.logger.Fatal(fmt.Sprintf("Ошибка поиска треда c ID = %d", thread.ID),
				zap.Error(err),
				zap.Any("Thread", thread))
			return nil, fmt.Errorf("Ошибка поиска треда: %w", err)
		}
		threads = append(threads, thread)
	}
	f.logger.Info("Успешное получение всех тредов")
	return threads, nil
}

func (f *forumRepository) GetThreadByID(id int) (entity.Thread, error) {
	f.logger.Debug("Получение треда по ID", zap.Int("id", id))
	query := `SELECT id, title, content, create_at, user_id 
              FROM threads 
              WHERE id = $1
              ORDER BY create_at DESC`

	var thread entity.Thread
	err := f.db.QueryRow(query, id).Scan(
		&thread.ID,
		&thread.Title,
		&thread.Content,
		&thread.CreateAt,
		&thread.UserID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			f.logger.Warn("Тред не найден", zap.Int("id", id))
			return entity.Thread{}, entity.ErrorNotFoundThread
		}
		f.logger.Error("Ошибка при получении треда по ID",
			zap.Int("id", id),
			zap.Error(err))
		return entity.Thread{}, fmt.Errorf("Ошибка поиска треда по id: %w", err)
	}

	f.logger.Debug("Тред успешно получен",
		zap.Int("id", id),
		zap.Any("thread", thread))
	return thread, nil
}

func (f *forumRepository) CreateThread(thread entity.Thread) (entity.Thread, error) {
	f.logger.Debug("Создание нового треда",
		zap.String("title", thread.Title),
		zap.Int("userID", thread.UserID))

	query :=
		`INSERT INTO threads (title, content, create_at, user_id)
         VALUES ($1, $2, $3, $4)
         RETURNING id, title, content, create_at, user_id`

	var createThread entity.Thread
	err := f.db.QueryRow(
		query,
		thread.Title,
		thread.Content,
		time.Now(),
		thread.UserID,
	).Scan(
		&createThread.ID,
		&createThread.Title,
		&createThread.Content,
		&createThread.CreateAt,
		&createThread.UserID,
	)

	if err != nil {
		f.logger.Error("Ошибка при создании треда",
			zap.Any("thread", thread),
			zap.Error(err))
		return entity.Thread{}, fmt.Errorf("Ошибка при создании треда: %w", err)
	}

	f.logger.Info("Тред успешно создан",
		zap.Int("id", createThread.ID),
		zap.String("title", createThread.Title))
	return createThread, nil
}

func (f *forumRepository) EditThread(thread entity.Thread, userID int) error {
	query := `UPDATE threads
			  SET title=$1, content=$2, create_at=$3
			  WHERE id=$4;`

	valid, err := f.CheckUserByID(thread, userID)
	if err != nil {
		return fmt.Errorf("ошибка проверки прав: %w", err)
	}
	if !valid {
		return fmt.Errorf("нет прав на редактирование")
	}

	exec, err := f.db.Exec(query, thread.Title, thread.Content, thread.CreateAt, thread.ID)
	if err != nil {
		return err
	}

	affected, err := exec.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (f *forumRepository) DeleteThreadByID(id int) error {
	f.logger.Debug("Удаление треда по ID", zap.Int("id", id))
	query := `DELETE FROM threads WHERE id = $1`
	result, err := f.db.Exec(query, id)
	if err != nil {
		f.logger.Error("Ошибка при удалении треда",
			zap.Int("id", id),
			zap.Error(err))
		return fmt.Errorf("Ошибка удаления треда: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		f.logger.Error("Ошибка при получении количества удаленных строк",
			zap.Int("id", id),
			zap.Error(err))
		return fmt.Errorf("Ошибка получения измененых строк: %w", err)
	}

	if rowsAffected == 0 {
		f.logger.Warn("Тред не найден для удаления", zap.Int("id", id))
		return entity.ErrorNotFoundThread
	}

	f.logger.Info("Тред успешно удален",
		zap.Int("id", id),
		zap.Int64("rowsAffected", rowsAffected))
	return nil
}

func (f *forumRepository) CreatePost(post entity.Post) (entity.Post, error) {
	f.logger.Debug("Создание нового поста",
		zap.Int("threadID", post.ThreadID),
		zap.Int("userID", post.UserID))

	query :=
		`INSERT INTO posts (content, create_at, thread_id, user_id)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, content, create_at, thread_id, user_id`

	var createdPost entity.Post
	err := f.db.QueryRow(
		query,
		post.Content,
		post.CreateAt,
		post.ThreadID,
		post.UserID,
	).Scan(
		&createdPost.ID,
		&createdPost.Content,
		&createdPost.CreateAt,
		&createdPost.ThreadID,
		&createdPost.UserID,
	)
	if err != nil {
		f.logger.Error("Ошибка при создании поста",
			zap.Any("post", post),
			zap.Error(err))
		return entity.Post{}, fmt.Errorf("Ошибка при создании поста: %w", err)
	}

	f.logger.Info("Пост успешно создан",
		zap.Int("id", createdPost.ID),
		zap.Int("threadID", createdPost.ThreadID))
	return createdPost, nil
}
func (f *forumRepository) GetPostsByThreadID(threadID int) ([]entity.Post, error) {
	f.logger.Debug("Получение постов по ID треда", zap.Int("threadID", threadID))
	query :=
		`SELECT id, content, create_at, thread_id, user_id 
		 FROM posts WHERE thread_id = $1`

	var posts []entity.Post
	rows, err := f.db.Query(query, threadID)
	if err != nil {
		f.logger.Error("Ошибка при запросе постов треда",
			zap.Int("threadID", threadID),
			zap.Error(err))
		return nil, fmt.Errorf("Ошибка поиска постов по id треда: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post entity.Post
		if err := rows.Scan(
			&post.ID,
			&post.Content,
			&post.CreateAt,
			&post.ThreadID,
			&post.UserID,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				f.logger.Warn("Посты для треда не найдены",
					zap.Int("threadID", threadID))
				return nil, entity.ErrorNotFoundPost
			}
			f.logger.Error("Ошибка при сканировании поста",
				zap.Int("threadID", threadID),
				zap.Error(err))
			return nil, fmt.Errorf("Ошибка поиска поста по id треда: %w", err)
		}
		posts = append(posts, post)
	}

	f.logger.Debug("Посты успешно получены",
		zap.Int("threadID", threadID),
		zap.Int("count", len(posts)))
	return posts, nil
}

func (f *forumRepository) GetPostsByUserID(id int) ([]entity.Post, error) {
	f.logger.Debug("Получение постов по ID пользователя", zap.Int("userID", id))
	query := `
        SELECT id, content, create_at, thread_id, user_id
        FROM posts
        WHERE user_id = $1
        ORDER BY create_at DESC`

	rows, err := f.db.Query(query, id)
	if err != nil {
		f.logger.Error("Ошибка выполнения запроса постов пользователя",
			zap.Int("userID", id),
			zap.Error(err))
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var posts []entity.Post

	for rows.Next() {
		var post entity.Post
		if err := rows.Scan(
			&post.ID,
			&post.Content,
			&post.CreateAt,
			&post.ThreadID,
			&post.UserID,
		); err != nil {
			f.logger.Error("Ошибка сканирования поста",
				zap.Int("userID", id),
				zap.Error(err))
			return nil, fmt.Errorf("ошибка сканирования поста: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		f.logger.Error("Ошибка при обработке результатов",
			zap.Int("userID", id),
			zap.Error(err))
		return nil, fmt.Errorf("ошибка при итерации по результатам: %w", err)
	}

	if len(posts) == 0 {
		f.logger.Warn("Посты пользователя не найдены",
			zap.Int("userID", id))
		return nil, entity.ErrorNotFoundUser
	}

	f.logger.Debug("Посты пользователя успешно получены",
		zap.Int("userID", id),
		zap.Int("count", len(posts)))
	return posts, nil
}

func (f *forumRepository) GetPostByID(id int) (entity.Post, error) {
	query := `SELECT * FROM posts WHERE id = $1`

	var post entity.Post
	err := f.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.Content,
		&post.CreateAt,
		&post.ThreadID,
		&post.UserID,
	)

	if err != nil {
		return entity.Post{}, err
	}

	return post, nil
}

func (f *forumRepository) DeletePostByID(id int) error {
	f.logger.Debug("Удаление поста по ID", zap.Int("id", id))
	query := `DELETE FROM posts WHERE id = $1`
	result, err := f.db.Exec(query, id)
	if err != nil {
		f.logger.Error("Ошибка при удалении поста",
			zap.Int("id", id),
			zap.Error(err))
		return fmt.Errorf("Ошибка при удалении поста: %w", err)
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		f.logger.Error("Ошибка при получении количества удаленных строк",
			zap.Int("id", id),
			zap.Error(err))
		return fmt.Errorf("Ошибка получения измененных строк: %w", err)
	}
	if rowAffected == 0 {
		f.logger.Warn("Пост не найден для удаления",
			zap.Int("id", id))
		return entity.ErrorNotFoundPost
	}

	f.logger.Info("Пост успешно удален",
		zap.Int("id", id),
		zap.Int64("rowsAffected", rowAffected))
	return nil
}

func (f *forumRepository) GetThreadsByUserID(userId int) ([]entity.Thread, error) {
	f.logger.Debug("Получение тредов по ID пользователя", zap.Int("userID", userId))
	query := `SELECT * FROM threads 
         	  WHERE user_ID = $1
         	  ORDER BY create_at DESC`
	threads, err := f.db.Query(query, userId)
	if err != nil {
		f.logger.Error("Ошибка при запросе тредов пользователя",
			zap.Int("userID", userId),
			zap.Error(err))
		return nil, err
	}
	defer threads.Close()

	var searchThreads []entity.Thread
	for threads.Next() {
		var thread entity.Thread
		if err = threads.Scan(&thread.ID, &thread.Title, &thread.Content, &thread.CreateAt, &thread.UserID); err != nil {
			f.logger.Error("Ошибка при сканировании треда",
				zap.Int("userID", userId),
				zap.Error(err))
			return nil, err
		}
		searchThreads = append(searchThreads, thread)
	}
	f.logger.Debug("Треды пользователя успешно получены",
		zap.Int("userID", userId),
		zap.Int("count", len(searchThreads)))
	return searchThreads, nil
}

func (f *forumRepository) LinkPostToChat(chat entity.Chat) error {
	f.logger.Debug("Привязка поста к чату",
		zap.Int("threadID", chat.ThreadID),
		zap.Int("postID", chat.PostID),
		zap.Int("userID", chat.UserID))

	query := `INSERT INTO chat (thread_id, user_id, post_id) VALUES ($1, $2, $3)`
	_, err := f.db.Exec(query, chat.ThreadID, chat.UserID, chat.PostID)
	if err != nil {
		f.logger.Error("Ошибка при привязке поста к чату",
			zap.Any("chat", chat),
			zap.Error(err))
		return err
	}

	f.logger.Info("Пост успешно привязан к чату",
		zap.Int("threadID", chat.ThreadID),
		zap.Int("postID", chat.PostID))
	return nil
}

func (f *forumRepository) GetChatPosts(threadID int) ([]entity.Post, error) {
	f.logger.Debug("Получение постов чата по ID треда", zap.Int("threadID", threadID))
	query := `
		SELECT p.id, p.content, p.create_at, p.thread_id, p.user_id
		FROM posts p
		JOIN chat c ON p.id = c.post_id
		WHERE c.thread_id = $1
		ORDER BY p.create_at ASC`

	rows, err := f.db.Query(query, threadID)
	if err != nil {
		f.logger.Error("Ошибка при запросе постов чата",
			zap.Int("threadID", threadID),
			zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var post entity.Post
		if err := rows.Scan(&post.ID, &post.Content, &post.CreateAt, &post.ThreadID, &post.UserID); err != nil {
			f.logger.Error("Ошибка при сканировании поста чата",
				zap.Int("threadID", threadID),
				zap.Error(err))
			return nil, err
		}
		posts = append(posts, post)
	}

	f.logger.Debug("Посты чата успешно получены",
		zap.Int("threadID", threadID),
		zap.Int("count", len(posts)))
	return posts, nil
}

func (f *forumRepository) CheckUserByID(any entity.User, id int) (bool, error) {
	query := `SELECT id, name, email, role FROM users WHERE id = $1`

	var user struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	err := f.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		return false, err
	}

	if user.Role == "admin" {
		return true, nil
	}

	if user.ID != any.USER_ID() {
		return false, fmt.Errorf("Нет прав: %d %d %s", user.ID, any.USER_ID(), user.Role)
	}
	return true, nil
}
