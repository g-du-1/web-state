// ==UserScript==
// @name         Scroll State Saver
// @namespace    http://tampermonkey.net/
// @version      2025-08-06
// @description  Saves page state.
// @author       You
// @match        https://www.reddit.com/*
// @icon         https://www.google.com/s2/favicons?sz=64&domain=reddit.com
// @grant        GM_xmlhttpRequest
// ==/UserScript==

(function () {
  "use strict";

  const baseUrl = "http://localhost:8080";
  const saveUrl = `${baseUrl}/pagestate`;
  const getLatestUrl = `${baseUrl}/pagestate/latest`;

  const savePageState = () => {
    const payload = {
      url: window.location.href,
      scrollPos: Math.trunc(window.scrollY),
      visibleText: "TODO: Not implemented in js",
    };

    GM_xmlhttpRequest({
      method: "POST",
      url: saveUrl,
      data: JSON.stringify(payload),
      headers: {
        "Content-Type": "application/json",
      },
      onload: (response) => {
        console.log("Saved Page State:", JSON.parse(response.responseText));
      },
      onerror: (error) => {
        console.error("Page State Save:", error);
      },
    });
  };

  const getLatestPageState = () => {
    GM_xmlhttpRequest({
      method: "GET",
      url: getLatestUrl + "?url=" + encodeURIComponent(window.location.href),
      onload: (response) => {
        console.log("Latest Page State:", JSON.parse(response.responseText));
      },
      onerror: (error) => {
        console.error("Latest Page State:", error);
      },
    });
  };

  window.addEventListener("load", getLatestPageState);
  window.addEventListener("scrollend", savePageState);
})();
