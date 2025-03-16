package main

import (
	"os"

	"github.com/lainio/err2/try"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

func main() {
	// try.To()
	os.RemoveAll("pb_data")
	// os.Args = append(os.Args, "serve")

	app := pocketbase.New()

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		e.InstallerFunc = func(app core.App, systemSuperuser *core.Record, baseURL string) error {
			return nil
		}
		initPublicCollection(app)
		initAppSetting(app)
		initUsers(app)
		return e.Next()
	})

	app.OnCollectionAfterDeleteError()

	try.To(app.Start())
}

func initPublicCollection(app core.App) {
	collection := core.NewBaseCollection("public")
	collection.ListRule = types.Pointer("")
	collection.ViewRule = types.Pointer("")
	collection.CreateRule = types.Pointer("")
	collection.UpdateRule = types.Pointer("")
	collection.DeleteRule = types.Pointer("")
	collection.Fields.Add(
		&core.TextField{
			Name:     "name",
			Required: true,
		},
		&core.AutodateField{
			Name:     "created",
			OnCreate: true,
		},
		&core.AutodateField{
			Name:     "updated",
			OnCreate: true,
			OnUpdate: true,
		},
	)
	try.To(app.Save(collection))
}

func initAppSetting(app core.App) {
	collection := try.To1(app.FindCollectionByNameOrId("users"))
	collection.AuthToken.Duration = 10
	try.To(app.Save(collection))
}

func initUsers(app core.App) {
	collection := try.To1(app.FindCollectionByNameOrId("users"))
	record := core.NewRecord(collection)
	record.Load(map[string]any{
		"email":           "test@test.invaild",
		"username":        "test",
		"password":        "testtest",
		"PasswordConfirm": "testtest",
	})
	try.To(app.Save(record))
}
