package mapper

import (
	"github.com/ProtobufMan/bufman/internal/dal"
	"github.com/ProtobufMan/bufman/internal/model"
)

type FileMapper interface {
	FindAllBlobsByCommitID(commitID string) (model.FileBlobs, error)
	FindManifestByCommitID(commitID string) (*model.FileManifest, error)
	FindBlobByCommitIDAndPath(commitID, path string) (*model.FileBlob, error)
}

type FileMapperImpl struct{}

func (f *FileMapperImpl) FindAllBlobsByCommitID(commitID string) (model.FileBlobs, error) {
	return dal.FileBlob.Where(dal.FileBlob.CommitID.Eq(commitID)).Find()
}

func (f *FileMapperImpl) FindManifestByCommitID(commitID string) (*model.FileManifest, error) {
	return dal.FileManifest.Where(dal.FileManifest.CommitID.Eq(commitID)).First()
}

func (f *FileMapperImpl) FindBlobByCommitIDAndPath(commitID, path string) (*model.FileBlob, error) {
	return dal.FileBlob.Where(dal.FileBlob.CommitID.Eq(commitID), dal.FileBlob.FileName.Eq(path)).First()
}
