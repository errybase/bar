package bar

import (
	"reflect"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

type relation[T any] struct {
	Model        any
	RelationName string
}

func (r relation[T]) rel(db bun.IDB) *rel {
	for name, bunRel := range r.baseTable(db).Relations {
		if name == r.RelationName {
			return &rel{
				Relation: bunRel,
				v:        r.baseValue(),
			}
		}
	}
	return nil
}

func (r relation[T]) baseValue() reflect.Value {
	return reflect.ValueOf(r.Model).Elem()
}

func (r relation[T]) baseTable(db bun.IDB) *schema.Table {
	return db.Dialect().Tables().ByModel(r.baseModel())
}

func (r relation[T]) baseModel() string {
	return reflect.TypeOf(r.Model).Elem().Name()
}

type rel struct {
	*schema.Relation

	v reflect.Value
}

func (r rel) appendRelModel(t any) {
	v := reflect.ValueOf(t).Elem()
	for _, f := range r.fields() {
		v.FieldByName(f.join.GoName).Set(f.value)
	}
}

func (r rel) appendBaseModel(t any) {
	v := reflect.ValueOf(t)
	for _, f := range r.fields() {
		r.v.FieldByName(f.base.GoName).Set(v.FieldByName(f.join.GoName))
	}
}

func (r rel) basePKValues() []any {
	var values []any
	for _, pk := range r.BasePKs {
		values = append(values, r.v.FieldByName(pk.GoName).Interface())
	}
	return values
}

func (r rel) fields() []field {
	var fields []field
	for i, joinPK := range r.JoinPKs {
		basePK := r.BasePKs[i]
		fields = append(fields, field{
			base:  basePK,
			join:  joinPK,
			value: r.v.FieldByName(basePK.GoName),
		})
	}
	return fields
}

func (r rel) joinIdentCols() []bun.Ident {
	var cols []bun.Ident
	for _, pk := range r.JoinPKs {
		cols = append(cols, bun.Ident(pk.Name))
	}
	return cols
}

func (r rel) joinCols() []string {
	var cols []string
	for _, pk := range r.JoinPKs {
		cols = append(cols, pk.Name)
	}
	return cols
}

type field struct {
	base  *schema.Field
	join  *schema.Field
	value reflect.Value
}
