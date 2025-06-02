package mocks

import (
	"github.com/DrusGalkin/forum-client/internal/entity"
	"github.com/stretchr/testify/mock"
)

type ForumRepository struct {
	mock.Mock
}

func (m *ForumRepository) GetAllThreads() ([]entity.Thread, error) {
	args := m.Called()
	return args.Get(0).([]entity.Thread), args.Error(1)
}

func (m *ForumRepository) GetThreadByID(id int) (entity.Thread, error) {
	args := m.Called(id)
	return args.Get(0).(entity.Thread), args.Error(1)
}

func (m *ForumRepository) CreateThread(thread entity.Thread) (entity.Thread, error) {
	args := m.Called(thread)
	return args.Get(0).(entity.Thread), args.Error(1)
}

func (m *ForumRepository) DeleteThreadByID(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ForumRepository) GetThreadsByUserID(userId int) ([]entity.Thread, error) {
	args := m.Called(userId)
	return args.Get(0).([]entity.Thread), args.Error(1)
}

func (m *ForumRepository) CreatePost(post entity.Post) (entity.Post, error) {
	args := m.Called(post)
	return args.Get(0).(entity.Post), args.Error(1)
}

func (m *ForumRepository) GetPostsByThreadID(threadID int) ([]entity.Post, error) {
	args := m.Called(threadID)
	return args.Get(0).([]entity.Post), args.Error(1)
}

func (m *ForumRepository) DeletePostByID(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ForumRepository) GetPostsByUserID(id int) ([]entity.Post, error) {
	args := m.Called(id)
	return args.Get(0).([]entity.Post), args.Error(1)
}

func (m *ForumRepository) GetChatPosts(threadID int) ([]entity.Post, error) {
	args := m.Called(threadID)
	return args.Get(0).([]entity.Post), args.Error(1)
}

func (m *ForumRepository) LinkPostToChat(chat entity.Chat) error {
	args := m.Called(chat)
	return args.Error(0)
}

func (m *ForumRepository) CheckUserByID(user entity.User, id int) (bool, error) {
	args := m.Called(user, id)
	return args.Bool(0), args.Error(1)
}

func (m *ForumRepository) GetPostByID(id int) (entity.Post, error) {
	args := m.Called(id)
	return args.Get(0).(entity.Post), args.Error(1)
}

func (m *ForumRepository) EditThread(thread entity.Thread, userID int) error {
	args := m.Called(thread, userID)
	return args.Error(0)
}
