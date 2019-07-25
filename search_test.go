package strato

import (
	"github.com/stretchr/testify/assert"
	"strato/proto"
	"testing"
)

func TestSearch(t *testing.T) {
	is := assert.New(t)

	t.Run("Document", func(t *testing.T) {
		doc := &Document{
			ID: "one",
			Content: "Here is a document",
		}

		docs := []*Document{doc}

		docP := &proto.Document{
			Id: "one",
			Content: "Here is a document",
		}

		docsP := []*proto.Document{docP}

		resP := &proto.SearchResults{
			Documents: docsP,
		}

		is.Equal(doc, docFromProto(docP))
		is.Equal(doc.toProto(), docP)
		is.Equal(docs, docsFromProto(docsP))
		is.Equal(docsToResults(docs), resP)
	})
}
