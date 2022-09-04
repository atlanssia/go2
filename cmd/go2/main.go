package main

import (
	"context"
	"fmt"

	"github.com/atlanssia/go2/pkg/log"
)

func main() {
	fmt.Println("go2")
	ctx := context.Background()
	log.Info(ctx, "init in %v", ctx)
	log.Sync()
}
