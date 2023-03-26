package napi

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/gofiber/fiber/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

// RoutePrinter route printer struct
type RoutePrinter struct {
	app   *fiber.App
	lock  *sync.RWMutex
	items []routeItem
}

// routeItem is used as a routing store
type routeItem struct {
	path     string
	methods  []string
	params   []string
	handlers []string
}

// routeMessage is used for debug printing
type routeMessage struct {
	name     string
	method   string
	path     string
	handlers string
}

// NewRoutePrinter create a new route printer
func NewRoutePrinter(app *fiber.App) *RoutePrinter {
	return &RoutePrinter{items: make([]routeItem, 0), app: app, lock: new(sync.RWMutex)}
}

// Hydrate iterates through the app route stack and hydrates acceptable router item structs.
func (p *RoutePrinter) Hydrate() *RoutePrinter {
	for _, routes := range p.app.Stack() {
		for _, route := range routes {
			p.addItem(route.Method, route.Path, route.Params, route.Handlers)
		}
	}

	return p
}

// Len get total routes
func (p *RoutePrinter) Len() int {
	return len(p.items)
}

// PrintPretty prints a simplified and color styled table
func (p *RoutePrinter) PrintPretty() *RoutePrinter {
	if len(p.items) == 0 {
		p.Hydrate()
	}

	sort.Slice(p.items, func(i, j int) bool {
		return p.items[i].path < p.items[j].path
	})

	data := make([][]string, len(p.items))
	dataIdx := 0
	for _, item := range p.items {
		var ctrls []string

		for _, handler := range item.handlers {
			c, f := parseFiberHandler(handler)
			if c != "" && f != "" {
				ctrls = append(ctrls, fmt.Sprintf("%s.%s", c, f))
			}
		}
		data[dataIdx] = []string{
			strings.Join(item.methods, "|"),
			item.path,
			strings.Join(item.params, ","),
			strings.Join(ctrls, ","),
		}
		dataIdx++
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Method", "URI", "Parameters", "Handlers"})
	t.Style().Format = table.FormatOptionsDefault
	t.Style().Options = table.OptionsNoBordersAndSeparators
	t.Style().Color.Header = text.Colors{text.BgHiGreen, text.FgBlack}

	for _, v := range data {
		if len(v) == 4 {
			t.AppendRow(table.Row{v[0], v[1], v[2], v[3]})
		}
	}
	t.Render() // Send output

	return p
}

// PrintDebug is ripped from fiber test codebase. Gives a more detailed route explaination.
func (p *RoutePrinter) PrintDebug() *RoutePrinter {
	const (
		// cBlack = "\u001b[90m"
		// cRed   = "\u001b[91m"
		cCyan   = "\u001b[96m"
		cGreen  = "\u001b[92m"
		cYellow = "\u001b[93m"
		cBlue   = "\u001b[94m"
		// cMagenta = "\u001b[95m"
		cWhite = "\u001b[97m"
		// cReset = "\u001b[0m"
	)
	var routes []routeMessage
	for _, routeStack := range p.app.Stack() {
		for _, route := range routeStack {
			var newRoute routeMessage
			newRoute.name = route.Name
			newRoute.method = route.Method
			newRoute.path = route.Path
			for _, handler := range route.Handlers {
				newRoute.handlers += runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name() + " "
			}
			routes = append(routes, newRoute)
		}
	}

	out := colorable.NewColorableStdout()
	if os.Getenv("TERM") == "dumb" || os.Getenv("NO_COLOR") == "1" || (!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd())) {
		out = colorable.NewNonColorable(os.Stdout)
	}

	w := tabwriter.NewWriter(out, 1, 1, 1, ' ', 0)
	// Sort routes by path
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].path < routes[j].path
	})
	_, _ = fmt.Fprintf(w, "%smethod\t%s| %spath\t%s| %sname\t%s| %shandlers\n", cBlue, cWhite, cGreen, cWhite, cCyan, cWhite, cYellow)
	_, _ = fmt.Fprintf(w, "%s------\t%s| %s----\t%s| %s----\t%s| %s--------\n", cBlue, cWhite, cGreen, cWhite, cCyan, cWhite, cYellow)
	for _, route := range routes {
		_, _ = fmt.Fprintf(w, "%s%s\t%s| %s%s\t%s| %s%s\t%s| %s%s\n", cBlue, route.method, cWhite, cGreen, route.path, cWhite, cCyan, route.name, cWhite, cYellow, route.handlers)
	}
	_ = w.Flush()

	return p
}

// addItem scans the given handlers and uses reflection to create human readable route items.
func (p *RoutePrinter) addItem(method string, path string, params []string, handlers []fiber.Handler) {
	p.lock.Lock()
	defer p.lock.Unlock()

	ifound := false
	for idx, item := range p.items {
		if item.path == path {
			ifound = true
			mfound := false
			for _, m := range item.methods {
				if m == method {
					mfound = true
					break
				}
			}

			if !mfound {
				if method != "CONNECT" && method != "OPTIONS" && method != "TRACE" {
					p.items[idx].methods = append(p.items[idx].methods, method)
				}
			}
			break
		}
	}

	if !ifound {
		var h []string
		for _, handler := range handlers {
			h = append(h, runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
		}

		i := routeItem{
			path:     path,
			methods:  []string{method},
			params:   params,
			handlers: h,
		}
		p.items = append(p.items, i)
	} else {
		for i, item := range p.items {
			if item.path == path && method != "HEAD" {
				for _, handler := range handlers {
					p.items[i].handlers = append(p.items[i].handlers, runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
				}
			}
		}
	}
}

// parseFiberHandler parse the handler from fiber's debug route data. See tests for more details
func parseFiberHandler(h string) (string, string) {
	p := ".([^.]+).([^.]+).func"
	rgx, err := regexp.Compile(p)
	if err != nil {
		return "", ""
	}

	matches := rgx.FindStringSubmatch(h)
	if len(matches) < 3 {
		return "", ""
	}

	if strings.Contains(matches[1], "/") {
		matches[1] = ""
		matches[2] = ""
	}

	return matches[1], matches[2]
}

// PrintRoutes this is the most useful function in our programs. Usage: rprint.PrintAuths(srv.App(), true)
func PrintRoutes(app *fiber.App, debug bool) {
	p := NewRoutePrinter(app)
	if debug {
		p.PrintDebug()
	} else {
		p.PrintPretty()
	}
}
