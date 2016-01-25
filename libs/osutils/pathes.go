//   Copyright (C) 2015 Evgeny M. Safonov
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package osutils

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getHomeDir() (string, error) {
	var homeDir string

	switch runtime.GOOS {
	case "windows":
		homeDir = os.Getenv("USERPROFILE")
	default:
		homeDir = os.Getenv("HOME")
	}

	if homeDir == "" {
		return "", errors.New("No home directory found - set $HOME (or the platform equivalent).")
	}

	return homeDir, nil
}

func getHomeSubDir(subDirs ...string) (string, error) {
	if homeDir, err := getHomeDir(); err == nil {
		return filepath.Join(homeDir, filepath.Join(subDirs...)), nil
	} else {
		return "", err
	}
}

func GetConfigDir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return os.Getenv("LOCALAPPDATA"), nil
	case "darwin":
		return getHomeSubDir("Library", "Application Support")
	default:
		return getHomeSubDir(".config")
	}
}

func GetAppConfigDir(appName string) (string, error) {
	if configDir, err := GetConfigDir(); err == nil {
		configDir = filepath.Join(configDir, appName)
		return configDir, nil
	} else {
		return "", err
	}
}

func GetAppDir() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}
