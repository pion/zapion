package zapion

import (
	"go.uber.org/zap"
)

func ExampleLogger() {
	zf := &ZapFactory{BaseLogger: zap.NewExample()}
	l := zf.NewLogger("scope")

	l.Error("test")
	l.Errorf("test printf %d", 1)

	// Output:
	// {"level":"error","logger":"scope","msg":"test"}
	// {"level":"error","logger":"scope","msg":"test printf 1"}
}
