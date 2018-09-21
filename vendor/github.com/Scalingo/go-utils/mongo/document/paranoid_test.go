package document

import (
	"context"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const ParanoidDocsCollection = "paranoid_docs"

type ParanoidDoc struct {
	Paranoid `bson:",inline"`
}

func NewTestParanoidDoc(t *testing.T) (*ParanoidDoc, func()) {
	d := ParanoidDoc{}
	require.NoError(t, Save(context.Background(), ParanoidDocsCollection, &d))
	return &d, func() {
		require.NoError(t, ReallyDestroy(context.Background(), ParanoidDocsCollection, &d))
	}
}

func TestParanoid_Find(t *testing.T) {
	examples := []struct {
		Name        string
		ParanoidDoc func(t *testing.T) (*ParanoidDoc, func())
		Error       string
	}{
		{
			Name: "it should find existing doc",
			ParanoidDoc: func(t *testing.T) (*ParanoidDoc, func()) {
				d, clean := NewTestParanoidDoc(t)
				return d, clean
			},
		}, {
			Name: "it should not find unsaved doc",
			ParanoidDoc: func(t *testing.T) (*ParanoidDoc, func()) {
				d := &ParanoidDoc{}
				d.ID = bson.NewObjectId()
				return d, func() {}
			},
			Error: "not found",
		}, {
			Name: "it should not find destroyed doc",
			ParanoidDoc: func(t *testing.T) (*ParanoidDoc, func()) {
				d, clean := NewTestParanoidDoc(t)
				err := Destroy(context.Background(), ParanoidDocsCollection, d)
				if err != nil {
					clean()
					require.NoError(t, err)
				}
				return d, clean
			},
			Error: "not found",
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			fixtureParanoidDoc, clean := example.ParanoidDoc(t)
			defer clean()

			var d ParanoidDoc
			err := Find(context.Background(), ParanoidDocsCollection, fixtureParanoidDoc.ID, &d)
			if example.Error != "" {
				assert.Contains(t, err.Error(), example.Error)
			} else {
				require.NoError(t, err)
				require.Equal(t, fixtureParanoidDoc.ID, d.ID)
			}
		})
	}
}

func TestParanoid_FindUnscoped(t *testing.T) {
	examples := []struct {
		Name        string
		ParanoidDoc func(t *testing.T) (*ParanoidDoc, func())
		Error       string
	}{
		{
			Name: "it should find existing doc",
			ParanoidDoc: func(t *testing.T) (*ParanoidDoc, func()) {
				d, clean := NewTestParanoidDoc(t)
				return d, clean
			},
		}, {
			Name: "it should not find unsaved doc",
			ParanoidDoc: func(t *testing.T) (*ParanoidDoc, func()) {
				d := &ParanoidDoc{}
				d.ID = bson.NewObjectId()
				return d, func() {}
			},
			Error: "not found",
		}, {
			Name: "it should find destroyed doc",
			ParanoidDoc: func(t *testing.T) (*ParanoidDoc, func()) {
				d, clean := NewTestParanoidDoc(t)
				err := Destroy(context.Background(), ParanoidDocsCollection, d)
				require.NoError(t, err)
				return d, clean
			},
		}, {
			Name: "it should not find really destroyed doc",
			ParanoidDoc: func(t *testing.T) (*ParanoidDoc, func()) {
				d := ParanoidDoc{}
				err := Save(context.Background(), ParanoidDocsCollection, &d)
				require.NoError(t, err)
				err = ReallyDestroy(context.Background(), ParanoidDocsCollection, &d)
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
			fixtureParanoidDoc, clean := example.ParanoidDoc(t)
			defer clean()

			var d ParanoidDoc
			err := FindUnscoped(context.Background(), ParanoidDocsCollection, fixtureParanoidDoc.ID, &d)
			if example.Error != "" {
				assert.Contains(t, err.Error(), example.Error)
			} else {
				require.NoError(t, err)
				require.Equal(t, fixtureParanoidDoc.ID, d.ID)
			}
		})
	}
}

func TestParanoid_Where(t *testing.T) {
	examples := []struct {
		Name         string
		Query        bson.M
		ParanoidDocs func(t *testing.T) ([]*ParanoidDoc, func())
		Count        int
	}{
		{
			Name: "it should find existing documents",
			ParanoidDocs: func(t *testing.T) ([]*ParanoidDoc, func()) {
				d1, clean1 := NewTestParanoidDoc(t)
				d2, clean2 := NewTestParanoidDoc(t)
				return []*ParanoidDoc{d1, d2}, func() {
					clean1()
					clean2()
				}
			},
			Count: 2,
		}, {
			Name: "it should not find paranoia-deleted documents",
			ParanoidDocs: func(t *testing.T) ([]*ParanoidDoc, func()) {
				d1, clean1 := NewTestParanoidDoc(t)
				err := Destroy(context.Background(), ParanoidDocsCollection, d1)
				require.NoError(t, err)
				d2, clean2 := NewTestParanoidDoc(t)
				err = Destroy(context.Background(), ParanoidDocsCollection, d2)
				require.NoError(t, err)
				return []*ParanoidDoc{d1, d2}, func() {
					clean1()
					clean2()
				}
			},
			Count: 0,
		}, {
			Name:  "it should find deleted document, if queried specifically",
			Query: bson.M{"deleted_at": bson.M{"$exists": true}},
			ParanoidDocs: func(t *testing.T) ([]*ParanoidDoc, func()) {
				d1, clean1 := NewTestParanoidDoc(t)
				err := Destroy(context.Background(), ParanoidDocsCollection, d1)
				require.NoError(t, err)
				d2, clean2 := NewTestParanoidDoc(t)
				err = Destroy(context.Background(), ParanoidDocsCollection, d2)
				return []*ParanoidDoc{d1, d2}, func() {
					clean1()
					clean2()
				}
			},
			Count: 2,
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			_, clean := example.ParanoidDocs(t)
			defer clean()

			query := bson.M{}
			if example.Query != nil {
				query = example.Query
			}
			var docs []*ParanoidDoc
			err := Where(context.Background(), ParanoidDocsCollection, query, &docs)
			require.NoError(t, err)
			require.Len(t, docs, example.Count)
		})
	}
}

func TestParanoid_Restore(t *testing.T) {
	ctx := context.Background()
	doc, clean := NewTestParanoidDoc(t)
	defer clean()
	err := Destroy(ctx, ParanoidDocsCollection, doc)
	assert.NoError(t, err)

	err = FindUnscoped(ctx, ParanoidDocsCollection, doc.ID, doc)
	assert.NoError(t, err)
	assert.True(t, doc.IsDeleted())

	err = Restore(ctx, ParanoidDocsCollection, doc)
	assert.NoError(t, err)

	err = Find(ctx, ParanoidDocsCollection, doc.ID, doc)
	assert.NoError(t, err)
	assert.False(t, doc.IsDeleted())
}
