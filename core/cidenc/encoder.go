package cidenc

import (
	cid "gx/ipfs/QmPSQnBKM9g7BaUcZCvswUJVscQ1ipjmwxN5PXCjkp9EQ7/go-cid"
	path "gx/ipfs/QmX7uSbkNz76yNwBhuwYwRbhihLnJqM73VTCjS3UMJud9A/go-path"
	mbase "gx/ipfs/QmekxXDhCxCJRNuzmHreuaT3BsuJcsjcXWNrtV9C8DRHtd/go-multibase"
)

// Encoder is a type used to encode or recode Cid as the user
// specifies
type Interface interface {
	Encode(c cid.Cid) string
	Recode(v string) (string, error)
}

// Basic is a basic Encoder that will encode Cid's using
// a specifed base, optionally upgrading a Cid if is Version 0
type Encoder struct {
	Base    mbase.Encoder
	Upgrade bool
}

var Default = Encoder{
	Base:    mbase.MustNewEncoder(mbase.Base58BTC),
	Upgrade: false,
}

func (enc Encoder) Encode(c cid.Cid) string {
	if enc.Upgrade && c.Version() == 0 {
		c = cid.NewCidV1(c.Type(), c.Hash())
	}
	return c.Encode(enc.Base)
}

func (enc Encoder) Recode(v string) (string, error) {
	skip, err := enc.NoopRecode(v)
	if skip || err != nil {
		return v, err
	}

	c, err := cid.Decode(v)
	if err != nil {
		return v, err
	}

	return enc.Encode(c), nil
}

func (enc Encoder) NoopRecode(v string) (bool, error) {
	if len(v) < 2 {
		return false, cid.ErrCidTooShort
	}
	ver := cidVer(v)
	skip := ver == 0 && !enc.Upgrade || ver == 1 && v[0] == byte(enc.Base.Encoding())
	return skip, nil
}

func cidVer(v string) int {
	if len(v) == 46 && v[:2] == "Qm" {
		return 0
	} else {
		return 1
	}
}

// On error an unmodified encoder will be returned so it is safe to
// ignore the error
func (enc Encoder) FromPath(p string) (Encoder, error) {
	v := extractCidString(p)
	if cidVer(v) == 0 {
		return Encoder{enc.Base, false}, nil
	}
	e, err := mbase.NewEncoder(mbase.Encoding(v[0]))
	if err != nil {
		return enc, err
	}
	return Encoder{e, true}, nil
}

func extractCidString(p string) string {
	segs := path.FromString(p).Segments()
	v := segs[0]
	if v == "ipfs" && len(segs) > 0 {
		v = segs[1]
	}
	return v
}

type WithOverride struct {
	base     Encoder
	override map[cid.Cid]string
}

func (enc Encoder) WithOverride(cids ...string) Interface {
	override := map[cid.Cid]string{}
	for _, p := range cids {
		v := p
		skip, err := enc.NoopRecode(v)
		if skip || err != nil {
			continue
		}
		c, err := cid.Decode(v)
		if err != nil {
			continue
		}
		println(v)
		override[c] = v
	}
	if len(override) == 0 {
		return enc
	}
	return WithOverride{enc, override}
}

func (enc WithOverride) Encode(c cid.Cid) string {
	v, ok := enc.override[c]
	if ok {
		return v
	}
	return enc.base.Encode(c)
}

func (enc WithOverride) Recode(v string) (string, error) {
	c, err := cid.Decode(v)
	if err != nil {
		return v, err
	}

	return enc.Encode(c), nil
}
