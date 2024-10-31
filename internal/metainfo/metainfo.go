package metainfo

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/stupoid/torrent/internal/bencode"
)

type MetaInfo struct {
	Announce     string
	AnnounceList [][]string
	Comment      string
	CreatedBy    string
	CreationDate time.Time
	Encoding     string
	Info         Info
}

func (m MetaInfo) String() string {
	return fmt.Sprintf("MetaInfo{Announce: %s, AnnounceList: %v, Comment: %s, CreatedBy: %s, CreationDate: %s, Encoding: %s, Info: %v}", m.Announce, m.AnnounceList, m.Comment, m.CreatedBy, m.CreationDate, m.Encoding, m.Info)
}

type Info struct {
	PieceLength int64
	Pieces      [][20]byte
	Private     bool

	Name string // Filename in Single File Mode, Directory Name in Multiple File Mode

	// Only present in Single File Mode
	Length int64
	MD5Sum []byte

	// Only present in Multiple File Mode
	Files []File
}

func (i Info) String() string {
	return fmt.Sprintf("Info{PieceLength: %d, Private: %t, Name: %s, Length: %d, MD5Sum: %x, Files: %v}", i.PieceLength, i.Private, i.Name, i.Length, i.MD5Sum, i.Files)
}

type File struct {
	Length int64
	MD5Sum []byte
	Path   string
}

func (f File) String() string {
	return fmt.Sprintf("File{Length: %d, MD5Sum: %x, Path: %s}", f.Length, f.MD5Sum, f.Path)
}

func Parse(r *bufio.Reader) (*MetaInfo, error) {
	decoder := bencode.NewDecoder(r)
	dict, err := decoder.DecodeDict()
	if err != nil {
		return nil, err
	}

	metaInfo := MetaInfo{}

	if _, ok := dict["announce"]; !ok {
		return nil, errors.New("missing announce key")
	}
	announce, ok := dict["announce"].(string)
	if !ok {
		return nil, errors.New("invalid announce")
	}
	metaInfo.Announce = announce

	if announceList, ok := dict["announce-list"]; ok {
		announceList, ok := announceList.([][]string)
		if ok {
			metaInfo.AnnounceList = announceList
		}
	}

	if comment, ok := dict["comment"]; ok {
		if comment, ok := comment.(string); ok {
			metaInfo.Comment = comment
		}
	}

	if createdBy, ok := dict["created by"]; ok {
		if createdBy, ok := createdBy.(string); ok {
			metaInfo.CreatedBy = createdBy
		}
	}

	if creationDate, ok := dict["creation date"]; ok {
		if creationDate, ok := creationDate.(int64); ok {
			metaInfo.CreationDate = time.Unix(creationDate, 0)
		}
	}

	if encoding, ok := dict["encoding"]; ok {
		if encoding, ok := encoding.(string); ok {
			metaInfo.Encoding = encoding
		}
	}

	dictInfo, ok := dict["info"].(map[string]interface{})
	if !ok {
		return nil, errors.New("missing info key")
	}
	info, err := ParseInfo(dictInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse info: %w", err)
	}
	metaInfo.Info = info

	return &metaInfo, nil
}

func ParseInfo(dict map[string]interface{}) (Info, error) {
	info := Info{}

	pieceLength, ok := dict["piece length"].(int64)
	if !ok {
		return info, errors.New("missing piece length")
	}
	info.PieceLength = pieceLength

	piecesString, ok := dict["pieces"].(string)
	if !ok {
		return info, errors.New("missing pieces")
	}
	for i := 0; i < len(piecesString); i += 20 {
		var piece [20]byte
		copy(piece[:], piecesString[i:i+20])
		info.Pieces = append(info.Pieces, piece)
	}

	if private, ok := dict["private"].(int64); ok {
		info.Private = private == 1
	}

	if name, ok := dict["name"].(string); ok {
		info.Name = name
	}

	if length, ok := dict["length"].(int64); ok {
		info.Length = length

		if md5sumHexString, ok := dict["md5sum"].(string); ok {
			md5sum, err := hex.DecodeString(md5sumHexString)
			if err != nil {
				return info, errors.New("invalid md5sum")
			}
			info.MD5Sum = md5sum
		}

	} else if filesList, ok := dict["files"].([]map[string]interface{}); ok {
		// Multiple File Mode
		for _, fileDict := range filesList {
			file := File{}

			length, ok := fileDict["length"].(int64)
			if !ok {
				return info, errors.New("missing file length")
			}
			file.Length = length

			if md5sumHexString, ok := fileDict["md5sum"].(string); ok {
				md5sum, err := hex.DecodeString(md5sumHexString)
				if err != nil {
					return info, errors.New("invalid md5sum")
				}
				file.MD5Sum = md5sum
			}

			pathList, ok := fileDict["path"].([]interface{})
			if !ok {
				return info, errors.New("missing file path")
			}
			for _, pathComponent := range pathList {
				pathComponent, ok := pathComponent.(string)
				if !ok {
					return info, errors.New("invalid file path")
				}
				file.Path += pathComponent
			}

			info.Files = append(info.Files, file)
		}

	} else {
		return info, errors.New("missing any file definition")
	}

	return info, nil
}
