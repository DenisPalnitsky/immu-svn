package pkg

import (
	"testing"

	"github.com/DenisPalnitsky/immu-svn/mocks"
	"github.com/DenisPalnitsky/immu-svn/pkg/data"
	"github.com/stretchr/testify/assert"
)

func TestSvn_New(t *testing.T) {
	storage := &mocks.Storage{}
	t.Run("CreateSvnInstance", func(t *testing.T) {
		storage.EXPECT().CreateRepo("test").Return(nil)
		svn, err := NewSnv(storage, "./testdata/repo")
		assert.NoError(t, err)
		assert.Equal(t, "repo", svn.RepoName)
	})

	t.Run("CreateSvnInstanceWithInvalidPath", func(t *testing.T) {
		_, err := NewSnv(storage, "")
		assert.Error(t, err)
	})
}

func TestSvn_Init(t *testing.T) {
	t.Run("InitRepoSuccessfully ", func(t *testing.T) {
		storage := &mocks.Storage{}
		storage.EXPECT().CreateRepo("repo").Return(nil)
		svn, err := NewSnv(storage, "./testdata/repo")
		assert.NoError(t, err)
		assert.Equal(t, "repo", svn.RepoName)
		err = svn.Init()
		assert.NoError(t, err)
	})

	t.Run("RepoAlreadyInitialized", func(t *testing.T) {
		storage := &mocks.Storage{}
		storage.EXPECT().CreateRepo("repo").Return(data.CollectionAlreadyExist)
		svn, err := NewSnv(storage, "./testdata/repo")
		assert.NoError(t, err)
		err = svn.Init()
		assert.ErrorContains(t, err, "repository already initialized")
	})
}

func TestSvn_Commit(t *testing.T) {

	t.Run("CommitFilesRunsOk", func(t *testing.T) {
		storage := &mocks.Storage{}
		storage.EXPECT().GetCollectionInfo("repo").Return(nil)
		storage.EXPECT().AddOrUpdateFiles("repo", map[string]string{
			"readme.txt": "hello go",
			"test.txt":   "Hello Universe\n",
		}).Return(2, 0, nil)
		svn, err := NewSnv(storage, "./testdata/repo")
		assert.NoError(t, err)
		assert.Equal(t, "repo", svn.RepoName)
		commitInfo, err := svn.Commit()
		assert.NoError(t, err)
		assert.Equal(t, 2, commitInfo.FilesAdded)
		assert.Equal(t, 0, commitInfo.FilesUpdated)
	})

	t.Run("CommitFileThatIsTooBig", func(t *testing.T) {
		storage := &mocks.Storage{}
		storage.EXPECT().GetCollectionInfo("bad-repo").Return(nil)

		svn, err := NewSnv(storage, "./testdata/bad-repo")
		assert.NoError(t, err)
		_, err = svn.Commit()
		assert.Error(t, err)
	})

	t.Run("RepositoryNotFound", func(t *testing.T) {
		storage := &mocks.Storage{}
		storage.EXPECT().GetCollectionInfo("repo").Return(data.RepositoryNotFound)
		svn, err := NewSnv(storage, "./testdata/repo")
		assert.NoError(t, err)
		_, err = svn.Commit()
		assert.Error(t, err)
	})
}
