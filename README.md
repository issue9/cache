# cache

[![codecov](https://codecov.io/gh/issue9/cache/branch/master/graph/badge.svg)](https://codecov.io/gh/issue9/cache)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/issue9/cache)](https://pkg.go.dev/github.com/issue9/cache)
![Go version](https://img.shields.io/github/go-mod/go-version/issue9/cache)
![License](https://img.shields.io/github/license/issue9/cache)

通用的缓存接口

目前支持以下组件：

名称       | 包                                   | 状态
-----------|--------------------------------------|-----
memory     | 内存                                 | [![memory](https://github.com/issue9/cache/workflows/memory/badge.svg)](https://github.com/issue9/cache/actions?query=workflow%3Amemory)
memcached  | github.com/bradfitz/gomemcache       | [![memcache](https://github.com/issue9/cache/workflows/memcached/badge.svg)](https://github.com/issue9/cache/actions?query=workflow%3Amemcached)
redis      | github.com/redis/go-redis            | [![memcache](https://github.com/issue9/cache/workflows/redis/badge.svg)](https://github.com/issue9/cache/actions?query=workflow%3Aredis)

```go
// memory
c, _ := memory.New(...)
c.Set("number", 1)
var v int
c.Get("number",&v)
print(v)

// memcached
c = memcache.New("localhost:11211")
c.Set("number", 1)
c.Get("number", &v)
print(v)
```

## 安装

```shell
go get github.com/issue9/cache
```

## 版权

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
