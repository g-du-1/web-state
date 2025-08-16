const apiInput = document.getElementById("api-url");

if (apiInput) {
  apiInput.addEventListener("input", (event) => {
    const target = event.target as HTMLInputElement;

    chrome.storage.local.set({ pageStateApiUrl: target.value });
  });
}

chrome.storage.local.get("pageStateApiUrl", (result) => {
  if (result.pageStateApiUrl) {
    const apiInput = document.getElementById("api-url");

    if (apiInput) {
      (apiInput as HTMLInputElement).value = result.pageStateApiUrl;
    }
  }
});
