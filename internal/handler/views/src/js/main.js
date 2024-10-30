import "@preline/collapse";

htmx.onLoad(function () {
  // Auto init preline JS, see https://preline.co/docs/preline-javascript.html
  // This need to be inside `htmx.onLoad` to be work together with hx-boost.
  HSCollapse.autoInit();
});
