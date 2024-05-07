package document

import "gorm.io/gorm"

type documentRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) DocumentRepository {
	return &documentRepository{db}
}

func (r *documentRepository) GetDocumentByID(id uint) (*Document, error) {
	var doc Document
	if err := r.db.First(&doc, id).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *documentRepository) GetAllDocuments() ([]Document, error) {
	var documents []Document
	if err := r.db.Find(&documents).Error; err != nil {
		return nil, err
	}
	return documents, nil
}

func (r *documentRepository) CreateDocument(doc *Document) error {
	if err := r.db.Create(doc).Error; err != nil {
		return err
	}
	return nil
}

func (r *documentRepository) UpdateDocument(doc *Document) error {
	if err := r.db.Save(doc).Error; err != nil {
		return err
	}
	return nil
}

func (r *documentRepository) DeleteDocument(id uint) error {
	if err := r.db.Delete(&Document{}, id).Error; err != nil {
		return err
	}
	return nil
}
