package pkg

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/DenisPalnitsky/immu-svn/pkg/data"
)

type Storage interface {
	CreateRepo(repoName string) error
	AddOrUpdateFiles(repoName string, files map[string]string) (added int, updated int, err error)
	Diff(repoName string, filePath string) ([]data.DiffLogItem, error)
	GetCollectionInfo(repoName string) error
}

type CommitInfo struct {
	FilesUpdated int
	FilesAdded   int
}

type Svn struct {
	storage    Storage
	RepoName   string
	workingDir string
}

func NewSnv(storage Storage, workingDir string) (*Svn, error) {
	repoName, err := getRepoName(workingDir)
	if err != nil {
		return nil, err
	}
	return &Svn{
		storage:    storage,
		RepoName:   repoName,
		workingDir: workingDir,
	}, nil
}

func (s *Svn) Init() error {
	err := s.storage.CreateRepo(s.RepoName)
	if err != nil {
		if err == data.CollectionAlreadyExist {
			return errors.New("repository already initialized")
		}
		return err
	}
	return nil
}

func (s *Svn) Commit() (*CommitInfo, error) {
	if err := s.storage.GetCollectionInfo(s.RepoName); err != nil {
		if err == data.RepositoryNotFound {
			return nil, fmt.Errorf("repository %s not found", s.RepoName)
		} else {
			return nil, fmt.Errorf("error getting collection info %w", err)
		}
	}

	files, err := listFiles(s.workingDir)
	if err != nil {
		return nil, fmt.Errorf("error listing files %w", err)
	}

	added, updated, err := s.storage.AddOrUpdateFiles(s.RepoName, files)
	if err != nil {
		return nil, fmt.Errorf("error adding or updating files %w", err)
	}

	return &CommitInfo{
		FilesAdded:   added,
		FilesUpdated: updated,
	}, nil
}

func (s *Svn) Diff(filePath string) ([]data.DiffLogItem, error) {
	return s.storage.Diff(s.RepoName, filePath)
}

func getRepoName(dir string) (string, error) {
	repoName := dir[strings.LastIndex(dir, "/")+1:]
	if len(repoName) == 0 {
		return "", errors.New("invalid repository name")
	}
	return repoName, nil
}

func listFiles(directoryPath string) (map[string]string, error) {
	files := make(map[string]string)

	err := filepath.Walk(directoryPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip .xxx directories
		if info.IsDir() && len(info.Name()) > 1 && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			if len(content) >= 512 {
				return fmt.Errorf("file %s is too big. We support only files up to 512 bytes", path)
			}

			relativePath, err := filepath.Rel(directoryPath, path)
			if err != nil {
				return err
			}
			files[relativePath] = string(content)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
