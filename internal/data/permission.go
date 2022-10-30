package data

const (
	PermissionProductsRead  = "products:read"
	PermissionProductsWrite = "products:write"
	PermissionProductsBuy   = "products:buy"
)

type Permissions []string

func (p Permissions) Includes(code string) bool {

	for i := range p {
		if code == p[i] {
			return true
		}
	}

	return false
}
