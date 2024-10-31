package metainfo

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	wd, _ := os.Getwd()
	testDir := filepath.Join(wd, "..", "..", "test")

	testFilepath := filepath.Join(testDir, "test.torrent")
	expectedAnnounce := "https://torrent.ubuntu.com/announce"
	expectedInfoName := "ubuntu-24.04.1-desktop-amd64.iso"
	expectedInfoLength := int64(6203355136)
	expectedCreationDate := time.Unix(1724947415, 0)

	t.Run("ubuntu", func(t *testing.T) {
		f, err := os.Open(testFilepath)
		if err != nil {
			t.Fatal(err)
		}
		m, err := Parse(bufio.NewReader(f))
		if err != nil {
			t.Fatal(err)
		}
		if m.Announce != expectedAnnounce {
			t.Errorf("expected %q, got %q", expectedAnnounce, m.Announce)
		}
		if m.Info.Name != expectedInfoName {
			t.Errorf("expected %q, got %q", expectedInfoName, m.Info.Name)
		}
		if m.Info.Length != expectedInfoLength {
			t.Errorf("expected %d, got %d", expectedInfoLength, m.Info.Length)
		}
		if m.CreationDate != expectedCreationDate {
			t.Errorf("expected %v, got %v", expectedCreationDate, m.CreationDate)
		}
	})
}
