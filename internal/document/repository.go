package document

import (
	"gorm.io/gorm"
)

// documentRepository implements DocumentRepository interface
type documentRepository struct {
	db *gorm.DB
}

// NewRepository creates a new instance of DocumentRepository
func NewRepository(db *gorm.DB) DocumentRepository {
	return &documentRepository{db}
}

// GetDocumentByID retrieves a document by ID
func (r *documentRepository) GetDocumentByID(id uint) (*Document, error) {
	var doc Document
	if err := r.db.First(&doc, id).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

// GetAllDocuments retrieves all documents
func (r *documentRepository) GetAllDocuments() ([]Document, error) {
	var documents []Document
	if err := r.db.Find(&documents).Error; err != nil {
		return nil, err
	}
	return documents, nil
}

// CreateDocument creates a new document
func (r *documentRepository) CreateDocument(doc *Document) error {
	if err := r.db.Create(doc).Error; err != nil {
		return err
	}
	return nil
}

// UpdateDocument updates an existing document
func (r *documentRepository) UpdateDocument(doc *Document) error {
	if err := r.db.Save(doc).Error; err != nil {
		return err
	}
	return nil
}

// DeleteDocument deletes a document by ID
func (r *documentRepository) DeleteDocument(id uint) error {
	if err := r.db.Delete(&Document{}, id).Error; err != nil {
		return err
	}
	return nil
}
