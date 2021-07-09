package ble

import (
	"context"
	"strings"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

func scan(ctx context.Context) (<-chan Peripheral, error) {
	d, err := linux.NewDevice()
	if err != nil {
		return nil, err
	}
	defer d.Stop()
	ble.SetDefaultDevice(d)

	ch := make(chan Peripheral)
	go func() {
		defer close(ch)
		err = ble.Scan(ctx, false, func(a ble.Advertisement) {
			ch <- Peripheral{
				Address:          strings.ToLower(a.Addr().String()),
				Name:             a.LocalName(),
				RSSI:             int16(a.RSSI()),
				ManufacturerData: a.ManufacturerData(),
			}
		}, nil)
		if err != nil {
		}
	}()

	return ch, nil
}