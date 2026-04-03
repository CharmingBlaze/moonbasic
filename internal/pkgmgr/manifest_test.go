package pkgmgr

import "testing"

func TestManifestValidate(t *testing.T) {
	_, err := ParseManifest([]byte(`{"name":"bad Name","version":"1","entry_mbc":"x.mbc"}`))
	if err == nil {
		t.Fatal("expected invalid name")
	}
	m, err := ParseManifest([]byte(`{"name":"math_extra","version":"1.0.0","entry_mbc":"math_extra.mbc"}`))
	if err != nil {
		t.Fatal(err)
	}
	if m.Name != "math_extra" {
		t.Fatalf("name %q", m.Name)
	}
}
