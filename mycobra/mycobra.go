package mycobra

import (
	"fmt"
	"net"
	"reflect"

	"github.com/spf13/cobra"
)

// IfReplace 如果设置了，则更新
func IfReplace(cmd *cobra.Command, key string, value interface{}) error {
	reqType := reflect.TypeOf(value)
	if reqType.Kind() != reflect.Ptr {
		return fmt.Errorf("req type not a pointer:%v", reqType)
	}

	// fmt.Println("------1", reflect.TypeOf(value))

	if cmd.Flags().Changed(key) {
		switch v := value.(type) {
		case *string:
			out, err := cmd.Flags().GetString(key)
			if err != nil {
				return err
			}
			*v = out
		case *[]string:
			out, err := cmd.Flags().GetStringSlice(key)
			if err != nil {
				return err
			}
			*v = out
		case *uint:
			out, err := cmd.Flags().GetUint(key)
			if err != nil {
				return err
			}
			*v = out
		case *[]uint:
			out, err := cmd.Flags().GetUintSlice(key)
			if err != nil {
				return err
			}
			*v = out
		case *bool:
			out, err := cmd.Flags().GetBool(key)
			if err != nil {
				return err
			}
			*v = out
		case *[]bool:
			out, err := cmd.Flags().GetBoolSlice(key)
			if err != nil {
				return err
			}
			*v = out
		case *int:
			out, err := cmd.Flags().GetInt(key)
			if err != nil {
				return err
			}
			*v = out
		case *[]int:
			out, err := cmd.Flags().GetIntSlice(key)
			if err != nil {
				return err
			}
			*v = out
		case *net.IP:
			out, err := cmd.Flags().GetIP(key)
			if err != nil {
				return err
			}
			*v = out
		case *[]net.IP:
			out, err := cmd.Flags().GetIPSlice(key)
			if err != nil {
				return err
			}
			*v = out
		case *float64:
			out, err := cmd.Flags().GetFloat64(key)
			if err != nil {
				return err
			}
			*v = out
		}

	}

	return nil
}
