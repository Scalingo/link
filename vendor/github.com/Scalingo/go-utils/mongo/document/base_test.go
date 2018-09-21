package document

import (
	"context"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const DocsCollection = "docs"

type Doc struct {
	Base `bson:",inline"`
	Data int `bson:"data,omitempty"`
}

func NewTestDoc(t *testing.T) (*Doc, func()) {
	return NewTestDocWithData(t, 0)
}

func NewTestDocWithData(t *testing.T, data int) (*Doc, func()) {
	d := Doc{
		Data: data,
	}
	require.NoError(t, Save(context.Background(), DocsCollection, &d))
	return &d, func() {
		require.NoError(t, ReallyDestroy(context.Background(), DocsCollection, &d))
	}
}

func TestBase_Save(t *testing.T) {
	examples := []struct {
		Name   string
		Doc    func(t *testing.T) (*Doc, func())
		Expect func(t *testing.T, d *Doc)
		Error  string
	}{
		{
			Name: "it should define an ID",
			Doc: func(t *testing.T) (*Doc, func()) {
				return NewTestDoc(t)
			},
			Expect: func(t *testing.T, d *Doc) {
				require.NotEmpty(t, d.ID)
			},
		}, {
			Name: "it should not replace an ID",
			Doc: func(t *testing.T) (*Doc, func()) {
				return NewTestDoc(t)
			},
			Expect: func(t *testing.T, d *Doc) {
				id := d.ID
				err := Save(context.Background(), DocsCollection, d)
				require.NoError(t, err)
				require.Equal(t, id, d.ID)
			},
		}, {
			Name: "it should define created_at",
			Doc: func(t *testing.T) (*Doc, func()) {
				return NewTestDoc(t)
			},
			Expect: func(t *testing.T, d *Doc) {
				require.NotZero(t, d.CreatedAt)
			},
		}, {
			Name: "it should not replace created_at",
			Doc: func(t *testing.T) (*Doc, func()) {
				return NewTestDoc(t)
			},
			Expect: func(t *testing.T, d *Doc) {
				ts := d.CreatedAt
				err := Save(context.Background(), DocsCollection, d)
				require.NoError(t, err)
				require.Equal(t, ts, d.CreatedAt)
			},
		}, {
			Name: "it should define updated_at",
			Doc: func(t *testing.T) (*Doc, func()) {
				return NewTestDoc(t)
			},
			Expect: func(t *testing.T, d *Doc) {
				require.NotZero(t, d.UpdatedAt)
			},
		}, {
			Name: "it should update updated_at",
			Doc: func(t *testing.T) (*Doc, func()) {
				return NewTestDoc(t)
			},
			Expect: func(t *testing.T, d *Doc) {
				ts := d.UpdatedAt
				err := Save(context.Background(), DocsCollection, d)
				require.NoError(t, err)
				require.NotEqual(t, ts, d.UpdatedAt)
			},
		},
	}
	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			fixtureDoc, clean := example.Doc(t)
			defer clean()
			example.Expect(t, fixtureDoc)
		})
	}
}

func TestBase_Find(t *testing.T) {
	examples := []struct {
		Name  string
		Doc   func(t *testing.T) (*Doc, func())
		Error string
	}{
		{
			Name: "it should find existing doc",
			Doc: func(t *testing.T) (*Doc, func()) {
				return NewTestDoc(t)
			},
		}, {
			Name: "it should not find unsaved doc",
			Doc: func(t *testing.T) (*Doc, func()) {
				d := &Doc{}
				d.ID = bson.NewObjectId()
				return d, func() {}
			},
			Error: "not found",
		}, {
			Name: "it should not find destroyed doc",
			Doc: func(t *testing.T) (*Doc, func()) {
				d, _ := NewTestDoc(t)
				err := Destroy(context.Background(), DocsCollection, d)
				if err != nil {
					require.NoError(t, err)
				}
				return d, func() {}
			},
			Error: "not found",
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			fixtureDoc, clean := example.Doc(t)
			defer clean()

			var d Doc
			err := Find(context.Background(), DocsCollection, fixtureDoc.ID, &d)
			if example.Error != "" {
				assert.Contains(t, err.Error(), example.Error)
			} else {
				require.NoError(t, err)
				require.Equal(t, fixtureDoc.ID, d.ID)
			}
		})
	}
}

func TestBase_FindOne_WithSort(t *testing.T) {
	examples := map[string]struct {
		Docs     func(t *testing.T) ([]*Doc, func())
		Expected int
		Sort     SortField
	}{
		"with only one document": {
			Docs: func(t *testing.T) ([]*Doc, func()) {
				var docs []*Doc
				doc, clean := NewTestDocWithData(t, 10)
				return append(docs, doc), clean
			},
			Expected: 10,
			Sort:     SortField("data"),
		},
		"with three document sort positive": {
			Docs: func(t *testing.T) ([]*Doc, func()) {
				var docs []*Doc
				var cleans []func()
				for i := 1; i < 4; i++ {
					doc, clean := NewTestDocWithData(t, i)
					cleans = append(cleans, clean)
					docs = append(docs, doc)
				}

				return docs, func() {
					for _, clean := range cleans {
						clean()
					}
				}
			},
			Expected: 1,
			Sort:     SortField("data"),
		},
		"with three document sort negative": {
			Docs: func(t *testing.T) ([]*Doc, func()) {
				var docs []*Doc
				var cleans []func()
				for i := 1; i < 4; i++ {
					doc, clean := NewTestDocWithData(t, i)
					cleans = append(cleans, clean)
					docs = append(docs, doc)
				}

				return docs, func() {
					for _, clean := range cleans {
						clean()
					}
				}
			},
			Expected: 3,
			Sort:     SortField("-data"),
		},
	}

	for name, example := range examples {
		t.Run(name, func(t *testing.T) {
			_, clean := example.Docs(t)
			defer clean()

			var doc Doc
			FindOne(context.Background(), DocsCollection, bson.M{}, &doc, example.Sort)
			assert.Equal(t, example.Expected, doc.Data)

		})
	}
}

func TestBase_FindUnscoped(t *testing.T) {
	examples := []struct {
		Name  string
		Doc   func(t *testing.T) (*Doc, func())
		Error string
	}{
		{
			Name: "it should find existing doc",
			Doc: func(t *testing.T) (*Doc, func()) {
				return NewTestDoc(t)
			},
		}, {
			Name: "it should not find unsaved doc",
			Doc: func(t *testing.T) (*Doc, func()) {
				d := &Doc{}
				d.ID = bson.NewObjectId()
				return d, func() {}
			},
			Error: "not found",
		}, {
			Name: "it should not find destroyed doc",
			Doc: func(t *testing.T) (*Doc, func()) {
				d := Doc{}
				err := Save(context.Background(), DocsCollection, &d)
				require.NoError(t, err)
				err = Destroy(context.Background(), DocsCollection, &d)
				if err != nil {
					require.NoError(t, err)
				}
				return &d, func() {}
			},
			Error: "not found",
		}, {
			Name: "it should not find really destroyed doc",
			Doc: func(t *testing.T) (*Doc, func()) {
				d := Doc{}
				err := Save(context.Background(), DocsCollection, &d)
				require.NoError(t, err)
				err = ReallyDestroy(context.Background(), DocsCollection, &d)
				if err != nil {
					require.NoError(t, err)
				}
				return &d, func() {}
			},
			Error: "not found",
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			fixtureDoc, clean := example.Doc(t)
			defer clean()

			var d Doc
			err := FindUnscoped(context.Background(), DocsCollection, fixtureDoc.ID, &d)
			if example.Error != "" {
				assert.Contains(t, err.Error(), example.Error)
			} else {
				require.NoError(t, err)
				require.Equal(t, fixtureDoc.ID, d.ID)
			}
		})
	}
}

func TestBase_Where(t *testing.T) {
	examples := []struct {
		Name  string
		Query bson.M
		Docs  func(t *testing.T) ([]*Doc, func())
		Count int
	}{
		{
			Name: "it should find existing documents",
			Docs: func(t *testing.T) ([]*Doc, func()) {
				d1, clean1 := NewTestDoc(t)
				d2, clean2 := NewTestDoc(t)
				return []*Doc{d1, d2}, func() {
					clean1()
					clean2()
				}
			},
			Count: 2,
		}, {
			Name: "it should not find deleted documents",
			Docs: func(t *testing.T) ([]*Doc, func()) {
				d1, _ := NewTestDoc(t)
				err := Destroy(context.Background(), DocsCollection, d1)
				require.NoError(t, err)
				d2, _ := NewTestDoc(t)
				err = Destroy(context.Background(), DocsCollection, d2)
				require.NoError(t, err)
				return []*Doc{d1, d2}, func() {}
			},
			Count: 0,
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			_, clean := example.Docs(t)
			defer clean()

			query := bson.M{}
			if example.Query != nil {
				query = example.Query
			}
			var docs []*Doc
			err := Where(context.Background(), DocsCollection, query, &docs)
			require.NoError(t, err)
			require.Len(t, docs, example.Count)
		})
	}
}
