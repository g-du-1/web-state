export const isUrlDisallowed = async (url: string) => {
  const { whitelistSites } = await chrome.storage.local.get("whitelistSites");

  if (whitelistSites === "") return true;

  let result = true;

  const whitelistedArr = whitelistSites?.split(",");

  whitelistedArr.forEach((site: string) => {
    if (url.includes(site)) {
      result = false;
    }
  })

  return result;
};
