cache [![Build Status](https://travis-ci.org/issue9/cache.svg?branch=master)](https://travis-ci.org/issue9/cache)
======

一组缓存接口，目前只实现了基于内存的缓存，后续会再添加 redis。

** memory.Set 可以存储任意类型，memcache.Set 只能存储字符串类型，memcache.Incr 不能为负数。
相同的接口，却达不到相同的功能，不如不做！**

### 安装

```shell
go get github.com/issue9/cache
```


### 文档

[![Go Walker](https://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/issue9/cache)
[![GoDoc](https://godoc.org/github.com/issue9/cache?status.svg)](https://godoc.org/github.com/issue9/cache)


### 版权

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
