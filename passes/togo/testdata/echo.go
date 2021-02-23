package translated

import "github.com/ctriv/php/passes/togo/internal/phpctx"

func Echo(ctx phpctx.PHPContext) {
	ctx.Echo.Write("test")
}
