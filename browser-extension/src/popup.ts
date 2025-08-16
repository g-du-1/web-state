const apiInput = document.getElementById("api-url");

if (apiInput) {
  apiInput.addEventListener("input", (event) => {
    const target = event.target as HTMLInputElement;

    chrome.storage.local.set({ redditPageStateUrl: target.value });
  });
}

chrome.storage.local.get("redditPageStateUrl", (result) => {
  if (result.redditPageStateUrl) {
    const apiInput = document.getElementById("api-url");

    if (apiInput) {
      (apiInput as HTMLInputElement).value = result.redditPageStateUrl;
    }
  }
});
