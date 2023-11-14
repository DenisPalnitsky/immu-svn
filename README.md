# ImmuDb based source management system

Immu-svn is a source code management system that stores code in ImmuDb which makes you code history immutable in compliant with every possible regulation.

# How it works
immu-svn is a cli application

- Build the cli
- Init `IMMUDB_API_KEY` env var with your ImmuDB Vault API key 
- Go inside the directory with your source code
- Run `immu-svn init` to initialize the repository 
- Run `immu-svn commit` to add or update all the files to the repository
- Run `immu-svn diff -f [filename]` to see the history of file changes

## Here is an example
```
export IMMUDB_API_KEY=<API KEY>
git clone https://github.com/DenisPalnitsky/immu-svn.git
cd immu-svn
go run main.go init -d pkg/testdata/repo
go run main.go commit -d pkg/testdata/repo
echo "Hello Universe">pkg/testdata/repo/test.txt
go run main.go commit -d pkg/testdata/repo
go run main.go diff -d pkg/testdata/repo -f test.txt
```

# Limitations
- Currently, only files with less than 512 caracteres are supported
- Not more than 100 files in repository


