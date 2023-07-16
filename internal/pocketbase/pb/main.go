package main

import (
	"os"

	"github.com/lainio/err2/try"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func main() {
	// try.To()
	os.RemoveAll("pb_data")
	// os.Args = append(os.Args, "serve")

	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		initPublicCollection(app)
		initAppSetting(app)
		initUsers(app)
		return nil
	})

	try.To(app.Start())
}

func initPublicCollection(app *pocketbase.PocketBase) {
	collection := &models.Collection{
		Name: "public",
		Type: models.CollectionTypeBase,

		ListRule:   types.Pointer(""),
		ViewRule:   types.Pointer(""),
		CreateRule: types.Pointer(""),
		UpdateRule: types.Pointer(""),
		DeleteRule: types.Pointer(""),

		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name:     "name",
				Type:     schema.FieldTypeText,
				Required: true,
			},
		),
	}
	form := forms.NewCollectionUpsert(app, collection)
	try.To(form.Submit())
}

func initAppSetting(app *pocketbase.PocketBase) {
	form := forms.NewSettingsUpsert(app)
	form.RecordAuthToken.Duration = 5
	try.To(form.Submit())
}

func initUsers(app *pocketbase.PocketBase) {
	collection := try.To1(app.Dao().FindCollectionByNameOrId("users"))
	record := models.NewRecord(collection)
	form := forms.NewRecordUpsert(app, record)
	form.Username = "test"
	form.Password = "testtest"
	form.PasswordConfirm = "testtest"
	try.To(form.Submit())
}
