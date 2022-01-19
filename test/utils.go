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

package testing

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

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
func GetListOfPids() []string {
	os.Setenv("INSTALL_PATH", "../../supertokens-root")
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
func SetUpST() {
	os.Setenv("INSTALL_PATH", "../../supertokens-root")
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
func StartUpST(host string, port string) {
	os.Setenv("INSTALL_PATH", "../../supertokens-root")
	defer os.Unsetenv("INSTALL_PATH")
	installationPath := os.Getenv("INSTALL_PATH") //*---> ../supertokens-root
	// pidsBefore := GetListOfPids()
	// returned := false
	command := fmt.Sprintf(`java -Djava.security.egd=file:/dev/urandom -classpath "./core/*:./plugin-interface/*" io.supertokens.Main ./ DEV host=%s port=%s test_mode`, host, port)

	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = installationPath
	err := cmd.Run()
	if err != nil {
		log.Fatalf(err.Error(), "could not initiate a supertokens instance")
	}
	// startTime := time.Now().Unix()
	// for (time.Now().Unix() - startTime) < 3 {
	// 	pidsAfter := GetListOfPids()
	// 	if len(pidsAfter) <= len(pidsBefore) {
	// 		time.Sleep(10 * time.Millisecond)
	// 		continue
	// 	}
	// 	var nonIntersection = []string{""}
	// 	for i := 0; i < len(pidsAfter); i++ {
	// 		for j := 0; j < len(pidsBefore); j++ {
	// 			if pidsAfter[i] != pidsBefore[j] {
	// 				nonIntersection = append(nonIntersection, pidsAfter[i])
	// 			}
	// 		}
	// 	}
	// 	if len(nonIntersection) != 1 {
	// 		if !returned {
	// 			returned = true
	// 			log.Fatalf("something went wrong while starting up the core")
	// 		}
	// 	} else {
	// 		if !returned {
	// 			returned = true
	// 			if len(nonIntersection) != 0 {
	// 				// return nonIntersection[0]
	// 			}
	// 		}
	// 	}
	// }
	// if !returned {
	// 	returned = true
	// 	log.Fatalf("something went wrong while starting up the core")
	// }
}

//*this kills a running instance of the root by taking it's pid
func StopST(pid string) {
	os.Setenv("INSTALL_PATH", "../../supertokens-root")
	defer os.Unsetenv("INSTALL_PATH")
	installationPath := os.Getenv("INSTALL_PATH") //*---> ../supertokens-root
	pidsBefore := GetListOfPids()
	if len(pidsBefore) == 0 {
		return
	}
	cmd := exec.Command("kill", pid)
	cmd.Dir = installationPath
	err := cmd.Run()
	if err != nil {
		log.Fatalf(err.Error(), "could not close the supertokens instance%s", pid)
	}
	// startTime := time.Now().Unix()
	// for (time.Now().Unix() - startTime) < 3 {
	// 	pidsAfter := GetListOfPids()
	// 	includes := false
	// 	for i := 0; i < len(pidsAfter); i++ {
	// 		if pidsAfter[i] == pid {
	// 			includes = true
	// 		}
	// 	}
	// 	if includes {
	// 		time.Sleep(10 * time.Millisecond)
	// 		continue
	// 	} else {
	// 		return
	// 	}
	// }
	// log.Fatalf(err.Error(), "error while stopping st with pid%s", pid)
}

func KillAllSTCoresOnly() {
	pids := GetListOfPids()
	for i := 0; i < len(pids); i++ {
		StopST(pids[i])
	}
}

//*this function cleans up all the files that were required to run the root instance and should be called after closing the pid
func CleanST() {
	os.Setenv("INSTALL_PATH", "../../supertokens-root")
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
func ResetAll() {
	supertokens.ResetForTest()
	emailpassword.ResetForTest()
	session.ResetForTest()
	thirdparty.ResetForTest()
	thirdpartyemailpassword.ResetForTest()
	openid.ResetForTest()
	passwordless.ResetForTest()
	jwt.ResetForTest()
}

func KillAllST() {
	pids := GetListOfPids()
	for i := 0; i < len(pids); i++ {
		StopST(pids[i])
	}
	ResetAll()
}
