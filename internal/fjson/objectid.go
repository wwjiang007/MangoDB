// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fjson

import (
	"bytes"
	"encoding/hex"
	"encoding/json"

	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
)

// fjsonObjectID represents BSON fjsonObjectID data type.
type fjsonObjectID types.ObjectID

// fjsontype implements fjsontype interface.
func (obj *fjsonObjectID) fjsontype() {}

type objectIDJSON struct {
	O string `json:"$o"`
}

// UnmarshalJSON implements fjsontype interface.
func (obj *fjsonObjectID) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		panic("null data")
	}

	r := bytes.NewReader(data)
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var o objectIDJSON
	if err := dec.Decode(&o); err != nil {
		return lazyerrors.Error(err)
	}
	if err := checkConsumed(dec, r); err != nil {
		return lazyerrors.Error(err)
	}

	b, err := hex.DecodeString(o.O)
	if err != nil {
		return lazyerrors.Error(err)
	}
	if len(b) != 12 {
		return lazyerrors.Errorf("fjson.ObjectID.UnmarshalJSON: %d bytes", len(b))
	}
	copy(obj[:], b)

	return nil
}

// MarshalJSON implements fjsontype interface.
func (obj *fjsonObjectID) MarshalJSON() ([]byte, error) {
	res, err := json.Marshal(objectIDJSON{
		O: hex.EncodeToString(obj[:]),
	})
	if err != nil {
		return nil, lazyerrors.Error(err)
	}
	return res, nil
}

// check interfaces
var (
	_ fjsontype = (*fjsonObjectID)(nil)
)