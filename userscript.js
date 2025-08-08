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

  let latestState = null;

  const createButton = () => {
    const button = document.createElement("button");

    button.textContent = "Show State";

    button.style.cssText = `
      position: fixed;
      bottom: 20px;
      right: 20px;
      z-index: 10000;
      padding: 10px;
      background: #0079d3;
      color: white;
      border: none;
      border-radius: 5px;
      cursor: pointer;
    `;

    button.onclick = () => {
      if (latestState) {
        showStateInfo(latestState);
      } else {
        alert("No saved state found");
      }
    };

    document.body.appendChild(button);
  };

  const showStateInfo = (state) => {
    const currentScroll = Math.trunc(window.scrollY);

    alert(
      `Saved State:\nScroll: ${state.scrollPos}px\nText: ${
        state.visibleText || "N/A"
      }\nCurrent Scroll: ${currentScroll}px`
    );
  };

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
        latestState = JSON.parse(response.responseText);
        console.log("Latest Page State:", latestState);
      },
      onerror: (error) => {
        console.error("Latest Page State:", error);
      },
    });
  };

  window.addEventListener("load", () => {
    createButton();
    getLatestPageState();
  });

  window.addEventListener("scrollend", savePageState);
})();
