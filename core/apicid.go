package core

import (
	"encoding/json"

	cidenc "github.com/ipfs/go-ipfs/core/cidenc"
	cid "gx/ipfs/QmPSQnBKM9g7BaUcZCvswUJVscQ1ipjmwxN5PXCjkp9EQ7/go-cid"
	mbase "gx/ipfs/QmekxXDhCxCJRNuzmHreuaT3BsuJcsjcXWNrtV9C8DRHtd/go-multibase"
)

// CidJSONBase is the base to use when Encoding into JSON.
//var CidJSONBase mbase.Encoder = mbase.MustNewEncoder(mbase.Base58BTC)
var CidJSONBase mbase.Encoder = mbase.MustNewEncoder(mbase.Base32)

// APICid is a type to respesnt CID in the API
type APICid struct {
	str string // always in CidJSONBase
}

// FromCid created an APICid from a Cid
func FromCid(c cid.Cid) APICid {
	return APICid{c.Encode(CidJSONBase)}
}

// Cid converts an APICid to a CID
func (c APICid) Cid() (cid.Cid, error) {
	return cid.Decode(c.str)
}

func (c APICid) String() string {
	return c.Encode(cidenc.Default)
}

func (c APICid) Encode(enc cidenc.Interface) string {
	if c.str == "" {
		return ""
	}
	str, err := enc.Recode(c.str)
	if err != nil {
		return c.str
	}
	return str
}

func (c *APICid) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &c.str)
}

func (c APICid) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.str)
}
