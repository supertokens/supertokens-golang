/*
 * Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package unittesting

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/jwt"
	"github.com/supertokens/supertokens-golang/recipe/openid"
	"github.com/supertokens/supertokens-golang/recipe/passwordless"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/supertokens"
)

//*returns the list of all the pids inside all the running instance files inside the .started directory of the root
func getListOfPids() []string {
	slashesNeededToGoUp := returnNumberOfDirsToGoUpFromCurrentWorkingDir()
	os.Setenv("INSTALL_PATH", slashesNeededToGoUp+"supertokens-root")
	defer os.Unsetenv("INSTALL_PATH")
	installationPath := os.Getenv("INSTALL_PATH") //*---> ../supertokens-root
	pathOfDirToRead := installationPath + "/.started/"
	files, err := ioutil.ReadDir(pathOfDirToRead)
	if err != nil {
		return []string{}
	}
	var result []string
	for _, file := range files {
		pathOfFileToBeRead := installationPath + "/.started/" + file.Name()
		data, err := ioutil.ReadFile(pathOfFileToBeRead)
		if err != nil {
			log.Fatalf(err.Error())
		}
		result = append(result, string(data))
	}
	return result
}

//*this copies the config.yaml from temp folder of the root and puts it into the root level.
func setUpST() {
	slashesNeededToGoUp := returnNumberOfDirsToGoUpFromCurrentWorkingDir()
	os.Setenv("INSTALL_PATH", slashesNeededToGoUp+"supertokens-root")
	defer os.Unsetenv("INSTALL_PATH")
	installationPath := os.Getenv("INSTALL_PATH") //*---> ../supertokens-root
	cmd := exec.Command("cp", "temp/config.yaml", "./config.yaml")
	cmd.Dir = installationPath
	err := cmd.Run()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

//*this runs the java command to start the root in testing environemnt
func startUpST(host string, port string) {
	slashesNeededToGoUp := returnNumberOfDirsToGoUpFromCurrentWorkingDir()
	os.Setenv("INSTALL_PATH", slashesNeededToGoUp+"supertokens-root")
	defer os.Unsetenv("INSTALL_PATH")
	installationPath := os.Getenv("INSTALL_PATH") //*---> ../supertokens-root
	command := fmt.Sprintf(`java -Djava.security.egd=file:/dev/urandom -classpath "./core/*:./plugin-interface/*" io.supertokens.Main ./ DEV host=%s port=%s test_mode`, host, port)

	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = installationPath
	err := cmd.Run()
	if err != nil {
		log.Fatalf(err.Error(), "could not initiate a supertokens instance")
	}
}

//*this kills a running instance of the root by taking it's pid
func stopST(pid string) {
	slashesNeededToGoUp := returnNumberOfDirsToGoUpFromCurrentWorkingDir()
	os.Setenv("INSTALL_PATH", slashesNeededToGoUp+"supertokens-root")
	defer os.Unsetenv("INSTALL_PATH")
	installationPath := os.Getenv("INSTALL_PATH") //*---> ../supertokens-root
	pidsBefore := getListOfPids()
	if len(pidsBefore) == 0 {
		return
	}
	cmd := exec.Command("kill", pid)
	cmd.Dir = installationPath
	err := cmd.Run()
	if err != nil {
		log.Fatalf(err.Error(), "could not close the supertokens instance%s", pid)
	}
}

func killAllSTCoresOnly() {
	pids := getListOfPids()
	for i := 0; i < len(pids); i++ {
		stopST(pids[i])
	}
}

//*this function cleans up all the files that were required to run the root instance and should be called after closing the pid
func cleanST() {
	slashesNeededToGoUp := returnNumberOfDirsToGoUpFromCurrentWorkingDir()
	os.Setenv("INSTALL_PATH", slashesNeededToGoUp+"supertokens-root")
	defer os.Unsetenv("INSTALL_PATH")
	installationPath := os.Getenv("INSTALL_PATH") //*---> ../supertokens-root
	cmd := exec.Command("rm", "config.yaml")
	cmd.Dir = installationPath
	err := cmd.Run()
	if err != nil {
		log.Fatalf(err.Error(), "could not delete the config yaml file")
	}

	cmd = exec.Command("rm", "-rf", ".webserver-temp-*")
	cmd.Dir = installationPath
	err = cmd.Run()
	if err != nil {
		log.Fatalf(err.Error(), "could not delete the webserver-temp files")
	}

	cmd = exec.Command("rm", "-rf", ".started")
	cmd.Dir = installationPath
	err = cmd.Run()
	if err != nil {
		log.Fatalf(err.Error(), "could not delete the .started file")
	}
}

//*eliminates the signelton instance for all recipes
func resetAll() {
	supertokens.ResetForTest()
	emailpassword.ResetForTest()
	session.ResetForTest()
	thirdparty.ResetForTest()
	thirdpartyemailpassword.ResetForTest()
	openid.ResetForTest()
	passwordless.ResetForTest()
	jwt.ResetForTest()
}

func killAllST() {
	pids := getListOfPids()
	for i := 0; i < len(pids); i++ {
		stopST(pids[i])
	}
	resetAll()
}

func StartingHelper() {
	killAllST()
	setUpST()
	startUpST("localhost", "8080")
}

func EndingHelper() {
	resetAll()
	killAllST()
	cleanST()
}

//*returns "../../../../" if needed to go 4 levels up
func returnNumberOfDirsToGoUpFromCurrentWorkingDir() string {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	arr := strings.Split(mydir, "/")
	counter := 0
	for i := 0; i < len(arr); i++ {
		if arr[i] == "supertokens-golang" {
			counter = i
			break
		}
	}
	numberOfElems := len(arr) - counter
	var dirUpSlash string
	for i := 0; i < numberOfElems; i++ {
		dirUpSlash += "../"
	}
	return dirUpSlash
}

func RemoveTrailingSlashFromTheEndOfString(input string) string {
	if input[len(input)-1:] == "/" {
		res := input[:len(input)-1] + ""
		return res
	} else {
		return input
	}
}

func ExtractInfoFromResponse(res *http.Response) map[string]string {
	antiCsrf := res.Header["Anti-Csrf"]
	cookies := res.Header["Set-Cookie"]
	var refreshToken string
	var idRefreshTokenFromCookie string
	for _, cookie := range cookies {
		if strings.Split(strings.Split(cookie, ";")[0], "=")[0] == "sRefreshToken" {
			refreshToken = strings.Split(strings.Split(cookie, ";")[0], "=")[1]
		} else if strings.Split(strings.Split(cookie, ";")[0], "=")[0] == "sIdRefreshToken" {
			idRefreshTokenFromCookie = strings.Split(strings.Split(cookie, ";")[0], "=")[1]
		} else {
			continue
		}
	}

	return map[string]string{
		"antiCsrf":        antiCsrf[0],
		"sRefreshToken":   refreshToken,
		"sIdRefreshToken": idRefreshTokenFromCookie,
	}
}
