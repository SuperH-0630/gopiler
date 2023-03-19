package main

import (
	"os"
	"path"
)

func copyVersionFile() (err error) {
	if len(ico) != 0 {
		err = copyFile(ico, path.Join(tmpPath, Base(ico)), getFileMode(ico))
		if err != nil {
			return err
		}
	}

	code := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
 <assemblyIdentity
   type="win32"
   name="Github.com.JosephSpurrier.GoVersionInfo"
   version="1.0.0.0"
   processorArchitecture="*"/>
<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
  <security>
    <requestedPrivileges>
      <requestedExecutionLevel
        level="asInvoker"
        uiAccess="false"/>
      </requestedPrivileges>
  </security>
</trustInfo>
`

	err = os.WriteFile(path.Join(tmpPath, "goversioninfo.exe.manifest"), []byte(code), 0666)
	if err != nil {
		return err
	}

	versioninfo := path.Join(projectPath, "versioninfo.json")
	err = copyFile(versioninfo, path.Join(tmpPath, "versioninfo.json"), getFileMode(versioninfo))
	if err != nil {
		return err
	}

	return
}
