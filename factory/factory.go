package factory

import (
	"gorm.io/gorm"
	"log"
	"reflect"
)

type Factory struct {
	factories factoryMap
	sql       *gorm.DB
}

// New instantiates a new factory struct with *gorm.DB
func New(sql *gorm.DB) *Factory {
	return &Factory{
		factories: factoryMap{},
		sql:       sql,
	}
}

// Add a new factory into the factoryMap.
func (f *Factory) Add(model ...interface{}) *Factory {
	for _, m := range model {
		f.factories.set(m, m.(IFactory))
	}
	return f
}

// Make will make a new model interface without saving to the database
func (f *Factory) Make(model interface{}) interface{} {
	if fac, ok := f.factories[getModelType(model)]; ok {
		model = overwriteStructFields(fac.Make(), model)
		return model
	}
	return nil
}

// Create will use the underlying gorm.DB and create a new model
func (f *Factory) Create(model interface{}) interface{} {
	model = f.Make(model)
	if tx := f.sql.Create(model); tx.Error != nil {
		log.Fatalln(tx.Error)
	}
	return model
}

type factoryMap map[string]IFactory

// set sets the factory for a given model
func (f factoryMap) set(model interface{}, factory IFactory) {
	f[getModelType(model)] = factory
}

// getModelType probably don't need this anymore because of the new interface{} mapping. Tuck it away for now.
func getModelType(model interface{}) string {
	value := reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr && value.IsNil() {
		value = reflect.New(value.Type().Elem())
	}
	modelType := reflect.Indirect(value).Type()

	if modelType.Kind() == reflect.Interface {
		modelType = reflect.Indirect(reflect.ValueOf(model)).Elem().Type()
	}

	return modelType.String()
}

// OverrideEmptyString useful in cases when you want to override empty strings. reflection causes issues without this.
const OverrideEmptyString = "OVERRIDE_EMPTY_STRING"

// overwriteStructFields src is the original model, dest are the new fields. must both be pointers
func overwriteStructFields(src interface{}, dest interface{}) interface{} {
	vModel := reflect.ValueOf(src).Elem()
	vOverride := reflect.ValueOf(dest)

	// vEmpty is an empty clone of the original model to ensure unused fields don't override unintended
	vEmpty := reflect.New(vModel.Type()).Elem()

	for i := 0; i < vModel.NumField(); i++ {
		fvEmpty := vEmpty.Field(i)
		fvOverride := vOverride.Field(i)

		if fvOverride.CanInterface() && fvOverride.IsValid() {
			fvModel := vModel.Field(i)
			if fvOverride.String() == OverrideEmptyString {
				fvModel.SetString("")
			} else if fvEmpty.Kind() == reflect.Struct && fvOverride.Kind() == reflect.Struct { // might evolve later
				if !reflect.DeepEqual(fvEmpty.Interface(), fvOverride.Interface()) {
					fvModel.Set(reflect.ValueOf(fvOverride.Interface()))
				}
			} else if fvEmpty.Kind() == reflect.Slice && fvOverride.Kind() == reflect.Slice {

			} else {
				if fvEmpty.Interface() != fvOverride.Interface() {
					fvModel.Set(reflect.ValueOf(fvOverride.Interface()))
				}
			}
		}

	}

	return src
}
