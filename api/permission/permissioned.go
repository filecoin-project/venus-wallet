package permission

import (
	"context"
	"golang.org/x/xerrors"
	"reflect"

	logging "github.com/ipfs/go-log/v2"
)

type permKey int

var (
	permCtxKey   permKey = 1
	permCtxLocal permKey = 2
)

// nolint
var log = logging.Logger("api")

type Permission = string

const (
	// When changing these, update docs/API.md too

	PermRead  Permission = "read" // default
	PermWrite Permission = "write"
	PermSign  Permission = "sign"  // Use wallet keys for signing
	PermAdmin Permission = "admin" // Manage permissions

)

var AllPermissions = []Permission{PermRead, PermWrite, PermSign, PermAdmin}
var defaultPerms = []Permission{PermRead, PermWrite}

func WithPerm(ctx context.Context, perms []Permission) context.Context {
	return context.WithValue(ctx, permCtxKey, perms)
}
func WithIPPerm(ctx context.Context) context.Context {
	return context.WithValue(ctx, permCtxLocal, true)
}

func HasPerm(ctx context.Context, perm Permission) bool {
	callerPerms, ok := ctx.Value(permCtxKey).([]Permission)
	if !ok {
		callerPerms = defaultPerms
	}

	for _, callerPerm := range callerPerms {
		if callerPerm == perm {
			return true
		}
	}
	return false
}

func HasIPPerm(ctx context.Context, required bool) bool {
	if !required {
		return true
	}
	localPerm, ok := ctx.Value(permCtxLocal).(bool)
	if !ok {
		return false
	}
	return localPerm
}
func PermissionedAny(in interface{}, out interface{}) {
	rint := reflect.ValueOf(out).Elem()
	ra := reflect.ValueOf(in)

	for f := 0; f < rint.NumField(); f++ {
		field := rint.Type().Field(f)
		requiredPerm := Permission(field.Tag.Get("perm"))
		if requiredPerm == "" {
			panic("missing 'perm' tag on " + field.Name) // ok
		}
		localTag := field.Tag.Get("local")
		required := false
		if localTag == "required" {
			required = true
		}
		// Validate perm tag
		ok := false
		for _, perm := range AllPermissions {
			if requiredPerm == perm {
				ok = true
				break
			}
		}
		if !ok {
			panic("unknown 'perm' tag on " + field.Name) // ok
		}
		fn := ra.MethodByName(field.Name)
		rint.Field(f).Set(reflect.MakeFunc(field.Type, func(args []reflect.Value) (results []reflect.Value) {
			ctx := args[0].Interface().(context.Context)
			errNum := 0
			if !HasPerm(ctx, requiredPerm) {
				errNum += 1
				goto ABORT
			}
			if !HasIPPerm(ctx, required) {
				errNum += 2
				goto ABORT
			}
			return fn.Call(args)
		ABORT:
			err := xerrors.Errorf("missing permission to invoke '%s'", field.Name)
			if errNum&1 == 1 {
				err = xerrors.Errorf("%s  (need '%s')", err, requiredPerm)
			}
			if errNum&2 == 2 {
				err = xerrors.Errorf("%s (only allow local access)", err)
			}
			rerr := reflect.ValueOf(&err).Elem()
			if field.Type.NumOut() == 2 {
				return []reflect.Value{
					reflect.Zero(field.Type.Out(0)),
					rerr,
				}
			} else {
				return []reflect.Value{rerr}
			}
		}))
	}
}
