package mocks

//
//import (
//	"github.com/DrusGalkin/forum-client/internal/entity"
//	"github.com/stretchr/testify/mock"
//)
//
//type ForumUseCase struct {
//	mock.Mock
//}
//
//func (m *ForumUseCase) GetAllThreads() ([]entity.Thread, error) {
//	args := m.Called()
//	return args.Get(0).([]entity.Thread), args.Error(1)
//}
//
//func (m *ForumUseCase) GetThreadByID(id int) (entity.Thread, error) {
//	args := m.Called(id)
//	return args.Get(0).(entity.Thread), args.Error(1)
//}
//
//func (m *ForumUseCase) CreateThread(thread entity.Thread) (entity.Thread, error) {
//	args := m.Called(thread)
//	return args.Get(0).(entity.Thread), args.Error(1)
//}
//
//func (m *ForumUseCase) DeleteThreadByID(id int, userID int) error {
//	args := m.Called(id, userID)
//	return args.Error(0)
//}
//
//func (m *ForumUseCase) CreatePost(post entity.Post) (entity.Post, error) {
//	args := m.Called(post)
//	return args.Get(0).(entity.Post), args.Error(1)
//}
//
//func (m *ForumUseCase) GetChatPosts(threadID int) ([]entity.Post, error) {
//	args := m.Called(threadID)
//	return args.Get(0).([]entity.Post), args.Error(1)
//}
//
//func (m *ForumUseCase) GetPostByThreadID(threadID int) ([]entity.Post, error) {
//	args := m.Called(threadID)
//	return args.Get(0).([]entity.Post), args.Error(1)
//}
//
//func (m *ForumUseCase) DeletePostByID(id int, userID int) error {
//	args := m.Called(id, userID)
//	return args.Error(0)
//}
//
//func (m *ForumUseCase) GetPostsByUserID(id int) ([]entity.Post, error) {
//	args := m.Called(id)
//	return args.Get(0).([]entity.Post), args.Error(1)
//}
//
//func (m *ForumUseCase) CheckUserByID(any entity.User, id int) (bool, error) {
//	args := m.Called(any, id)
//	return args.Bool(0), args.Error(1)
//}
//
//func (m *ForumUseCase) GetUserThreads(userId int) ([]entity.Thread, error) {
//	args := m.Called(userId)
//	return args.Get(0).([]entity.Thread), args.Error(1)
//}
//
//func (m *ForumUseCase) EditThread(thread entity.Thread, userID int) error {
//	args := m.Called(thread, userID)
//	return args.Error(0)
//}
