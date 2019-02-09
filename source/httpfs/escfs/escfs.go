// Code generated by "esc -o escfs/escfs.go -pkg escfs -prefix escfs/ escfs/"; DO NOT EDIT.

package escfs

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	if !f.isDir {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is not directory", f.name)
	}

	fis, ok := _escDirs[f.local]
	if !ok {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is directory, but we have no info about content of this dir, local=%s", f.name, f.local)
	}
	limit := count
	if count <= 0 || limit > len(fis) {
		limit = len(fis)
	}

	if len(fis) == 0 && count > 0 {
		return nil, io.EOF
	}

	return []os.FileInfo(fis[0:limit]), nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/Test/1_foobar.down.sql": {
		name:    "1_foobar.down.sql",
		local:   "escfs/Test/1_foobar.down.sql",
		size:    7,
		modtime: 1549747705,
		compressed: `
H4sIAAAAAAAC/zJUSMkvz+MCBAAA//+UkzbRBwAAAA==
`,
	},

	"/Test/1_foobar.up.sql": {
		name:    "1_foobar.up.sql",
		local:   "escfs/Test/1_foobar.up.sql",
		size:    5,
		modtime: 1549747662,
		compressed: `
H4sIAAAAAAAC/zJUKC3gAgQAAP//RvOBZwUAAAA=
`,
	},

	"/Test/3_foobar.up.sql": {
		name:    "3_foobar.up.sql",
		local:   "escfs/Test/3_foobar.up.sql",
		size:    5,
		modtime: 1549747673,
		compressed: `
H4sIAAAAAAAC/zJWKC3gAgQAAP//JqBBHQUAAAA=
`,
	},

	"/Test/4_foobar.down.sql": {
		name:    "4_foobar.down.sql",
		local:   "escfs/Test/4_foobar.down.sql",
		size:    7,
		modtime: 1549747714,
		compressed: `
H4sIAAAAAAAC/zJRSMkvz+MCBAAA//8zvA6DBwAAAA==
`,
	},

	"/Test/4_foobar.up.sql": {
		name:    "4_foobar.up.sql",
		local:   "escfs/Test/4_foobar.up.sql",
		size:    5,
		modtime: 1549747679,
		compressed: `
H4sIAAAAAAAC/zJRKC3gAgQAAP//NnxhrwUAAAA=
`,
	},

	"/Test/5_foobar.up.sql": {
		name:    "5_foobar.up.sql",
		local:   "escfs/Test/5_foobar.up.sql",
		size:    5,
		modtime: 1549747685,
		compressed: `
H4sIAAAAAAAC/zJVKC3gAgQAAP//hlUBkgUAAAA=
`,
	},

	"/Test/7_foobar.down.sql": {
		name:    "7_foobar.down.sql",
		local:   "escfs/Test/7_foobar.down.sql",
		size:    7,
		modtime: 1549747727,
		compressed: `
H4sIAAAAAAAC/zJXSMkvz+MCBAAA//+upuayBwAAAA==
`,
	},

	"/Test/7_foobar.up.sql": {
		name:    "7_foobar.up.sql",
		local:   "escfs/Test/7_foobar.up.sql",
		size:    5,
		modtime: 1549747691,
		compressed: `
H4sIAAAAAAAC/zJXKC3gAgQAAP//5gbB6AUAAAA=
`,
	},

	"/escfs.go": {
		name:    "escfs.go",
		local:   "escfs/escfs.go",
		size:    0,
		modtime: 1549748432,
		compressed: `
H4sIAAAAAAAC/wEAAP//AAAAAAAAAAA=
`,
	},

	"/": {
		name:  "/",
		local: `escfs/`,
		isDir: true,
	},

	"/Test": {
		name:  "Test",
		local: `escfs/Test`,
		isDir: true,
	},
}

var _escDirs = map[string][]os.FileInfo{

	"escfs/": {
		_escData["/Test"],
		_escData["/escfs.go"],
	},

	"escfs/Test": {
		_escData["/Test/1_foobar.down.sql"],
		_escData["/Test/1_foobar.up.sql"],
		_escData["/Test/3_foobar.up.sql"],
		_escData["/Test/4_foobar.down.sql"],
		_escData["/Test/4_foobar.up.sql"],
		_escData["/Test/5_foobar.up.sql"],
		_escData["/Test/7_foobar.down.sql"],
		_escData["/Test/7_foobar.up.sql"],
	},
}
