package translated

import "github.com/ctriv/php/passes/togo/internal/phpctx"

func Shell(ctx phpctx.PHPContext) {
	ctx.Shell(`ls -al`)
}
