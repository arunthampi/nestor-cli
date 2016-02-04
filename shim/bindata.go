// Code generated by go-bindata.
// sources:
// bindata.go
// index.js
// DO NOT EDIT!

package shim

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _bindataGo = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func bindataGoBytes() ([]byte, error) {
	return bindataRead(
		_bindataGo,
		"bindata.go",
	)
}

func bindataGo() (*asset, error) {
	bytes, err := bindataGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "bindata.go", size: 0, mode: os.FileMode(420), modTime: time.Unix(1454610410, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _indexJs = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x6c\x52\x3d\x8f\xd4\x30\x10\xed\xf3\x2b\xdc\xc5\x91\x22\x6f\x7f\x2b\x0a\x0a\xca\x43\xe8\x80\x0a\x21\xcb\x6b\xcf\x2e\x11\x89\x1d\x6c\x27\x04\xad\xf6\xbf\x33\x63\xe7\x53\xb7\x29\xac\xcc\xbc\xf7\xe6\x7b\x54\x9e\x49\x0f\x57\xf6\x81\x79\xf8\x33\x34\x1e\x78\x29\x4e\x41\xfb\xa6\x8f\x27\xeb\x0c\xc8\xce\x99\xa1\x85\x70\x52\xe1\xb7\x85\x10\x9d\x2f\xab\x73\x41\xb2\x37\x77\x71\x11\x75\x24\x17\xc9\xc8\xfe\x6f\x30\xc5\x57\x08\x41\xdd\x60\x41\x77\xae\xcc\xf9\x1e\xc0\x2f\x20\xfd\x67\xef\xe7\x14\xff\xa3\x51\x7d\xdc\xe0\x83\xf3\x5c\x14\x30\xf5\xce\xc7\x20\x7e\x29\x6b\x5a\xca\x70\x1d\xac\x8e\x8d\xb3\x1c\x46\xb0\xb1\x66\x3a\x4e\x15\xbb\x17\x8c\x5d\x9d\xe7\x14\x16\xec\xf8\xc5\xbb\x9e\x35\x96\x25\x8a\x90\x32\x77\x22\x11\xc9\x54\xc6\x7a\xef\x34\x56\x28\xd0\xf5\x63\x16\xfc\xc4\xe0\xef\x05\x2b\x7a\x46\xdd\xa3\xc0\xc7\x43\xab\xa6\x4f\x44\xdc\x09\x92\x53\x26\x8b\x88\x43\x6e\xd8\xc2\xdf\xd4\x3b\xdf\x34\x82\x20\x39\x34\xa6\x4e\x75\x1c\xbf\xfb\x13\x1f\x26\x74\xae\x7b\xd9\xa5\x15\x1a\x87\x61\xa1\xa5\x28\x4f\x04\x8f\x8a\x2a\xe8\xd6\x95\x50\x11\xbb\x8d\x70\x2a\xa0\xde\x87\x8b\x08\x26\x8d\x9f\x57\x4c\x8a\xb4\x61\x7e\x60\xa9\x2e\xd5\xbd\x97\xe6\xb6\x91\x99\x91\x32\xcf\x0d\xed\xb2\x5e\x47\x63\xe0\x32\xdc\x5e\xf1\xb6\xb6\x1c\x42\xad\x4b\xa7\x5c\x87\x9d\xf3\xc4\x40\x6e\x1a\xf5\x7c\xa3\xf9\x42\xcb\xea\x00\xa6\x48\x1e\x34\x34\x23\xf0\xb9\xdf\x7a\xbb\x0f\xe3\x2c\x2c\xeb\xc6\x23\x11\x61\xd0\x1a\xc0\xf0\x65\xc8\xd1\xc9\x00\xd6\xbc\xcc\x81\xa2\xfb\x8a\x56\xbd\x81\x1e\xfa\xf6\xdf\x86\xbe\x91\x59\x6c\x03\xa6\xf7\x51\xfc\x0f\x00\x00\xff\xff\x58\x9f\xa1\x8e\x51\x03\x00\x00")

func indexJsBytes() ([]byte, error) {
	return bindataRead(
		_indexJs,
		"index.js",
	)
}

func indexJs() (*asset, error) {
	bytes, err := indexJsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "index.js", size: 849, mode: os.FileMode(420), modTime: time.Unix(1454610401, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"bindata.go": bindataGo,
	"index.js":   indexJs,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"bindata.go": &bintree{bindataGo, map[string]*bintree{}},
	"index.js":   &bintree{indexJs, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
