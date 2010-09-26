/*
Copyright (c) 2010, Jim Teeuwen.
All rights reserved.

This code is subject to a 1-clause BSD license.
The contents of which can be found in the LICENSE file.
*/

package xmlx

import "os"
import "io"

type ILoader interface {
	LoadUrl(string) os.Error
	LoadFile(string) os.Error
	LoadString(string) os.Error
	LoadStream(*io.Reader) os.Error
}

type ISaver interface {
	SaveFile(string) os.Error
	SaveString(string) (string, os.Error)
	SaveStream(*io.Writer) os.Error
}

type ILoaderSaver interface {
	ILoader
	ISaver
}
