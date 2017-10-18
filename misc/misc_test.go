package misc

import "testing"

func Test_Imax(t *testing.T) {

	x := Imax(2, 3)
	if x != 3 {
		t.Errorf("Error: imax returned wrong value\n")
	}
	x = Imax(4, 3)
	if x != 4 {
		t.Errorf("Error: imax returned wrong value\n")
	}

}

func Test_SetDebug(t *testing.T) {

	SetDebugFlags("aa,bb,cc")

	if !Db["aa"] {
		t.Errorf("Error: SetDebugFlags - should be true\n")
	}
	if !Db["bb"] {
		t.Errorf("Error: SetDebugFlags - should be true\n")
	}
	if !Db["cc"] {
		t.Errorf("Error: SetDebugFlags - should be true\n")
	}
	if Db["dd"] {
		t.Errorf("Error: SetDebugFlags - should be false\n")
	}

}

func Test_SVar(t *testing.T) {
	mdata := make(map[string]string)
	mdata["abc"] = "def"
	s := SVar(mdata)
	if s != `{"abc":"def"}` {
		t.Errorf("Error: SVar\n")
	}
}

func Test_SVarI(t *testing.T) {
	mdata := make(map[string]string)
	mdata["abc"] = "def"
	s := SVarI(mdata)
	if s != `{
	"abc": "def"
}` {
		t.Errorf("Error: SVar\n")
	}
}

/* vim: set noai ts=4 sw=4: */
