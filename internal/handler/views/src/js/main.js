import "@preline/collapse";
import "@preline/overlay";

window.addEventListener("load", () => {
  var clipboard = new ClipboardJS(".clipboard");
  clipboard.on("success", function (e) {
    let btnText = e.trigger.textContent;
    let successText = e.trigger.getAttribute("data-success-text");
    if (successText === null) {
      successText = "Copied 👍";
    }
    e.trigger.textContent = successText;
    e.trigger.disabled = true;
    setTimeout(function () {
      e.trigger.textContent = btnText;
      e.trigger.disabled = false;
    }, 1000);
  });
  clipboard.on("error", function (e) {
    console.error("Clipboard error", e.trigger);
  });
});

htmx.onLoad(function () {
  // Auto init preline JS, see https://preline.co/docs/preline-javascript.html
  // This need to be inside `htmx.onLoad` to be work together with hx-boost.
  HSCollapse.autoInit();
  HSOverlay.autoInit();
});
