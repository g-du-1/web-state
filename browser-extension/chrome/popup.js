const apiInput = document.getElementById("api-url");

if (apiInput) {
  apiInput.addEventListener("input", (event) => {
    const target = /** @type {HTMLInputElement} */ (event.target);

    chrome.storage.local.set({ redditPageStateUrl: target.value });
  });
}

chrome.storage.local.get("redditPageStateUrl", (result) => {
  if (result.redditPageStateUrl) {
    const apiInput = document.getElementById("api-url");

    if (apiInput) {
      /** @type {HTMLInputElement} */ (apiInput).value =
        result.redditPageStateUrl;
    }
  }
});
