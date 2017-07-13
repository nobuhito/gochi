package gochi

// type task struct {
// 	method    string
// 	key       string
// 	delayFunc *delay.Function
// }
//
// func (g *Gochi) registTask(method, key string, f interface{}) {
// 	t := task{
// 		method:    method,
// 		key:       strings.Split(key, "?")[0],
// 		delayFunc: delay.Func(uuid.NewV4().String(), f),
// 	}
// 	g.Tasks = append(g.Tasks, t)
// }
//
// func (g *Gochi) SearchTask(r *http.Request) (*delay.Function, error) {
// 	match := mux.RouteMatch{}
// 	if g.Router.Match(r, &match) {
// 		for _, task := range g.Tasks {
// 			t, _ := match.Route.GetPathTemplate()
// 			if t == task.key && r.Method == task.method {
// 				return task.delayFunc, nil
// 			}
// 		}
// 	}
// 	return nil, errors.New("not registed task")
// }
//
// func (g *Gochi) DelayGET(
// 	path string,
// 	h func(ctx context.Context, r *http.Request) Response,
// 	f interface{},
// ) {
// 	g.registTask("GET", path, f)
// 	g.Router.HandleFunc(path, responseToHandler(g, h)).Methods("GET")
// }
//
// func (g *Gochi) DelayPOST(
// 	path string,
// 	h func(ctx context.Context, r *http.Request) Response,
// 	f interface{},
// ) {
// 	g.registTask("POST", path, f)
// 	g.Router.HandleFunc(path, responseToHandler(g, h)).Methods("POST")
// }
//
// func (g *Gochi) DelayPUT(
// 	path string,
// 	h func(ctx context.Context, r *http.Request) Response,
// 	f interface{},
// ) {
// 	g.registTask("PUT", path, f)
// 	g.Router.HandleFunc(path, responseToHandler(g, h)).Methods("PUT")
// }
//
// func (g *Gochi) DelayDELETE(
// 	path string,
// 	h func(ctx context.Context, r *http.Request) Response,
// 	f interface{},
// ) {
// 	g.registTask("DELETE", path, f)
// 	g.Router.HandleFunc(path, responseToHandler(g, h)).Methods("DELETE")
// }
//
// func (g *Group) DelayGET(
// 	path string,
// 	h func(ctx context.Context, r *http.Request) Response,
// 	f interface{},
// ) {
// 	g.Gochi.registTask("GET", g.Parent+path, f)
// 	g.Gochi.Router.PathPrefix(g.Parent).Subrouter().HandleFunc(path, responseToHandler(g.Gochi, h)).Methods("GET")
// }
//
// func (g *Group) DelayPOST(
// 	path string,
// 	h func(ctx context.Context, r *http.Request) Response,
// 	f interface{},
// ) {
// 	g.Gochi.registTask("POST", g.Parent+path, f)
// 	g.Gochi.Router.PathPrefix(g.Parent).Subrouter().HandleFunc(path, responseToHandler(g.Gochi, h)).Methods("POST")
// }
//
// func (g *Group) DelayPUT(
// 	path string,
// 	h func(ctx context.Context, r *http.Request) Response,
// 	f interface{},
// ) {
// 	g.Gochi.registTask("PUT", g.Parent+path, f)
// 	g.Gochi.Router.PathPrefix(g.Parent).Subrouter().HandleFunc(path, responseToHandler(g.Gochi, h)).Methods("PUT")
// }
//
// func (g *Group) DelayDELETE(
// 	path string,
// 	h func(ctx context.Context, r *http.Request) Response,
// 	f interface{},
// ) {
// 	g.Gochi.registTask("DELETE", g.Parent+path, f)
// 	g.Gochi.Router.PathPrefix(g.Parent).Subrouter().HandleFunc(path, responseToHandler(g.Gochi, h)).Methods("DELETE")
// }
