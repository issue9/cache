// SPDX-License-Identifier: MIT

package memcache

import "github.com/issue9/cache"

var _ cache.Cache = &Memcache{}
