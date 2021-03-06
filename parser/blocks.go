package parser

import (
	"github.com/ctriv/php/ast"
	"github.com/ctriv/php/token"
)

func (p *Parser) parseBlock() *ast.Block {
	p.expect(token.BlockBegin)
	b := p.parseStatementsUntil(token.BlockEnd)
	p.expectCurrent(token.BlockEnd)
	return b
}

func (p *Parser) parseStatementsUntil(endTokens ...token.Token) *ast.Block {
	block := &ast.Block{}

	// this option exists to allow parser tests to pass while scope tests may be failing
	if !p.disableScoping {
		block.Scope = p.scope
	}

	breakTypes := map[token.Token]bool{}
	for _, Typ := range endTokens {
		breakTypes[Typ] = true
	}
	for {
		p.next()
		if _, ok := breakTypes[p.current.Typ]; ok {
			break
		}
		stmt := p.parseStmt()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		if p.current.Typ == token.EOF {
			break
		}
	}
	return block
}

func (p *Parser) parseExpressionsUntil(separator token.Token, endTokens ...token.Token) []ast.Expr {
	exprs := make([]ast.Expr, 0, 1)
	breakTypes := map[token.Token]bool{}
	for _, Typ := range endTokens {
		breakTypes[Typ] = true
	}
	p.next()
	first := true
	for {
		if _, ok := breakTypes[p.current.Typ]; ok {
			break
		} else if first {
			first = false
		} else {
			p.expectCurrent(separator)
			p.next()
		}
		expr := p.parseExpression()
		if expr == nil {
			return exprs
		}
		exprs = append(exprs, expr)
		p.next()
	}
	return exprs
}
