package common

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

func CreateID(prefix string) string {
	str := uuid.New().String()
	str = strings.ReplaceAll(str, "-", "")
	return fmt.Sprintf("%s_%s", prefix, str)
}