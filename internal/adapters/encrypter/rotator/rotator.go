package rotator

import (
	"fmt"
	"time"
)

type Rotate func() (key []byte, keyEncr []byte, err error)

type Rotator struct {
	duration time.Duration
	rotate   Rotate

	expiredTime  time.Time
	key          []byte
	keyEncrypted []byte
}

func New(duration time.Duration, rotate Rotate) (*Rotator, error) {
	key, keyEncr, err := rotate()
	if err != nil {
		return nil, fmt.Errorf("rotator - New - rotate: %w", err)
	}
	return &Rotator{duration: duration, rotate: rotate, expiredTime: time.Now().Add(duration), key: key, keyEncrypted: keyEncr}, nil
}

func (r *Rotator) Keys() (key []byte, keyEncr []byte, err error) {
	if time.Now().After(r.expiredTime) {
		k, ke, err := r.rotate()
		if err != nil {
			return nil, nil, fmt.Errorf("rotator - Keys - rotate: %w", err)
		}
		r.key, r.keyEncrypted = k, ke
		r.expiredTime = time.Now().Add(r.duration)
	}
	return r.key, r.keyEncrypted, nil
}
