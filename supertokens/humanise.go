/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
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

package supertokens

import (
	"fmt"
	"math"
)

func formatOneDecimalFloat(n float64) (string, string) {
	n = math.Floor(n*10) / 10
	if float64(int(n)) == n {
		if n == 1.0 {
			return fmt.Sprintf("%d", int(n)), ""
		} else {
			return fmt.Sprintf("%d", int(n)), "s"
		}
	}
	return fmt.Sprintf("%.1f", n), "s"
}

func HumaniseSeconds(t uint64) string {
	var suffix string = ""
	if t < 60 {
		if t > 1 {
			suffix = "s"
		}
		return fmt.Sprintf("%d second%s", t, suffix)
	} else if t < 3600 {
		if t/60 > 1 {
			suffix = "s"
		}
		return fmt.Sprintf("%d minute%s", t/60, suffix)
	}
	if t/3600 > 1 {
		suffix = "s"
	}
	h := float64(t) / 3600
	hStr, suffix := formatOneDecimalFloat(h)
	return fmt.Sprintf("%s hour%s", hStr, suffix)
}

func HumaniseMilliseconds(t uint64) string {
	return HumaniseSeconds(t / 1000)
}
