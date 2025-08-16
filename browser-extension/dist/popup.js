// src/popup.ts
var apiInput = document.getElementById("api-url");
var whitelistInput = document.getElementById("whitelist-sites");
if (apiInput) {
  apiInput.addEventListener("input", (event) => {
    const target = event.target;
    chrome.storage.local.set({ pageStateApiUrl: target.value });
  });
}
if (whitelistInput) {
  whitelistInput.addEventListener("input", (event) => {
    const target = event.target;
    chrome.storage.local.set({ whitelistSites: target.value });
  });
}
chrome.storage.local.get(["pageStateApiUrl", "whitelistSites"], (result) => {
  if (result.pageStateApiUrl) {
    const apiInput2 = document.getElementById("api-url");
    if (apiInput2) {
      apiInput2.value = result.pageStateApiUrl;
    }
  }
  if (result.whitelistSites) {
    const whitelistInput2 = document.getElementById("whitelist-sites");
    if (whitelistInput2) {
      whitelistInput2.value = result.whitelistSites;
    }
  }
});
