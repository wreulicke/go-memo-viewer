package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	gas "github.com/jessevdk/go-assets"
	"github.com/mattes/migrate"
	"github.com/wreulicke/go-memo-viewer/assets"
	"github.com/wreulicke/go-memo-viewer/driver"
	"github.com/wreulicke/go-memo-viewer/memo"
	"github.com/wreulicke/go-memo-viewer/migration_assets"

	_ "github.com/mattes/migrate/database/mysql"
)

type FS struct {
	gas.FileSystem
}

func (*FS) Exists(prefix string, path string) bool {
	return path == "/" || path == "/index.html"
}

//go:generate go-assets-builder -s="/public" -p assets -o assets/bindata.go public
//go:generate go-assets-builder -s="/migrations" -p migration_assets -o migration_assets/bindata.go migrations
func main() {
	d, err := driver.WithInstance(migration_assets.Assets)

	if err != nil {
		fmt.Println("Migration file system broken.")
		return
	}

	m, err := migrate.NewWithSourceInstance("go-assets", d, "mysql://root@tcp(127.0.0.1:3306)/test")

	if err != nil {
		fmt.Println(err)
		return
	}

	err = m.Migrate(d.Version)

	if err != nil && err.Error() != "no change" {
		fmt.Println("migration is failed: ", err)
		return
	}

	router := gin.Default()
	// store := sessions.NewCookieStore([]byte("secret"))
	// router.Use(sessions.Sessions("go-memo-viewer", store))

	router.NoRoute(static.Serve("/", &FS{FileSystem: *assets.Assets}))
	conn, err := sql.Open("mysql", "root@tcp(localhost:3306)/test")
	memo.Route(router.Group("/memo"), conn)

	router.Run(":8080")
}
