package register

import (
	"log"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

//const ...
const (
	pathRegKey      = "Path"
	HWNDBroadcast   = uintptr(0xffff)
	WMSettingChange = uintptr(0x001A)
)

func getRegEnvValue(key string) (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `System\CurrentControlSet\Control\Session Manager\Environment`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}

	defer k.Close()
	s, _, err := k.GetStringValue(key)
	return s, err
}

func saveRegEnvValue(key string, value string) error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `System\CurrentControlSet\Control\Session Manager\Environment`, registry.SET_VALUE)
	if err != nil {
		return err
	}

	defer k.Close()
	return k.SetStringValue(key, value)
}

// AppendPath ...
func AppendPath(newPath string) {
	path, err := getRegEnvValue(pathRegKey)
	if err != nil {
		log.Fatal("GetStringValue", err)
	}
	if err := saveRegEnvValue(pathRegKey, path+string(os.PathListSeparator)+newPath); err != nil {
		log.Fatal("SetStringValue", err)
	}
	syscall.NewLazyDLL("user32.dll").NewProc("SendMessageW").Call(HWNDBroadcast, WMSettingChange, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("ENVIRONMENT"))))
}
