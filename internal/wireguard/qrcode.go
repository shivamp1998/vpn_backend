package wireguard

import (
	"encoding/base64"

	"github.com/skip2/go-qrcode"
)

func GeneateQRCode(config string) (string, error) {
	qr, err := qrcode.New(config, qrcode.Medium)

	if err != nil {
		return "", err
	}

	pngBytes, err := qr.PNG(256)
	if err != nil {
		return "", err
	}

	base64String := base64.StdEncoding.EncodeToString(pngBytes)
	return base64String, nil
}
