package driver

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	gas "github.com/jessevdk/go-assets"
	"github.com/mattes/migrate/source"
)

func init() {
	source.Register("go-assets", &Bindata{})
}

type Bindata struct {
	path       string
	fs         *gas.FileSystem
	migrations *source.Migrations
	Version    uint
}

func (b *Bindata) Open(url string) (source.Driver, error) {
	return nil, fmt.Errorf("not yet implemented")
}

var (
	ErrNoAssetSource = fmt.Errorf("expects *AssetSource")
)

func WithInstance(fileSystem *gas.FileSystem) (*Bindata, error) {
	bn := &Bindata{
		path:       "<go-assets>",
		fs:         fileSystem,
		migrations: source.NewMigrations(),
	}

	var i uint
	for _, f := range fileSystem.Files {
		m, err := source.DefaultParse(f.Name())
		if err != nil {
			continue // ignore files that we can't parse
		}

		if !bn.migrations.Append(m) {
			return nil, fmt.Errorf("unable to parse file %v", f)
		}
		i = i + 1
	}

	bn.Version = i

	return bn, nil
}

func (b *Bindata) Close() error {
	return nil
}

func (b *Bindata) First() (version uint, err error) {
	if v, ok := b.migrations.First(); !ok {
		return 0, &os.PathError{"first", b.path, os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (b *Bindata) Prev(version uint) (prevVersion uint, err error) {
	if v, ok := b.migrations.Prev(version); !ok {
		return 0, &os.PathError{fmt.Sprintf("prev for version %v", version), b.path, os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (b *Bindata) Next(version uint) (nextVersion uint, err error) {
	if v, ok := b.migrations.Next(version); !ok {
		return 0, &os.PathError{fmt.Sprintf("next for version %v", version), b.path, os.ErrNotExist}
	} else {
		return v, nil
	}
}

func (b *Bindata) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := b.migrations.Up(version); ok {
		file, err := b.fs.Open("/" + m.Raw)
		if err != nil {
			return nil, "", err
		}

		bs, err := ioutil.ReadAll(file)

		if err != nil {
			return nil, "", err
		}

		return ioutil.NopCloser(bytes.NewReader(bs)), m.Identifier, nil
	}
	return nil, "", &os.PathError{fmt.Sprintf("read version %v", version), b.path, os.ErrNotExist}
}

func (b *Bindata) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := b.migrations.Down(version); ok {

		file, err := b.fs.Open("/" + m.Raw)
		if err != nil {
			return nil, "", err
		}
		bs, err := ioutil.ReadAll(file)

		if err != nil {
			return nil, "", err
		}

		return ioutil.NopCloser(bytes.NewReader(bs)), m.Identifier, nil
	}
	return nil, "", &os.PathError{fmt.Sprintf("read version %v", version), b.path, os.ErrNotExist}
}
