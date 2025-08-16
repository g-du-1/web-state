const apiInput = document.getElementById("api-url");
const whitelistInput = document.getElementById("whitelist-sites");

if (apiInput) {
  apiInput.addEventListener("input", (event) => {
    const target = event.target as HTMLInputElement;

    chrome.storage.local.set({ pageStateApiUrl: target.value });
  });
}

if (whitelistInput) {
  whitelistInput.addEventListener("input", (event) => {
    const target = event.target as HTMLInputElement;

    chrome.storage.local.set({ whitelistSites: target.value });
  });
}

chrome.storage.local.get(["pageStateApiUrl", "whitelistSites"], (result) => {
  if (result.pageStateApiUrl) {
    const apiInput = document.getElementById("api-url");

    if (apiInput) {
      (apiInput as HTMLInputElement).value = result.pageStateApiUrl;
    }
  }

  if (result.whitelistSites) {
    const whitelistInput = document.getElementById("whitelist-sites");

    if (whitelistInput) {
      (whitelistInput as HTMLInputElement).value = result.whitelistSites;
    }
  }
});
