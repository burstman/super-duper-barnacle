package csvdatabase

import "testing"

func TestGetFileList(t *testing.T) {
	// fileNames:=[]string{"My_first_shop",}
	// fileName2 :=
	list, err := GetFileList("/home/hamed/goExercices/shops")
	if err != nil {
		t.Error(err)
	}
	if len(list) != 2 {
		t.Errorf("len list must be 2 len(list)=%v", len(list))
	}
	for _, v := range list {
		t.Log(v)
	}
}
