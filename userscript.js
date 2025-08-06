// ==UserScript==
// @name         Scroll State Saver
// @namespace    http://tampermonkey.net/
// @version      2025-08-06
// @description  try to take over the world!
// @author       You
// @match        https://www.reddit.com/*
// @icon         https://www.google.com/s2/favicons?sz=64&domain=reddit.com
// @grant        GM_xmlhttpRequest
// ==/UserScript==

(function () {
  "use strict";

  const isElementInViewport = (element) => {
    const rect = element.getBoundingClientRect();

    return (
      rect.top >= 0 &&
      rect.left >= 0 &&
      rect.bottom <=
        (window.innerHeight || document.documentElement.clientHeight) &&
      rect.right <= (window.innerWidth || document.documentElement.clientWidth)
    );
  };

  const getElementsInMiddle = () => {
    const viewportMidpointY = window.innerHeight / 2;

    const elementsAtMidpoint = document.elementsFromPoint(
      window.innerWidth / 2,
      viewportMidpointY
    );

    return elementsAtMidpoint;
  };

  const captureVisibleText = () => {
    const elementsInMiddle = getElementsInMiddle();

    const visibleText = Array.from(elementsInMiddle).reduce((text, element) => {
      if (isElementInViewport(element)) {
        text += element.textContent.trim() + " ";
      }
      return text;
    }, "");

    return visibleText.replace(/\s+/g, " ");
  };

  const savePageState = async () => {
    const payload = {
      url: window.location.href,
      scrollPos: Math.trunc(window.scrollY),
    };

    GM_xmlhttpRequest({
      method: "POST",
      url: "http://localhost:8080/pagestate",
      data: JSON.stringify(payload),
      headers: {
        "Content-Type": "application/json",
      },
      onload: function (response) {
        console.log("Success:", response.responseText);
      },
      onerror: function (error) {
        console.error("Error:", error);
      },
    });
  };

  document.addEventListener("scrollend", async () => {
    console.log("Visible: ", captureVisibleText());
    savePageState;
  });
})();
