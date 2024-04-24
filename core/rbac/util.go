package rbac

import (
	"strconv"
)

// NormalizeRole corrects role ID for RBAC service
func NormalizeRole(r int) string {
	return "r" + strconv.Itoa(r)
}

// DenormalizeRole converts RBAC role back to normal
func DenormalizeRole(r string) int {
	iRole, _ := strconv.Atoi(r[1:])
	return iRole
}

// NormalizeUser corrects user ID for RBAC service
func NormalizeUser(uid int) string {
	return strconv.Itoa(uid)
}
