// +build qiniu

package libs

import "testing"

func TestUpload(t *testing.T) {
	t.Run("Qiniu Upload", func(t *testing.T) {
		key, hash, err := Upload("./test_upload.png", "test_upload.png")
		if err != nil {
			t.Errorf("TestUpload() get error %+v\n", err)
		}

		if key == "" || hash == "" {
			t.Errorf("TestUpload() get key: %s\n ,hash:%s", key, hash)
		}
	})
}
