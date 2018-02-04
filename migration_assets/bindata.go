package migration_assets

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assetsa44f6bb8b9a99d00200cea3235b735b6532ab1c1 = "create table memo(id varchar(128) primary key, title text, text longtext);\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"1_initial_table.up.sql"}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1517737033, 1517737033766598484),
		Data:     nil,
	}, "/1_initial_table.up.sql": &assets.File{
		Path:     "/1_initial_table.up.sql",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1517743629, 1517743629881845000),
		Data:     []byte(_Assetsa44f6bb8b9a99d00200cea3235b735b6532ab1c1),
	}}, "")
