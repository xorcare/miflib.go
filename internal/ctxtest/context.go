// Copyright (c) 2020 Vasiliy Vasilyuk All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ctxtest

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var contextKey = new(struct{})

// Background the method is similar context.Background(),
// but it differs in that it adds special keys to check
// the context.
func Background() context.Context {
	return context.WithValue(context.Background(), contextKey, contextKey)
}

// Is checks that the current context is a test context.
func Is(ctx context.Context) bool {
	val, ok := ctx.Value(contextKey).(*struct{})
	return ok && val == contextKey
}

// Match this is a function to use with mock.Mock.
var Match = mock.MatchedBy(Is)
