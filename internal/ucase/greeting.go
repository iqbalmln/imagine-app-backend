package ucase

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.privy.id/go_graphql/internal/appctx"
)

type greeting struct{}

func (greeting) Serve(data *appctx.Data) appctx.Response {
	// name := data.Request.URL.Query().Get("name")
	// params := mux.Vars(data.Request)
	// name := params["name"]

	/* Key Value Form Body Text*/
	// name := data.Request.FormValue("name")

	/* Key Value Form Body File (Image)*/
	// f, header, err := data.Request.FormFile("img")
	// if err != nil {
	// 	// Do Error Handling
	// }
	// defer f.Close()

	// buf := new(bytes.Buffer)

	// buf.ReadFrom(f)

	// err = os.WriteFile(header.Filename, buf.Bytes(), 0644)
	// if err != nil {
	// 	// Do Error Handling
	// }

	/* Read Data with JSON type */
	var payload struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(data.Request.Body).Decode(&payload)
	if err != nil {
		// Do Error Handling
	}

	return *appctx.NewResponse().
		WithCode(http.StatusOK).
		WithData(map[string]string{
			"message": fmt.Sprintf("Hello %s!", payload.Name),
		})
}

func NewGreeting() *greeting {
	return &greeting{}
}
