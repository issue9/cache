// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package locales 提供了本地化的内容
package locales

import "embed"

//go:embed *.yaml
var Locales embed.FS
