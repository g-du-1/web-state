// src/popup.ts
var apiInput = document.getElementById("api-url");
if (apiInput) {
  apiInput.addEventListener("input", (event) => {
    const target = event.target;
    chrome.storage.local.set({ pageStateApiUrl: target.value });
  });
}
chrome.storage.local.get("pageStateApiUrl", (result) => {
  if (result.pageStateApiUrl) {
    const apiInput2 = document.getElementById("api-url");
    if (apiInput2) {
      apiInput2.value = result.pageStateApiUrl;
    }
  }
});
