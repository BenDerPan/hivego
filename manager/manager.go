package manager

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/web"
	"github.com/rprp/hive/schedule"
	"log"
	"net/http"
	"time"
)

func StartManager(sl *schedule.ScheduleManager) {
	m := martini.Classic()
	m.Use(Logger)
	m.Use(martini.Static("web/public"))
	m.Use(web.ContextWithCookieSecret(""))
	m.Use(render.Renderer(render.Options{
		Directory:       "templates",                // Specify what path to load the templates from.
		Extensions:      []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		Delims:          render.Delims{"{[{", "}]}"},
		Charset:         "UTF-8",     // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON:      true,        // Output human readable JSON
		IndentXML:       true,        // Output human readable XML
		HTMLContentType: "text/html", // Output XHTML content type instead of default "text/html"
	}))

	m.Map(sl)
	controller(m)

	err := http.ListenAndServe(":3000", m)
	if err != nil {
		log.Fatal("Fail to start server: %v", err)
	}
}

func controller(m *martini.ClassicMartini) {
	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})

	m.Group("", func(r martini.Router) {
		r.Get("/schedules", getSchedules)
		r.Post("/schedules", addSchedule)
		r.Get("/schedules/ID", getScheduleById)
		r.Put("/schedules/:ID", updateSchedule)
		r.Delete("/schedules/ID", deleteSchedule)
	})

	m.Post("/hello", func(ctx *web.Context) {
		fmt.Println(ctx.Params)
		ctx.WriteString("Hello World!")
	})

}

func getSchedules(ctx *web.Context, r render.Render, res http.ResponseWriter, Ss *schedule.ScheduleManager) {
	sl := make([]*schedule.Schedule, 0)
	for _, s := range Ss.ScheduleList {
		d := &schedule.Schedule{}
		schedule.Copy(d, s)
		d.Job = nil
		sl = append(sl, d)
	}
	r.JSON(200, sl)

}

func getScheduleById(ctx *web.Context, res http.ResponseWriter, Ss *schedule.ScheduleManager) {
	for _, s := range Ss.ScheduleList {
		res.Write([]byte(s.String()))
	}

}

func addSchedule(ctx *web.Context, res http.ResponseWriter, Ss *schedule.ScheduleManager) {
	fmt.Println(ctx.Params)
	fmt.Println(ctx.Request)

}

func deleteSchedule(ctx *web.Context, res http.ResponseWriter, Ss *schedule.ScheduleManager) {
	for _, s := range Ss.ScheduleList {
		res.Write([]byte(s.String()))
	}

}

func updateSchedule(params martini.Params, ctx *web.Context, res http.ResponseWriter) {
	fmt.Println(params)

}

func Logger() martini.Handler { // {{{
	return func(res http.ResponseWriter, req *http.Request, ctx martini.Context, log *log.Logger) {

		start := time.Now()
		log.Printf("Started %s %s", req.Method, req.URL.Path)

		rw := res.(martini.ResponseWriter)
		ctx.Next()

		content := fmt.Sprintf("Completed %v %s in %v", rw.Status(), http.StatusText(rw.Status()), time.Since(start))
		switch rw.Status() {
		case 200:
			content = fmt.Sprintf("\033[1;32m%s\033[0m", content)
		case 304:
			content = fmt.Sprintf("\033[1;33m%s\033[0m", content)
		case 404:
			content = fmt.Sprintf("\033[1;31m%s\033[0m", content)
		case 500:
			content = fmt.Sprintf("\033[1;36m%s\033[0m", content)
		}
		log.Println(content)
	}
} // }}}