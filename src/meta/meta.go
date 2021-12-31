package meta

type FileMeta struct {
	FileMd5   string
	FileName  string
	FileSize  int64
	Localtion string
	AddAt     string
}

var FileMetaMap map[string]FileMeta

func init() {
	FileMetaMap = make(map[string]FileMeta)
}

func UpdateFileMeta(meta FileMeta) {
	FileMetaMap[meta.FileMd5] = meta
}

func DeleteFileMeta(metaMd5 string) {
	delete(FileMetaMap, metaMd5)
}

func GetFileMeta(metaMd5 string) FileMeta {
	return FileMetaMap[metaMd5]
}
