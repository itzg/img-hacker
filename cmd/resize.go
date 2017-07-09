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

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/Sirupsen/logrus"
	"github.com/itzg/img-hacker/internal"
)

// resizeCmd represents the resize command
var resizeCmd = &cobra.Command{
	Use:   "resize",
	Short: "Resizes all PNG files in the current and sub directories",

	Run: func(cmd *cobra.Command, args []string) {
		dfrom := internal.RectangleFromXFormat(resizeCfg.from)
		if dfrom == nil {
			logrus.Fatalf("Bad 'from' %v, needs to be like 16x16", resizeCfg.from)
		}
		dto := internal.RectangleFromXFormat(resizeCfg.to)
		if dto == nil {
			logrus.Fatalf("Bad 'to' %v, needs to be like 32x32", resizeCfg.to)
		}

		err := internal.Resize(dfrom, dto)
		if err != nil {
			logrus.WithError(err).Fatalln("Failed to resize")
		}
	},
}

var resizeCfg = struct{
	from string
	to string
}{}

func init() {
	RootCmd.AddCommand(resizeCmd)

	resizeCmd.Flags().StringVar(&resizeCfg.from, "from", "16x16", "Specifies the only size of image to resize")
	resizeCmd.MarkFlagRequired("from")
	resizeCmd.Flags().StringVar(&resizeCfg.to, "to", "32x32", "Specifies the desired size of the image")
	resizeCmd.MarkFlagRequired("to")
}
