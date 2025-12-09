package forge

import "bytes"

type Readable interface {
	Read(reader *bytes.Reader)
}
