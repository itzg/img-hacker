// Copyright © 2017 Geoff Bourne <itzgeoff@gmail.com>
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

package internal

import (
	"path/filepath"
	"os"
	"strings"
	"github.com/Sirupsen/logrus"
	"image/png"
	"github.com/nfnt/resize"
)

type Resizer struct {
	from, to *Dimension
}

func Resize(from, to *Dimension) error {

	resizer := &Resizer{from:from, to:to}
	filepath.Walk(".", resizer.walker)

	return nil
}

func (r *Resizer) walker(path string, info os.FileInfo, err error) error {

	if strings.HasSuffix(strings.ToLower(info.Name()), ".png") {
		f, err := os.Open(path)
		if err != nil {
			logrus.WithError(err).WithField("path", path).Error("Unable to open file")
			return nil
		}

		img, err := png.Decode(f)
		if err != nil {
			f.Close()
			logrus.WithError(err).WithField("path", path).Error("Unable to decode as PNG")
			return nil
		}
		f.Close()

		imgMax := img.Bounds().Max
		if imgMax.X == r.from.Width && imgMax.Y == r.from.Height {
			logrus.WithField("path", path).Info("Converting")

			resizedImg := resize.Resize(uint(r.to.Width), uint(r.to.Height), img, resize.NearestNeighbor)

			outFile, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				logrus.WithError(err).WithField("path", path).Error("Unable to re-open image file for writing")
				return nil
			}
			defer outFile.Close()

			err = png.Encode(outFile, resizedImg)
			if err != nil {
				logrus.WithError(err).WithField("path", path).Warn("Unable to re-encode image file as PNG")
			}
		} else {
			logrus.WithField("path", path).Infof("Skipping since original is not %v", r.from)
		}
	}
	return nil
}