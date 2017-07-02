package myreviewer

import(
	"bytes"
	"html/template"
)

func reviewTemplateHandler(message string, team []*Team, reviews []*Review) bytes.Buffer {
	var content bytes.Buffer

	reviewScripts := []string{}
	reviewScripts = append(reviewScripts, "/assets/scripts/review.js")

	args := getDefaultTemplate(reviewScripts, []string{})

	args["message"] = message
	args["team"] = team
	args["review"] = reviews

	templates.ExecuteTemplate(&content, "review.html", args)

	return content
}


func manageTemplateHandler(message string, team []*Team) bytes.Buffer {
	var content bytes.Buffer

	manageScripts := []string{}
	manageScripts = append(manageScripts, "/assets/scripts/manage.js")

	args := getDefaultTemplate(manageScripts, []string{})

	args["message"] = message
	args["team"] = team

	templates.ExecuteTemplate(&content, "manage.html", args)

	return content
}


func getDefaultTemplate(additionalScripts []string, additionalCss []string) map[string]interface{} {
	var header, navigation, footer bytes.Buffer

	//init args for template
	args := make(map[string]interface{})
	headerArgs := make(map[string]interface{})

	scripts := []string{}
	css := []string{}

	//append global js. global css first.
	css = append(css, "/assets/css/global.css")
	

	//append global js. global js first.
	scripts = append(scripts, "/assets/scripts/jquery-3.2.1.min.js")

	//append additional asset
	for _, v := range additionalScripts {
		scripts = append(scripts, v)
	}
	for _, v := range additionalCss {
		css = append(css, v)
	}

	//assign args for header
	headerArgs["scripts"] = scripts
	headerArgs["css"] = css

	templates.ExecuteTemplate(&header, "header.html", headerArgs)
	templates.ExecuteTemplate(&footer, "footer.html", nil)
	templates.ExecuteTemplate(&navigation, "navigation.html", nil)

	args["header"] = template.HTML(header.String())
	args["footer"] = template.HTML(footer.String())
	args["navigation"] = template.HTML(navigation.String())

	return args
}

func closeTemplateHandler() bytes.Buffer {
	var content bytes.Buffer

	templates.ExecuteTemplate(&content, "close_window.html", nil)

	return content
}