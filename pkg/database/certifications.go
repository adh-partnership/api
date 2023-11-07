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
	certCache []models.Certification
	mutex     = &sync.RWMutex{}
)

func GetCertifications() []models.Certification {
	if certCache == nil {
		mutex.Lock()
		certCache = make([]models.Certification, 0)
		certs := []models.Certification{}
		if err := DB.Model(&models.Certification{}).Find(&certs).Error; err != nil {
			log.Errorf("Error getting certifications: %v", err)
			return nil
		}
		certCache = certs
		log.Infof("Populated certifications cache with %d entries: %+v", len(certCache), certCache)
		mutex.Unlock()
	}

	mutex.RLock()
	defer mutex.RUnlock()

	return certCache
}

func InvalidateCertCache() {
	mutex.Lock()
	certCache = nil
	log.Infof("Certifications cache invalidated")
	mutex.Unlock()
}

func ValidCertification(key string) bool {
	certifications := GetCertifications()
	for _, cert := range certifications {
		if cert.Name == key {
			return true
		}
	}

	return false
}
