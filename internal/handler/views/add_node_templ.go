// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.906
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func AddNode() templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!-- Hero --><section class=\"relative overflow-hidden pt-6\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = heroGradient().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<div class=\"relative z-10\"><div class=\"max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-10 lg:py-16\"><div class=\"text-center\"><div class=\"mt-5\"><h1 class=\"block font-extrabold text-4xl md:text-5xl lg:text-6xl text-neutral-200\">Add Monero Node</h1></div><div class=\"mt-5\"><p class=\"text-lg text-neutral-300\">You can use this page to add known remote node to the system so my bots can monitor it.</p></div></div><hr class=\"my-6 border-orange-400 mx-auto max-w-3xl\"><div class=\"max-w-4xl mx-auto px-4\"><div class=\"p-4 bg-blue-800/10 border border-blue-900 text-sm text-white rounded-lg\" role=\"alert\" tabindex=\"-1\" aria-labelledby=\"add-node-notice\"><div class=\"flex\"><div class=\"ms-4\"><h2 id=\"add-node-notice\" class=\"text-xl font-bold text-center\">Important Note</h2><div class=\"mt-2 text-sm\"><ul class=\"list-disc space-y-1 ps-5\"><li>As an administrator of this instance, I have full rights to delete, and blacklist any submitted node with or without providing any reason.</li><li>I2P nodes monitoring is beta.</li></ul></div></div></div></div></div><div class=\"max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-6\"><p class=\"mt-1 text-center\">Enter your Monero node information below:</p><div class=\"mt-12\"><form method=\"put\" hx-swap=\"transition:true\" hx-target=\"#form-result\" hx-disabled-elt=\".form\" hx-on::after-request=\"this.reset()\"><div class=\"grid grid-cols-1 sm:grid-cols-4 gap-6\"><div><label for=\"protocol\" class=\"block text-neutral-200\">Protocol *</label> <select id=\"protocol\" name=\"protocol\" class=\"frameless form\" autocomplete=\"off\"><option value=\"http\">HTTP</option> <option value=\"https\">HTTPS</option></select></div><div class=\"md:col-span-2\"><label for=\"hostname\" class=\"block text-neutral-200\">Host / IP *</label> <input type=\"text\" name=\"hostname\" id=\"hostname\" class=\"frameless form\" autocomplete=\"off\" placeholder=\"Eg: node.example.com or 172.16.17.18\" required></div><div><label for=\"port\" class=\"block text-neutral-200\">Port *</label> <input type=\"text\" name=\"port\" id=\"port\" class=\"frameless form\" autocomplete=\"off\" placeholder=\"Eg: 18081\" required></div></div><div class=\"mt-6 grid\"><button type=\"submit\" class=\"form w-full py-3 px-4 inline-flex justify-center items-center gap-x-2 text-sm font-bold rounded-lg border border-transparent bg-orange-600 text-white hover:bg-orange-500 focus:outline-none disabled:opacity-60 disabled:pointer-events-none\">Submit</button></div></form><div id=\"form-result\" class=\"max-w-4xl mx-auto my-6\"></div><div class=\"mt-3 text-center\"><p class=\"text-sm text-gray-500 dark:text-neutral-500\">Existing remote nodes can be found in <a href=\"/remote-nodes\" class=\"link\">/remote-nodes</a> page.</p></div></div></div></div></div></section><!-- End Hero -->")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
