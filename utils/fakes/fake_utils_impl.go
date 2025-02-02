package fakes

import (
	"errors"
	"os/exec"
	"strings"
)

type FakeCfJavaPluginUtil struct {
	SshEnabled           bool
	Jmap_jvmmon_present  bool
	Container_path_valid bool
	Fspath               string
	LocalPathValid       bool
}

func (fakeUtil FakeCfJavaPluginUtil) CheckRequiredTools(app string) (bool, error) {

	if !fakeUtil.SshEnabled {
		return false, errors.New("ssh is not enabled for app: '" + app + "', please run below 2 shell commands to enable ssh and try again(please note application should be restarted before take effect):\ncf enable-ssh " + app + "\ncf restart " + app)
	}

	if !fakeUtil.Jmap_jvmmon_present {
		return false, errors.New(`jvmmon or jmap are required for generating heap dump, you can modify your application manifest.yaml on the 'JBP_CONFIG_OPEN_JDK_JRE' environment variable. This could be done like this:
		---
		applications:
		- name: <APP_NAME>
		  memory: 1G
		  path: <PATH_TO_BUILD_ARTIFACT>
		  buildpack: https://github.com/cloudfoundry/java-buildpack
		  env:
			JBP_CONFIG_OPEN_JDK_JRE: '{ jre: { repository_root: "https://java-buildpack.cloudfoundry.org/openjdk-jdk/bionic/x86_64", version: 11.+ } }'
		
		`)
	}

	return true, nil
}

func (fake FakeCfJavaPluginUtil) GetAvailablePath(data string, userpath string) (string, error) {
	if !fake.Container_path_valid && len(userpath) > 0 {
		return "", errors.New("the container path specified doesn't exist or have no read and write access, please check and try again later")
	}

	if len(fake.Fspath) > 0 {
		return fake.Fspath, nil
	}

	return "/tmp", nil
}

func (fake FakeCfJavaPluginUtil) CopyOverCat(app string, src string, dest string) error {

	if !fake.LocalPathValid {
		return errors.New("Error occured during create desination file: " + dest + ", please check you are allowed to create file in the path.")
	}

	return nil
}

func (fake FakeCfJavaPluginUtil) DeleteRemoteFile(app string, path string) error {
	_, err := exec.Command("cf", "ssh", app, "-c", "rm "+path).Output()

	if err != nil {
		return errors.New("error occured while removing dump file generated")

	}

	return nil
}

// func (fake FakeCfJavaPluginUtil) FindDumpFile(app string, path string, fspath string) (string, error) {
// 	cmd := " [ -f '" + path + "' ] && echo '" + path + "' ||  find -name 'java_pid*.hprof' -printf '%T@ %p\\0' | sort -zk 1nr | sed -z 's/^[^ ]* //' | tr '\\0' '\\n' | head -n 1  "

// 	output, err := exec.Command("cf", "ssh", app, "-c", cmd).Output()

// 	if err != nil {
// 		return "", errors.New("error occured while checking the generated file")

// 	}

// 	return strings.Trim(string(output[:]), "\n"), nil

// }

func (checker FakeCfJavaPluginUtil) FindDumpFile(app string, fullpath string, fspath string) (string, error) {
	cmd := " [ -f '" + fullpath + "' ] && echo '" + fullpath + "' ||  find " + fspath + " -name 'java_pid*.hprof' -printf '%T@ %p\\0' | sort -zk 1nr | sed -z 's/^[^ ]* //' | tr '\\0' '\\n' | head -n 1  "

	output, err := exec.Command("cf", "ssh", app, "-c", cmd).Output()

	if err != nil {
		return "", errors.New("error while checking the generated file" + (err.Error()))
	}

	return strings.Trim(string(output[:]), "\n"), nil

}
