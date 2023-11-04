/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package database

import (
	"sync"

	"github.com/adh-partnership/api/pkg/database/models"
)

var (
	certCache []string
	mutex     = &sync.Mutex{}
)

func GetCertifications() []string {
	mutex.Lock()
	if certCache == nil {
		certCache = make([]string, 0)
		certs := []models.Certification{}
		if err := DB.Model(&models.Certification{}).Find(&certs).Error; err != nil {
			log.Errorf("Error getting certifications: %v", err)
			return nil
		}
		for _, cert := range certs {
			certCache = append(certCache, cert.Name)
		}
	}
	mutex.Unlock()

	return certCache
}

func InvalidateCertCache() {
	mutex.Lock()
	certCache = nil
	mutex.Unlock()
}

func ValidCertification(key string) bool {
	certifications := GetCertifications()
	for _, cert := range certifications {
		if cert == key {
			return true
		}
	}

	return false
}
