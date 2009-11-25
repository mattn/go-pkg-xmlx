package xmlx

import "os"
import "io"

type ILoader interface {
	LoadUrl(string) os.Error;
	LoadFile(string) os.Error;
	LoadString(string) os.Error;
	LoadStream(*io.Reader) os.Error;
}

type ISaver interface {
	SaveFile(string) os.Error;
	SaveString(string) (string, os.Error);
	SaveStream(*io.Writer) os.Error;
}

type ILoaderSaver interface {
	ILoader;
	ISaver;
}
