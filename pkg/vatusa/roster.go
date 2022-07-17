package vatusa

import "fmt"

func RemoveController(cid string, by uint, reason string) (int, error) {
	status, _, err := handle("DELETE", "/facility/ZDV/roster/"+cid, map[string]string{
		"by":     fmt.Sprint(by),
		"reason": reason,
	})

	return status, err
}

func RemoveVisitingController(cid string, by uint, reason string) (int, error) {
	status, _, err := handle("DELETE", "/facility/ZDV/roster/visiting/"+cid, map[string]string{
		"by":     fmt.Sprint(by),
		"reason": reason,
	})

	return status, err
}

func AddVisitingController(cid string) (int, error) {
	status, _, err := handle("POST", "/facility/KZDV/roster/manageVisitor/"+cid, nil)

	return status, err
}
