// SPDX-License-Identifier: MIT

// Package file 文件缓存的实现
package file

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/issue9/cache"
)

const modePerm = os.ModePerm

type file struct {
	root   string
	errlog *log.Logger

	ticker *time.Ticker
	done   chan struct{}

	items *sync.Map
}

type item struct {
	path   string
	dur    time.Duration
	expire time.Time // 过期的时间
}

func (i *item) update() {
	i.expire = time.Now().Add(i.dur)
}

func (i *item) isExpired(now time.Time) bool {
	return i.expire.Before(now)
}

// New 返回基于文件系统的缓存
//
// root 文件系统的根目录；
// gc 表示执行回收操作的间隔。
func New(root string, gc time.Duration, errlog *log.Logger) cache.Cache {
	f := &file{
		root:   root,
		errlog: errlog,
		ticker: time.NewTicker(gc),
		done:   make(chan struct{}, 1),
		items:  &sync.Map{},
	}

	go func(f *file) {
		for {
			select {
			case <-f.ticker.C:
				f.gc()
			case <-f.done:
				return
			}
		}
	}(f)

	return f
}

func (f *file) Get(key string) (val interface{}, err error) {
	bs, err := ioutil.ReadFile(f.getPath(key))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, cache.ErrCacheMiss
		}

		return nil, err
	}

	if err := cache.GoDecode(bs, &val); err != nil {
		return nil, err
	}

	return val, nil
}

func (f *file) Set(key string, val interface{}, seconds int) error {
	key = f.getPath(key)

	bs, err := cache.GoEncode(&val)
	if err != nil {
		return err
	}

	if seconds != cache.Forever {
		if i, found := f.items.Load(key); found {
			i.(*item).update()
		} else {
			dur := time.Second * time.Duration(seconds)
			f.items.Store(key, &item{
				path:   key,
				dur:    dur,
				expire: time.Now().Add(dur),
			})
		}
	}

	return ioutil.WriteFile(key, bs, modePerm)
}

func (f *file) Delete(key string) error {
	key = f.getPath(key)
	f.items.Delete(key)
	return os.Remove(key)
}

func (f *file) Exists(key string) bool {
	_, err := os.Stat(f.getPath(key))
	return err == nil || !errors.Is(err, os.ErrNotExist)
}

func (f *file) Clear() error {
	if err := os.RemoveAll(f.root); err != nil {
		return err
	}
	return os.Mkdir(f.root, modePerm)
}

func (f *file) Close() error {
	// NOTE: 关闭服务，不能清除服务器的内容

	f.ticker.Stop()
	close(f.done)
	return nil
}

func (f *file) gc() {
	now := time.Now()

	f.items.Range(func(key, val interface{}) bool {
		if v := val.(*item); v.isExpired(now) {
			f.items.Delete(key)
			if err := os.Remove(key.(string)); err != nil {
				f.errlog.Println(err)
			}
		}
		return true
	})
}

func (f *file) getPath(key string) string {
	return filepath.Join(f.root, md5String(key))
}

var m = md5.New()

func md5String(key string) string {
	m.Reset()
	io.WriteString(m, key)
	return hex.EncodeToString(m.Sum(nil))
}
