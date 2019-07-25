package strato

import (
	"strato/proto"
	"strings"
)

type (
	Search interface {
		Index(doc *Document)
		Query(q string) []*Document
	}

	Document struct {
		ID      string
		Content string
	}
)

func (d *Document) prepare() *Document {
	return &Document{
		ID: d.ID,
		Content: strings.ToLower(d.Content),
	}
}

func (d *Document) toProto() *proto.Document {
	return &proto.Document{
		Id:      d.ID,
		Content: d.Content,
	}
}

func docFromProto(docP *proto.Document) *Document {
	return &Document{
		ID:      docP.Id,
		Content: docP.Content,
	}
}

func docsFromProto(docsP []*proto.Document) []*Document {
	docs := make([]*Document, 0)

	for _, d := range docsP {
		docs = append(docs, docFromProto(d))
	}

	return docs
}

func docsToResults(docs []*Document) *proto.SearchResults {
	docsP := make([]*proto.Document, 0)

	for _, d := range docs {
		docsP = append(docsP, d.toProto())
	}

	return &proto.SearchResults{
		Documents: docsP,
	}
}
