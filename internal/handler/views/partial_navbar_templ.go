// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.857
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func navbar(pageIdentifier string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<header class=\"fixed top-4 inset-x-0 flex flex-wrap z-50 w-full md:justify-start md:flex-nowrap before:absolute before:inset-0 before:max-w-7xl before:mx-2 before:lg:mx-auto before:rounded-[26px] before:bg-neutral-800/70 before:backdrop-blur-sm before:shadow-md before:shadow-orange-400/40\"><nav class=\"relative max-w-7xl w-full py-2.5 px-5 md:flex md:items-center md:justify-between md:py-0 mx-2 lg:mx-auto\"><div class=\"flex items-center justify-between\"><div class=\"flex-none inline-block\"><a class=\"text-xl font-semibold text-white focus:outline-none\" href=\"/\" aria-label=\"XMR Nodes\">XMR Nodes</a><div id=\"hx-indicator-main\" class=\"htmx-indicator animate-spin ml-2 inline-block size-4 border-[3px] border-current border-t-transparent text-orange-400 rounded-full\" role=\"status\" aria-label=\"loading indicator\"><span class=\"sr-only\">Loading...</span></div></div><div class=\"md:hidden\"><button type=\"button\" class=\"hs-collapse-toggle size-8 flex justify-center items-center text-sm font-semibold rounded-full bg-neutral-800 text-white disabled:opacity-50 disabled:pointer-events-none\" id=\"hs-navbar-floating-dark-collapse\" aria-expanded=\"false\" aria-controls=\"hs-navbar-floating-dark\" aria-label=\"Toggle navigation\" data-hs-collapse=\"#main-navbar\"><svg class=\"hs-collapse-open:hidden shrink-0 size-4\" xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><line x1=\"3\" x2=\"21\" y1=\"6\" y2=\"6\"></line><line x1=\"3\" x2=\"21\" y1=\"12\" y2=\"12\"></line><line x1=\"3\" x2=\"21\" y1=\"18\" y2=\"18\"></line></svg> <svg class=\"hs-collapse-open:block hidden shrink-0 size-4\" xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><path d=\"M18 6 6 18\"></path><path d=\"m6 6 12 12\"></path></svg></button></div></div><div id=\"main-navbar\" class=\"hs-collapse hidden overflow-hidden transition-all duration-300 basis-full grow md:block\" aria-labelledby=\"main-navbar-collapse\"><div class=\"flex flex-col md:flex-row md:items-center md:justify-end gap-2 md:gap-3 mt-3 md:mt-0 py-2 md:py-0 md:ps-7\"><a href=\"/\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if pageIdentifier == "/" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, " class=\"active\" aria-current=\"page\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, ">Home</a> <a href=\"/remote-nodes\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if pageIdentifier == "/remote-nodes" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, " class=\"active\" aria-current=\"page\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, ">Remote Nodes</a> <a href=\"/add-node\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if pageIdentifier == "/add-node" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, " class=\"active\" aria-current=\"page\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, ">Add Node</a></div></div></nav></header>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
