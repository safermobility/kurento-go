package kurento

type VideoInfo struct {
	IsSeekable   bool
	SeekableInit int64
	SeekableEnd  int64
	Duration     int64
}
