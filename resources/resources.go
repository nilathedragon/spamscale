package resources

import (
	"embed"
	"io/fs"
)

//go:embed vfs/**
var vfs embed.FS

func Read(path string) ([]byte, error) {
	return vfs.ReadFile("vfs/" + path)
}

func ReadDir(path string) ([]fs.DirEntry, error) {
	return vfs.ReadDir("vfs/" + path)
}
