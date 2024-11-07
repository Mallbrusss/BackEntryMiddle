package service

import (
	"github.com/Mallbrusss/BackEntryMiddle/models"
	"github.com/stretchr/testify/mock"
)

type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) CreateDocument(document *models.Document, grant []string) error {
	args := m.Called(document, grant)
	return args.Error(0)
}

func (m *MockDocumentRepository) GetDocuments(login string, filter map[string]any, limit int) ([]models.Document, error) {
	args := m.Called(login, filter, limit)

	if doc, ok := args.Get(0).([]models.Document); ok {
		return doc, args.Error(1)
	}
	return nil, args.Error(1)

}

func (m *MockDocumentRepository) GetDocumentByID(documentID, login string) (*models.Document, error) {
	args := m.Called(documentID, login)

	if doc, ok := args.Get(0).(*models.Document); ok {
		return doc, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockDocumentRepository) DeleteDocument(document *models.Document) error {
	args := m.Called(document)
	return args.Error(0)
}

func (m *MockDocumentRepository) GetDocumentAccessByID(documentID string) ([]models.DocumentAccess, error) {
	args := m.Called(documentID)

	if docA, ok := args.Get(0).([]models.DocumentAccess); ok {
		return docA, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDocumentRepository) FindByToken(token string) (*models.User, error) {
	args := m.Called(token)

	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDocumentRepository) FindByLogin(login string) (*models.User, error) {
	args := m.Called(login)

	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDocumentRepository) IsPermission(documentId string, user *models.User) (bool, error) {
	args := m.Called(documentId, user)

	return args.Bool(0), args.Error(1)
}
