package database

import "gorm.io/gorm"

/*
 *
 *   management of data security connexion
 *
 */
type SettingSecurity struct {
	ID int

	// Certificate X509
	Cert_activate bool
	Cert_list     string

	// Knowhost
	Kh_activate bool
	Kh_dns      bool
	Kh_list     string
}

// Setting secuity APPlication (Cache)
var settingSecurityAPP *SettingSecurity

func SaveSettingSecurity(record SettingSecurity) error {
	record.ID = 1
	ret := DB.Debug().Save(&record)

	if ret.Error != nil {
		return ret.Error
	}
	settingSecurityAPP.Cert_activate = record.Cert_activate
	settingSecurityAPP.Cert_list = record.Cert_list
	settingSecurityAPP.Kh_activate = record.Kh_activate
	settingSecurityAPP.Kh_dns = record.Kh_dns
	settingSecurityAPP.Kh_list = record.Kh_list
	return nil
}

func GetSettingSecurity() (*SettingSecurity, error) {
	if settingSecurityAPP == nil {
		settingSecurityAPP = &SettingSecurity{}
		ret := DB.First(settingSecurityAPP)
		if ret.Error != gorm.ErrRecordNotFound {
			if ret.Error != nil {
				return nil, ret.Error
			}
		}
	}
	return settingSecurityAPP, nil
}
