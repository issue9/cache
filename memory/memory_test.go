// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package memory

import (
	"github.com/issue9/cache"
)

var _ cache.Cache = &Memory{}
