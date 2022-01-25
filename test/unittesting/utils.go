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

func getListOfPids() []string {
	// slashesNeededToGoUp := returnNumberOfDirsToGoUpFromCurrentWorkingDir()
	// os.Setenv("INSTALL_PATH", slashesNeededToGoUp+"supertokens-root")
	// defer os.Unsetenv("INSTALL_PATH")
	installationPath := getInstallationDir()
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

func setUpST() {
	installationPath := getInstallationDir()
	cmd := exec.Command("cp", "temp/config.yaml", "./config.yaml")
	cmd.Dir = installationPath
	err := cmd.Run()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func startUpST(host string, port string) {
	installationPath := getInstallationDir()
	command := fmt.Sprintf(`java -Djava.security.egd=file:/dev/urandom -classpath "./core/*:./plugin-interface/*" io.supertokens.Main ./ DEV host=%s port=%s test_mode`, host, port)

	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = installationPath
	err := cmd.Run()
	if err != nil {
		log.Fatalf(err.Error(), "could not initiate a supertokens instance")
	}
}

func stopST(pid string) {
	installationPath := getInstallationDir()
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

func cleanST() {
	installationPath := getInstallationDir()
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

func returnNumberOfDirsToGoUpFromCurrentWorkingDir() string {
	mydir, err := os.Getwd()
	if err != nil {
		log.Fatal(err.Error())
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

func returnNumberOfSlashesRequiredToGoToRootOfTheProject() string {
	mydir, err := os.Getwd()
	if err != nil {
		log.Fatal(err.Error())
	}
	arr := strings.Split(mydir, "/")
	counter := 0
	for i := 0; i < len(arr); i++ {
		if arr[i] == "supertokens-golang" {
			counter = i
			break
		}
	}
	numberOfElems := len(arr) - counter - 1
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
	var refreshTokenExpiry string
	var refreshTokenDomain string
	var refreshTokenHttpOnly = "false"
	var idRefreshTokenFromCookie string
	var idRefreshTokenExpiry string
	var idRefreshTokenDomain string
	var idRefreshTokenHttpOnly = "false"
	var accessToken string
	var accessTokenExpiry string
	var accessTokenDomain string
	var accessTokenHttpOnly = "false"
	for _, cookie := range cookies {
		if strings.Split(strings.Split(cookie, ";")[0], "=")[0] == "sRefreshToken" {
			refreshToken = strings.Split(strings.Split(cookie, ";")[0], "=")[1]
			if strings.Split(strings.Split(cookie, ";")[2], "=")[0] == " Expires" {
				refreshTokenExpiry = strings.Split(strings.Split(cookie, ";")[2], "=")[1]
			} else if strings.Split(strings.Split(cookie, ";")[2], "=")[0] == " expires" {
				refreshTokenExpiry = strings.Split(strings.Split(cookie, ";")[2], "=")[1]
			} else {
				refreshTokenExpiry = strings.Split(strings.Split(cookie, ";")[3], "=")[1]
			}
			if strings.Split(strings.Split(cookie, ";")[1], "=")[0] == " Path" {
				fmt.Println("need to be figured out")
			}
			for _, property := range strings.Split(cookie, ";") {
				if strings.Index(property, "HttpOnly") == 1 {
					refreshTokenHttpOnly = "true"
					break
				}
			}
		} else if strings.Split(strings.Split(cookie, ";")[0], "=")[0] == "sIdRefreshToken" {
			idRefreshTokenFromCookie = strings.Split(strings.Split(cookie, ";")[0], "=")[1]
			if strings.Split(strings.Split(cookie, ";")[2], "=")[0] == " Expires" {
				idRefreshTokenExpiry = strings.Split(strings.Split(cookie, ";")[2], "=")[1]
			} else if strings.Split(strings.Split(cookie, ";")[2], "=")[0] == " expires" {
				idRefreshTokenExpiry = strings.Split(strings.Split(cookie, ";")[2], "=")[1]
			} else {
				idRefreshTokenExpiry = strings.Split(strings.Split(cookie, ";")[3], "=")[1]
			}
			if strings.Split(strings.Split(cookie, ";")[1], "=")[0] == " Path" {
				fmt.Println("need to be figured out")
			}
			for _, property := range strings.Split(cookie, ";") {
				if strings.Index(property, "HttpOnly") == 1 {
					idRefreshTokenHttpOnly = "true"
					break
				}
			}
		} else if strings.Split(strings.Split(cookie, ";")[0], "=")[0] == "sAccessToken" {
			accessToken = strings.Split(strings.Split(cookie, ";")[0], "=")[1]
			if strings.Split(strings.Split(cookie, ";")[2], "=")[0] == " Expires" {
				accessTokenExpiry = strings.Split(strings.Split(cookie, ";")[2], "=")[1]
			} else if strings.Split(strings.Split(cookie, ";")[2], "=")[0] == " expires" {
				accessTokenExpiry = strings.Split(strings.Split(cookie, ";")[2], "=")[1]
			} else {
				accessTokenExpiry = strings.Split(strings.Split(cookie, ";")[3], "=")[1]
			}
			if strings.Split(strings.Split(cookie, ";")[1], "=")[0] == " Path" {
				fmt.Println("need to be figured out")
			}
			for _, property := range strings.Split(cookie, ";") {
				if strings.Index(property, "HttpOnly") == 1 {
					accessTokenHttpOnly = "true"
					break
				}
			}
		}
	}
	return map[string]string{
		"antiCsrf":               antiCsrf[0],
		"sAccessToken":           accessToken,
		"sRefreshToken":          refreshToken,
		"sIdRefreshToken":        idRefreshTokenFromCookie,
		"refreshTokenExpiry":     refreshTokenExpiry,
		"refreshTokenDomain":     refreshTokenDomain,
		"refreshTokenHttpOnly":   refreshTokenHttpOnly,
		"idRefreshTokenExpiry":   idRefreshTokenExpiry,
		"idRefreshTokenDomain":   idRefreshTokenDomain,
		"idRefreshTokenHttpOnly": idRefreshTokenHttpOnly,
		"accessTokenExpiry":      accessTokenExpiry,
		"accessTokenDomain":      accessTokenDomain,
		"accessTokenHttpOnly":    accessTokenHttpOnly,
	}
}

func getInstallationDir() string {
	slashesNeededToGoUp := returnNumberOfDirsToGoUpFromCurrentWorkingDir()
	installationDir := os.Getenv("INSTALL_DIR")
	if installationDir == "" {
		installationDir = slashesNeededToGoUp + "supertokens-root"
	} else {
		installationDir = returnNumberOfSlashesRequiredToGoToRootOfTheProject() + installationDir
	}
	return installationDir
}
