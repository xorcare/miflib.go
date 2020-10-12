// Copyright (c) 2020 Vasiliy Vasilyuk All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package book

import (
	"encoding/json"
	"fmt"
)

var _ json.Unmarshaler = (*Addresses)(nil)

// Addresses just a set of addresses with the ability to parse not a set,
// but one address as such cases are in the API.
type Addresses []Address

// UnmarshalJSON implements json.Unmarshaler.
func (a *Addresses) UnmarshalJSON(bs []byte) error {
	if a == nil {
		return fmt.Errorf("%T: UnmarshalJSON on nil pointer", a)
	}

	address := Address{}
	if err := json.Unmarshal(bs, &address); err == nil {
		*a = append((*a)[:], address)
		return nil
	}

	addresses := make([]Address, 0, 1)
	err := json.Unmarshal(bs, &addresses)
	*a = append((*a)[:], Addresses(addresses)...)

	return err
}
