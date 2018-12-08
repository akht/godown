package evaluator

import (
	"godown/ast"
	"godown/object"
)

func Eval(node ast.Node) object.Document {
	switch node := node.(type) {
	case *ast.Document:
		return evalDocument(node)
	}

	return object.Document{}
}

func evalDocument(document *ast.Document) object.Document {
	var evaluated object.Document

	for _, block := range document.Blocks {
		result := EvalBlock(block)

		switch result := result.(type) {
		case *object.Heading:
			evaluated.Objects = append(evaluated.Objects, result)
		case *object.DiscList:
			evaluated.Objects = append(evaluated.Objects, result)
		case *object.CodeBlock:
			evaluated.Objects = append(evaluated.Objects, result)
		case *object.Paragraph:
			evaluated.Objects = append(evaluated.Objects, result)
		case *object.HorizontalRule:
			evaluated.Objects = append(evaluated.Objects, result)
		}
	}

	return evaluated
}

func EvalBlock(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Heading:
		return &object.Heading{Value: node.String()}
	case *ast.DiscList:
		return &object.DiscList{Value: node.String()}
	case *ast.CodeBlock:
		return &object.CodeBlock{Value: node.String()}
	case *ast.Paragraph:
		return &object.Paragraph{Value: node.String()}
	case *ast.HorizontalRule:
		return &object.HorizontalRule{}
	}

	return nil
}
