package trbac

import (
    "testing"
    "os"
    "io/ioutil"
)

func TestTOMLPrivileges(t *testing.T) {
	tomlData := `
reader = [
  { actions = ["read"], resource_types = ["document"], constraints = [] }
]

writer = [
  { actions = ["read", "write"], resource_types = ["document"], constraints = ["business_hours"] }
]
`

    tmpFile, err := ioutil.TempFile("", "privileges*.toml")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpFile.Name())

    if _, err := tmpFile.Write([]byte(tomlData)); err != nil {
        t.Fatal(err)
    }
    if err := tmpFile.Close(); err != nil {
        t.Fatal(err)
    }

    privs, err := NewTOMLPrivileges(tmpFile.Name())
    if err != nil {
        t.Fatalf("Failed to load privileges: %v", err)
    }

    readerPerms := privs.GetPermissions([]string{"reader"})
    if len(readerPerms) != 1 || readerPerms[0].Actions[0] != "read" {
        t.Error("Expected reader to have read permission on document")
    }

    writerPerms := privs.GetPermissions([]string{"writer"})
    if len(writerPerms) != 1 || writerPerms[0].Actions[1] != "write" {
        t.Error("Expected writer to have write permission on document")
    }
}
